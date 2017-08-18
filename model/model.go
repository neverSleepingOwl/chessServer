package model

import (
	"chessServer/utility"
)


var(
	KingProbableStepList = [8]utility.Point{{-1,-1},{-1,0},{-1,1},
											{0, -1},{0, 1},{1,-1},
											{1,0},{1,1},}
	KnightProbableStepList = [8]utility.Point{{-2,-1},{-1, -2},{1, -2},
											  {2, -1},{-2, 1},{-1,2},
											  {1, 2},{2, 1},}
	PawnProbableShortStepList = [2]utility.Point{{0,1},	//	probable short steps for black pawns
											   {0,-1},	//	probable short steps for white pawns
											  }
	PawnProbableLongStepList = [2]utility.Point{{0,2},	//	probable steps for black pawns
												{0,-2},	//	probable steps for white pawns
												}
	PawnProbableAttackList = [2][2]utility.Point{
		{{1,1}, {-1,1}},	//	probable steps for black pawns
		{{1,-1}, {-1,-1}},	//	probable steps for white pawns
	}
)
// TODO add classes Straight figure with collisions and figure with collisions in place only
//	interface for all figures, just checking if figure can go somewhere/attack a field
type StepMaker interface{
	CheckStepAvailable(point utility.Point)(bool)	//	check if figure can go to the following field
	CheckAttackAvailable(point utility.Point)(bool)	//	check if figure can attack the following field

	StepsAvailable()([]utility.Point)	//	list all available fields to go
	AttacksAvailable()([]utility.Point)	//	list all available fields to attack

	CheckForCollision(destination,obstacle utility.Point)(bool)	//	check if we can collide with other figure while doing step

	RetCoords()utility.Point

	RetColour()Colour

	isEnemy(maker StepMaker)bool
}

type array struct{
	arr [32]utility.Point
	size uint8
}

type Figure struct{	//	parent class for all figures(all figures inherits Figure and implement StepMaker)
	utility.Point	//	figure coordinates
	Colour_ Colour
	pSteps array
	pAttacks array
}

func(f Figure)RetColour()Colour{
	return f.Colour_
}

func (f Figure)isEnemy(maker StepMaker)bool{	//	true if two figures are enemies
	return f.Colour_ != maker.RetColour()
}

func ConstructFigure(x,y int ,colour Colour)(Figure) {
	return Figure{Point: utility.ConstructPoint(x,y),Colour_: colour}
}

func (f Figure)checkAvailable(available []utility.Point, point utility.Point)(bool){	//	utility function to check whereas given point is available for step
	for _,element:=range available{
		if element.Equal(point){
			return true
		}
	}
	return false
}

func (f Figure)RetCoords()utility.Point{
	return f.Point
}

type LinearFigure struct{	//	bishops and rooks
	Figure
	Direction []utility.Vector
}

func (l LinearFigure)StepsAvailable()(Buffer []utility.Point){
	Buffer = make([]utility.Point,0,16)
	for _,element:=range l.Direction{
		i,j:=l.Point.Add(element.Point),l.Point.Subtract(element.Point)
		for i.CheckFieldBoundaries(){
			Buffer = append(Buffer,i)
			i = i.Add(element.Point)
		}
		for j.CheckFieldBoundaries(){
			Buffer = append(Buffer,j)
			j = j.Subtract(element.Point)
		}
	}
	return Buffer
}

func (l LinearFigure)AttacksAvailable()(Buffer []utility.Point){
	return l.StepsAvailable()
}

func (l LinearFigure)CheckStepsAvailable(point utility.Point)(bool){
	return l.checkAvailable(l.StepsAvailable(),point)
}

func (l LinearFigure)CheckAttacksAvailable(point utility.Point)(bool){
	return l.CheckStepsAvailable(point)
}

func (l LinearFigure)CheckForCollision(destination, obstacle utility.Point)(bool){
	way:=utility.ConstructLine(l.Point, destination) // only valid destinations are checked for collision
	return way.Intersect(obstacle)                      // so we can skip checking for validity
}

type NonLinearFigure struct{	//	Kings/Knights, figures, which can collide only if they rich collision place
	Figure
	ProbableSteps [8]utility.Point	//	both kings and knights can visit only 8 places
}

func (n NonLinearFigure)StepsAvailable()(Buffer []utility.Point){
	Buffer = make([]utility.Point,0,8)
	for _,element:= range n.ProbableSteps{
		if n.Point.Add(element).CheckFieldBoundaries(){
			Buffer=append(Buffer,n.Point.Add(element))
		}
	}
	return Buffer
}

func (n NonLinearFigure)AttacksAvailable()(Buffer []utility.Point){
	return n.StepsAvailable()
}

func (n NonLinearFigure)CheckStepsAvailable(point utility.Point)(bool){
	return n.checkAvailable(n.StepsAvailable(),point)
}

