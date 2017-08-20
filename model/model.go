package model

import (
	"chessServer/utility/geometry"
	"chessServer/parser"
)


var(	//	definition of probable steps for figures, which can't go straight
	KingProbableStepList = [8]geometry.Point{{-1,-1},{-1,0},{-1,1},
											{0, -1},{0, 1},{1,-1},
											{1,0},{1,1},}
	KnightProbableStepList = [8]geometry.Point{{-2,-1},{-1, -2},{1, -2},
											  {2, -1},{-2, 1},{-1,2},
											  {1, 2},{2, 1},}
	PawnProbableShortStepList = [2]geometry.Point{{0,1},	//	probable short steps for black pawns
											   {0,-1},	//	probable short steps for white pawns
											  }
	PawnProbableLongStepList = [2]geometry.Point{{0,2},	//	probable steps for black pawns
												{0,-2},	//	probable steps for white pawns
												}
	PawnProbableAttackList = [2][2]geometry.Point{
		{{1,1}, {-1,1}},	//	probable steps for black pawns
		{{1,-1}, {-1,-1}},	//	probable steps for white pawns
	}
)
//	interface for all figures, just checking if figure can go somewhere/attack a field
type StepMaker interface{
	CheckStepAvailable(point geometry.Point)(bool)	//	check if figure can go to the following

	CheckAttackAvailable(point geometry.Point)(bool)	//	check if figure can attack the following field

	ListStepsAvailable()([]geometry.Point)	//	list all available fields to go
	AttacksAvailable()([]geometry.Point)	//	list all available fields to attack

	CheckForCollision(destination,obstacle geometry.Point)(bool)	//	check if we can collide with other figure while doing step

	RetCoords()geometry.Point	//	interface doesn't storage any data, but we sometimes will need coordinates and colours of all figures

	RetColour()Colour

	isEnemy(maker StepMaker)bool	//	check whereas other figure can make step or not
}

type Figure struct{	//	parent class for all figures(all figures inherits Figure and implement StepMaker)
	geometry.Point	//	figure coordinates
	Colour_ Colour	//	figure colour
}

func(f Figure)RetColour()Colour{	//	implement interfaces method
	return f.Colour_
}

func (f Figure)isEnemy(maker StepMaker)bool{	//	true if two figures are enemies
	return f.Colour_ != maker.RetColour()
}

func ConstructFigure(x,y int ,colour Colour)(Figure) {	//	constructor, returns value
	return Figure{Point: geometry.ConstructPoint(x,y),Colour_: colour}
}

func (f Figure)checkAvailable(available []geometry.Point, point geometry.Point)(bool){	//	utility function to check whereas given point is available for step
	for _,element:=range available{
		if element.Equal(point){
			return true
		}
	}
	return false
}

func (f Figure)RetCoords()geometry.Point{	//	return figure's coordinates
	return f.Point
}


// Parent class, representing behavior of linear figures:Rooks, Bishops and Queens
// Linear figures have list of directions vectors and just check if destination point
// lays on any of directions relative to self coordinates.
//To check collision it just checks if obstacle lays on given direction.
// (Direction = line, given by initial point(self coordinate) and unit vector)
type LinearFigure struct{
	Figure	//	figure data (base class)
	Direction []geometry.Vector	//	directions
}

func (l LinearFigure)ListStepsAvailable()(Buffer []geometry.Point){
	Buffer = make([]geometry.Point,0,16)
	for _,element:=range l.Direction{	//	check all directions
		i,j:=l.Point.Add(element.Point),l.Point.Subtract(element.Point)	//set iterators

		for i.CheckFieldBoundaries(){	//	iterate through all cells lying on given direction,(only in field boundaries)
			Buffer = append(Buffer,i)	// add all cells lying on given direction
			i = i.Add(element.Point)	//	add direction vector to iterator
		}
		for j.CheckFieldBoundaries(){	//	the same operation, but in opposite direction
			Buffer = append(Buffer,j)
			j = j.Subtract(element.Point)	//	subtract initial direction equals to adding opposite direction
		}
	}
	return Buffer
}

