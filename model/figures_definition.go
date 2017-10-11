package model

import (
	"chessServer/utility/geometry"
	"chessServer/utility/logger"
)

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

	logger.WriteLog(5, "From line 110: ", "Pawn step list: ", Buffer)
	return Buffer
}

func (p Pawn) AttacksAvailable()(Buffer []geometry.Point){
	Buffer = make([]geometry.Point,0,2)
	for _,element:=range PawnProbableAttackList[p.Colour_]{
		if element.Add(p.Point).CheckFieldBoundaries(){
			Buffer = append(Buffer, element.Add(p.Point))
		}
	}

	logger.WriteLog(5, "From line 120: ", "Pawn attack list: ", Buffer)
	return Buffer
}

func (p Pawn)CheckForCollision(destination,obstacle geometry.Point)(bool){
	logger.WriteLog(5, "From line 123: ", "Pawn collision check Destination: ", destination, "Obstacle: ", obstacle)
	way:=geometry.ConstructLine(p.Point, destination)
	return way.Intersect(obstacle)
}

func (p * Pawn)Step(point geometry.Point){
	p.Figure.Step(point)
	p.didStep = true
}