func (n NonLinearFigure)CheckAttacksAvailable(point utility.Point)(bool){
	return n.CheckStepsAvailable(point)
}

func (n NonLinearFigure)CheckForCollision(destination, obstacle utility.Point)(bool){
	return destination.Equal(obstacle)
}

// list of child classes representing figures, classes named as figures in english so
//I've added some russian cursive comments so i'll understand myself
type King struct{
	NonLinearFigure
}

func ConstructKing(x,y int, colour Colour)(King){
	return King{NonLinearFigure:NonLinearFigure{Figure:ConstructFigure(x,y,colour),ProbableSteps:KingProbableStepList}}
}

type Queen struct{
	LinearFigure
}

func ConstructQueen(x,y int, colour Colour)(Queen){
	direction:=make([]utility.Vector,0,4)
	direction = append(direction,utility.Vector{Point:utility.ConstructPoint(0,1)})
	direction = append(direction,utility.Vector{Point:utility.ConstructPoint(1,0)})
	direction = append(direction,utility.Vector{Point:utility.ConstructPoint(1,1)})
	direction = append(direction,utility.Vector{Point:utility.ConstructPoint(1,-1)})
	return Queen{LinearFigure:LinearFigure{Figure:ConstructFigure(x,y,colour),Direction:direction}}
}


type Bishop struct{
	LinearFigure	//	slon
}

func ConstructBishop(x,y int ,colour Colour)(Bishop){
	direction:=make([]utility.Vector,0,2)
	direction = append(direction,utility.Vector{Point:utility.ConstructPoint(0,1)})
	direction = append(direction,utility.Vector{Point:utility.ConstructPoint(1,0)})
	return Bishop{LinearFigure{ConstructFigure(x,y,colour),direction}}
}

type Knight struct{
	NonLinearFigure	//	kon
}

func ConstructKnight(x,y int ,colour Colour)(Knight){
	return Knight{NonLinearFigure{ConstructFigure(x,y,colour),KnightProbableStepList}}
}


type Rook struct{
	LinearFigure	//	ladya
}

func ConstructRook(x,y int ,colour Colour)(Rook){
	direction:=make([]utility.Vector,0,2)
	direction = append(direction,utility.Vector{Point:utility.ConstructPoint(1,1)})
	direction = append(direction,utility.Vector{Point:utility.ConstructPoint(1,-1)})
	return Rook{LinearFigure{ConstructFigure(x,y,colour), direction}}
}

type Pawn struct{
	Figure
	didStep bool
}

func ConstructPawn(x,y int ,colour Colour)(Pawn){
	return Pawn{ConstructFigure(x,y,colour),false}
}

func(p Pawn)CheckStepAvailable(point utility.Point)(bool){
	for _,element := range p.StepsAvailable(){
		if element.Equal(point){
			return true
		}
	}
	return false
}

func(p Pawn)CheckAttackAvailable(point utility.Point)(bool){
	for _,element := range p.AttacksAvailable(){
		if element.Equal(point){
			return true
		}
	}
	return false
}

func(p Pawn)StepsAvailable()(Buffer []utility.Point){
	Buffer = make([]utility.Point,0,2)
	if element:=PawnProbableShortStepList[p.Colour_];element.Add(p.Point).CheckFieldBoundaries(){
		Buffer = append(Buffer, element.Add(p.Point))
	}
	if element:=PawnProbableLongStepList[p.Colour_];!p.didStep && element.Add(p.Point).CheckFieldBoundaries(){
		Buffer = append(Buffer, element.Add(p.Point))
	}
	return Buffer
}

func (p Pawn) AttacksAvailable()(Buffer []utility.Point){
	Buffer = make([]utility.Point,0,2)
	for _,element:=range PawnProbableAttackList[p.Colour_]{
		if element.Add(p.Point).CheckFieldBoundaries(){
			Buffer = append(Buffer, element.Add(p.Point))
		}
	}
	return Buffer
}

func (p Pawn)CheckForCollision(destination,obstacle utility.Point)(bool){
	way:=utility.ConstructLine(p.Point, destination)
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
	Finished Colour
	StepDone bool
}

func (g GameSession)At(position utility.Point)(StepMaker,bool){
	for _,element := range g.Figures{
		if element.RetCoords().Equal(position){
			return element,true
		}
	}
	return nil,false
}

func (g GameSession)CanGo(destination utility.Point, fig  StepMaker)(bool){
	return fig.CheckStepAvailable(destination)
}

func (g GameSession)CheckStepForCollisions(destination utility.Point, fig StepMaker)(collidedIndex  int, collide bool){
	for i,element := range g.Figures{
		if element.RetCoords().Equal(fig.RetCoords()){	//	ignore collision with self
		}else if fig.CheckForCollision(destination,element.RetCoords()){
			return i,true
		}
	}
	return -1,false
}

func (g GameSession)CanAttack(destination utility.Point, fig  StepMaker)(bool){
	return fig.CheckAttackAvailable(destination)
}