func (l LinearFigure)AttacksAvailable()(Buffer []geometry.Point){	//	linear figures make step and attack the same way
	return l.ListStepsAvailable()
}

func (l LinearFigure)CheckStepAvailable(point geometry.Point)(bool){	// check if step to given cell is permitted
	return l.checkAvailable(l.ListStepsAvailable(),point)
}

func (l LinearFigure)CheckAttackAvailable(point geometry.Point)(bool){	//	attack and step are the same
	return l.CheckStepAvailable(point)
}


// Collision happens then obstacle lays on way of linear figure from initial point to destination
func (l LinearFigure)CheckForCollision(destination, obstacle geometry.Point)(bool){
	way:=geometry.ConstructLine(l.Point, destination) // only valid destinations are checked for collision
	return way.Intersect(obstacle)                      // so we can skip checking for validity
}
// NonlinearFigure represents figures that can go to the same cells:Kings and Knights, relative to it's coordinates,
// they can both attack and step to the same fields.
// Nonlinear figures collide obstacles placed only in destination point.
// Note, that Pawns aren't neither linear, nor Nonlinear figures, since
// they can step just forward, but attack only by diagonal, and also they can make long steps first time,
// so collision case is enough complex to make other implementation for Pawns.
type NonLinearFigure struct{	//	Kings/Knights, figures, which can collide only if they rich collision place
	Figure
	ProbableSteps [8]geometry.Point	//	both kings and knights can visit only 8 places
}

//
func (n NonLinearFigure)ListStepsAvailable()(Buffer []geometry.Point){
	Buffer = make([]geometry.Point,0,8)
	for _,element:= range n.ProbableSteps{
		if n.Point.Add(element).CheckFieldBoundaries(){
			Buffer=append(Buffer,n.Point.Add(element))
		}
	}
	return Buffer
}

func (n NonLinearFigure)AttacksAvailable()(Buffer []geometry.Point){
	return n.ListStepsAvailable()
}

func (n NonLinearFigure)CheckStepAvailable(point geometry.Point)(bool){
	return n.checkAvailable(n.ListStepsAvailable(),point)
}

func (n NonLinearFigure)CheckAttackAvailable(point geometry.Point)(bool){
	return n.CheckStepAvailable(point)
}

func (n NonLinearFigure)CheckForCollision(destination, obstacle geometry.Point)(bool){
	return destination.Equal(obstacle)
}

// list of child classes representing figures, classes named as figures in english so
//I've added some russian cursive comments so i'll understand myself
type King struct{
	NonLinearFigure
}

func ConstructKing(x,y int, colour Colour)(*King){
	return &King{NonLinearFigure:NonLinearFigure{Figure:ConstructFigure(x,y,colour),ProbableSteps:KingProbableStepList}}
}

type Queen struct{
	LinearFigure
}

func ConstructQueen(x,y int, colour Colour)(*Queen){
	direction:=make([]geometry.Vector,0,4)
	direction = append(direction,geometry.Vector{Point:geometry.ConstructPoint(0,1)})
	direction = append(direction,geometry.Vector{Point:geometry.ConstructPoint(1,0)})
	direction = append(direction,geometry.Vector{Point:geometry.ConstructPoint(1,1)})
	direction = append(direction,geometry.Vector{Point:geometry.ConstructPoint(1,-1)})
	return &Queen{LinearFigure:LinearFigure{Figure:ConstructFigure(x,y,colour),Direction:direction}}
}


type Bishop struct{
	LinearFigure	//	slon
}

func ConstructBishop(x,y int ,colour Colour)(*Bishop){
	direction:=make([]geometry.Vector,0,2)
	direction = append(direction,geometry.Vector{Point:geometry.ConstructPoint(0,1)})
	direction = append(direction,geometry.Vector{Point:geometry.ConstructPoint(1,0)})
	return &Bishop{LinearFigure{ConstructFigure(x,y,colour),direction}}
}

