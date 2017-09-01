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

	SetCoords(p geometry.Point)

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

func (f *Figure)SetCoords(p geometry.Point){
	f.Point = p
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
	Buffer = make([]geometry.Point,0,16)	//	allocate memory for output

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
	Figure	//	coordinates and colour, base class
	ProbableSteps [8]geometry.Point	//	both kings and knights can visit only 8 places
}

//To evaluate available steps we just add to all template coordinates coordinates of figure itself
// and check if result fits field boundaries
func (n NonLinearFigure)ListStepsAvailable()(Buffer []geometry.Point){
	Buffer = make([]geometry.Point,0,8)	//	allocate memory for list of available coordinates

	for _,element:= range n.ProbableSteps{	//	for each template step coordinate
		if n.Point.Add(element).CheckFieldBoundaries(){	//	add coordinate of figure itself and check bounds
			Buffer=append(Buffer,n.Point.Add(element))	//	add correct coordinate to list
		}
	}
	return Buffer
}

// NonLinear Figure attack places, where they can make step
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
//All classes except Pawns inherit Linear or NonLinear figures, so, I'Ve Just needed to implement some constructors
// All constructors return pointer to an object in order to assign result to interface
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
	direction:=make([]geometry.Vector,0,4)	//	allocate vector of directions

	direction = append(direction,geometry.Vector{Point:geometry.ConstructPoint(0,1)})	//	Queen behaves like Bishop and Rook same time
	direction = append(direction,geometry.Vector{Point:geometry.ConstructPoint(1,0)})	// fill vector of directions
	direction = append(direction,geometry.Vector{Point:geometry.ConstructPoint(1,1)})
	direction = append(direction,geometry.Vector{Point:geometry.ConstructPoint(1,-1)})

	return &Queen{LinearFigure:LinearFigure{Figure:ConstructFigure(x,y,colour),Direction:direction}}
}


type Bishop struct{
	LinearFigure	//	slon
}

func ConstructBishop(x,y int ,colour Colour)(*Bishop){
	direction:=make([]geometry.Vector,0,2)	//	allocate vector of directions

	direction = append(direction,geometry.Vector{Point:geometry.ConstructPoint(0,1)})
	direction = append(direction,geometry.Vector{Point:geometry.ConstructPoint(1,0)})	// fill vector of directions

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


//Pawn differs from both Linear and Nonlinear figures

type Pawn struct{
	Figure
	didStep bool
}


func ConstructPawn(x,y int ,colour Colour)(*Pawn){	// constructor
	return &Pawn{ConstructFigure(x,y,colour),false}
}

//check if pawn can step into the following cell
func(p Pawn)CheckStepAvailable(point geometry.Point)(bool){
	for _,element := range p.ListStepsAvailable(){	//	just iterate through all available positions
		if element.Equal(point){
			return true
		}
	}
	return false
}

func(p Pawn)CheckAttackAvailable(point geometry.Point)(bool){	//	the same for attacks
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

//Class, representing a single game session
//Does the following actions:
//load from database, then performs step if player has already picked a figure
//else if player picks figure returns success or failure if there is no figure on given coordinate
type GameSession struct{
	Figures []StepMaker
	AuthToken string	//actually it's a number of session
	Password string		// password to authenticate
	PlayingNow Colour	// which players turn
	StepDone bool		// is someone have finished step, or just picked a figure
	chosenFigure StepMaker	//	if someone has picked figure
	higlightedFigs []geometry.Point	//	when you pick a figure you should know
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


//function, checking if there is a figure in given position
func (g GameSession)At(position geometry.Point)(StepMaker,bool){
	for _,element := range g.Figures{
		if element.RetCoords().Equal(position) && g.PlayingNow == element.RetColour(){
			return element,true
		}
	}
	return nil,false
}


// Check if figure collides with other figure while making step
//return index of collided figure if collision happens
func (g GameSession)CheckStepForCollisions(destination geometry.Point, fig StepMaker)(collidedIndex  int, collide bool){
	for i,element := range g.Figures{
		if element.RetCoords().Equal(fig.RetCoords()){	//	ignore collision with self
		}else if fig.CheckForCollision(destination,element.RetCoords()){
			return i,true
		}
	}
	return -1,false
}


// Recognises, whereas we can perform the attack/step to a given position
// performs standart step/attack action with collision check, removes attacked figure
// in case of attack, but saves it to a temporary variable
// if action doesn't case check to attacking side, then return true and index of attacked figure
// if no figures has been attacked returns negative index
// returns result of attack (true if success) and index of attacked figure
//WARNING, can't attack king
//TODO add castling
func (g * GameSession)CanAct(destination geometry.Point, fig StepMaker)(bool, int){
	var (
		temporaryDeletedFig StepMaker
		prevCoords geometry.Point
		deleted bool = false	//	flag to measure if we eated ( deleted figure from main array
	)
	//find out if we can step to a given position
	step:= fig.CheckStepAvailable(destination)
	//find out if we collide with an obstacle while performing step
	collision,yes:=g.CheckStepForCollisions(destination,fig)
	//we can attack if our destination has the same coordinates
	//as the figure we collide, we can attack to a given destination and
	//figure, placed at destination has different colour
	attackAble:=fig.CheckAttackAvailable(destination) && yes &&
			g.Figures[collision].isEnemy(fig) && g.Figures[collision].RetCoords().Equal(destination)
	switch t:=g.Figures[collision].(type) {
	//Can't attack king
	case King:
		attackAble = false
	default:	//	to prevent compile error, we don't need t variable

	}

	switch{
	case attackAble:
		temporaryDeletedFig = g.Figures[collision] // perform attack virtually
		g.Figures = append(g.Figures[:collision], g.Figures[collision+1:]...)	//	delete element from main array
		prevCoords = fig.RetCoords()	//	remember previous coordinates
		deleted = true
		fig.SetCoords(destination)
	case step && !yes:
		prevCoords = fig.RetCoords()
		fig.SetCoords(destination)//perform step virtually
	default:
		return false,-1
	}

	//If check occurs abandon
	if g.CheckForCheck(){
		fig.SetCoords(prevCoords)	//	undo step
		if deleted{
			g.Figures = append(g.Figures,temporaryDeletedFig)
		}
		return false,-1
	}else{
		fig.SetCoords(prevCoords)
		ret:=-1
		if deleted{
			g.Figures = append(g.Figures,temporaryDeletedFig)
			ret = collision
		}
		return true,ret
	}
}

//Simple check if king is under attack or not
func (g GameSession)CheckForCheck()bool{
	var (
		flag bool = false
		king King
	)
	for _,element:= range g.Figures{
		switch t := element.(type) {
		case King:
			if t.Colour_ == g.PlayingNow{
				king = t
				flag = true
				break
			}
		}
	}
	if flag{
		for _,element:=range g.Figures{	//	search if any figure attacks king
			if element.CheckAttackAvailable(king.Point){
				if _,collides:=g.CheckStepForCollisions(king.Point, element);!collides{
					return true	//	if king is under attack
				}
			}
		}
		return false
	}else{
		return false//this can't happen, TODO generate exception or sth
	}
}

func (g GameSession)Act(clicked geometry.Point){

}


