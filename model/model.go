package model

import (
	"chessServer/utility"
)

var(
	KingProbableStepList = [8]utility.Point{utility.Point{-1,-1},utility.Point{-1,0},utility.Point{-1,1},
															 utility.Point{0, -1},utility.Point{0, 1},utility.Point{1,-1},
															 utility.Point{1,0},utility.Point{1,1},}
	KnightProbableStepList = [8]utility.Point{utility.Point{-2,-1},utility.Point{-1,-2},utility.Point{1,-2},
											  utility.Point{2, -1},utility.Point{-2, 1},utility.Point{-1,2},
											  utility.Point{1,2},utility.Point{2,1},}
)
// TODO add classes Straight figure with collisions and figure with collisions in place only
//	interface for all figures, just checking if figure can go somewhere/attack a field
type StepMaker interface{
	CheckStepAvailable(point utility.Point)(bool)	//	check if figure can go to the following field
	CheckAttackAvailable(point utility.Point)(bool)	//	check if figure can attack the following field

	StepsAvailable()([]utility.Point)	//	list all available fields to go
	AttacksAvailable()([]utility.Point)	//	list all available fields to attack

	CheckForCollision(destination,obstacle utility.Point)(bool)	//	check if we can collide with other figure while doing step
}



type Figure struct{	//	parent class for all figures(all figures inherits Figure and implement StepMaker)
	utility.Point	//	figure coordinates
	Colour_ Colour
}

func ConstructFigure(x,y int ,colour Colour)(Figure) {
	return Figure{utility.Point{x,y},colour}
}

func (f Figure)checkAvailable(available []utility.Point, point utility.Point)(bool){	//	utility function to check whereas given point is available for step
	for _,element:=range available{
		if element.Equal(point){
			return true
		}
	}
	return false
}

type LinearFigure struct{	//	bishops and rooks
	Figure
	Direction [2]utility.Vector
}

func (l LinearFigure)StepsAvailable()(Buffer []utility.Point){
	Buffer = make([]utility.Point, 16)
	for _,element:=range l.Direction{
		i,j:=l.Point.Add(element.Point),l.Point.Substract(element.Point)
		for i.CheckFieldBoundaries() && j.CheckFieldBoundaries(){
			Buffer = append(Buffer,i)
			Buffer = append(Buffer,i)
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
	way:=utility.Line{l.Point,destination}	// only valid destinations are checked for collision
	return way.Intersect(obstacle)						// so we can skip checking for validity
}

type NonLinearFigure struct{	//	Kings/Knights, figures, which can collide only if they rich collision place
	Figure
	ProbableSteps [8]utility.Point	//	both kings and knights can visit only 8 places
}

func (n NonLinearFigure)StepsAvailable()(Buffer []utility.Point){
	Buffer = make([]utility.Point, 8)
	for i,element:= range n.ProbableSteps{
		if n.Point.Add(element).CheckFieldBoundaries(){
			Buffer[i]=n.Point.Add(element)
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
	return King{NonLinearFigure{ConstructFigure(x,y,colour),KingProbableStepList}}
}

type Queen struct{
	Figure
}

type Bishop struct{
	LinearFigure	//	slon
}

func ConstructBishop(x,y int ,colour Colour)(Bishop){
	direction:=[2]utility.Vector{utility.Vector{utility.Point{0,1}},utility.Vector{utility.Point{1,0}}}
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
	direction:=[2]utility.Vector{utility.Vector{utility.Point{X:1,Y:1}},utility.Vector{utility.Point{X: 1, Y: -1}}}
	return Rook{LinearFigure{ConstructFigure(x,y,colour), direction}}
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
}