type Knight struct{
	NonLinearFigure	//	kon
}

func ConstructKnight(x,y int ,colour Colour)(*Knight){
	return &Knight{NonLinearFigure{ConstructFigure(x,y,colour),KnightProbableStepList}}
}


type Rook struct{
	LinearFigure	//	ladya
}

func ConstructRook(x,y int ,colour Colour)(*Rook){
	direction:=make([]geometry.Vector,0,2)
	direction = append(direction,geometry.Vector{Point:geometry.ConstructPoint(1,1)})
	direction = append(direction,geometry.Vector{Point:geometry.ConstructPoint(1,-1)})
	return &Rook{LinearFigure{ConstructFigure(x,y,colour), direction}}
}

type Pawn struct{
	Figure
	didStep bool
}

func ConstructPawn(x,y int ,colour Colour)(*Pawn){
	return &Pawn{ConstructFigure(x,y,colour),false}
}

func(p Pawn)CheckStepAvailable(point geometry.Point)(bool){
	for _,element := range p.ListStepsAvailable(){
		if element.Equal(point){
			return true
		}
	}
	return false
}

func(p Pawn)CheckAttackAvailable(point geometry.Point)(bool){
	for _,element := range p.AttacksAvailable(){
		if element.Equal(point){
			return true
		}
	}
	return false
}

func(p Pawn)ListStepsAvailable()(Buffer []geometry.Point){
	Buffer = make([]geometry.Point,0,2)
	if element:=PawnProbableShortStepList[p.Colour_];element.Add(p.Point).CheckFieldBoundaries(){
		Buffer = append(Buffer, element.Add(p.Point))
	}
	if element:=PawnProbableLongStepList[p.Colour_];!p.didStep && element.Add(p.Point).CheckFieldBoundaries(){
		Buffer = append(Buffer, element.Add(p.Point))
	}
	return Buffer
}

func (p Pawn) AttacksAvailable()(Buffer []geometry.Point){
	Buffer = make([]geometry.Point,0,2)
	for _,element:=range PawnProbableAttackList[p.Colour_]{
		if element.Add(p.Point).CheckFieldBoundaries(){
			Buffer = append(Buffer, element.Add(p.Point))
		}
	}
	return Buffer
}

func (p Pawn)CheckForCollision(destination,obstacle geometry.Point)(bool){
	way:=geometry.ConstructLine(p.Point, destination)
	return way.Intersect(obstacle)
}

type Colour uint8
const (
	BLACK = iota
	WHITE
)

type GameSession struct{
	Figures []StepMaker
	AuthToken string
	PlayingNow Colour
	StepDone bool
}

func InitFromString(s string)(gs GameSession,ok bool){
	figures, token, playingNow, stepDone,ok := parser.GenerateGameSessionData(s)
	if !ok{
		return gs,false
	}
	gs.Figures = figures
	gs.AuthToken = token
	gs.PlayingNow = playingNow
	gs.StepDone = stepDone
	return gs, true
}


//function, checking if there is
func (g GameSession)At(position geometry.Point)(StepMaker,bool){
	for _,element := range g.Figures{
		if element.RetCoords().Equal(position) && g.PlayingNow == element.RetColour(){
			return element,true
		}
	}
	return nil,false
}

func (g GameSession)CanGo(destination geometry.Point, fig  StepMaker)(bool){
	return fig.CheckStepAvailable(destination)
}

func (g GameSession)CheckStepForCollisions(destination geometry.Point, fig StepMaker)(collidedIndex  int, collide bool){
	for i,element := range g.Figures{
		if element.RetCoords().Equal(fig.RetCoords()){	//	ignore collision with self
		}else if fig.CheckForCollision(destination,element.RetCoords()){
			return i,true
		}
	}
	return -1,false
}

func (g GameSession)CanAttack(destination geometry.Point, fig  StepMaker)(bool){
	return fig.CheckAttackAvailable(destination)
}
