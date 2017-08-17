package model

//	interface for all figures, just checking if figure can go somewhere/attack a field
type StepMaker interface{
	CheckStepAvailable(point Point)(bool)	//	check if figure can go to the following field
	CheckAttackAvailable(point Point)(bool)	//	check if figure can attack the following field
	StepsAvailable()([]Point, bool)	//	list all available fields to go
	AttacksAvailable()([]Point, bool)	//	list all available fields to attack
}

type Point struct{
	X uint
	Y uint
}

type Figure struct{	//	parent class for all figures(all figures inherits Figure and implement StepMaker)
	Point	//	figure coordinates
	Colour_ Colour
}


// list of child classes representing figures, classes named as figures in english so
//I've added some russian cursive comments so i'll understand myself
type King struct{
	Figure
}

type Queen struct{
	Figure
}

type Bishop struct{
	Figure	//	slon
}

type Knight struct{
	Figure	//	kon
}

type R


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



