package model

import (
	"chessServer/utility"
	"debug/plan9obj"
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


// list of child classes representing figures, classes named as figures in english so
//I've added some russian cursive comments so i'll understand myself
type King struct{
	Figure
}

func (k King)StepsAvailable()(Buffer []utility.Point){
	Buffer = make([]utility.Point, 8)
	for i:=k.X-1;i <= k.X+1;i++{
		for j:=k.Y-1;j <= k.Y+1;j++{
			temp:=utility.Point{i,j}
			if !temp.Equal(k.Point){
				if temp.CheckFieldBoundaries(){
					Buffer = append(Buffer, temp)
				}
			}
		}
	}
	return Buffer
}

func (k King)AttacksAvailable()(Buffer []utility.Point){
	Buffer = k.StepsAvailable()
	return Buffer
}

func (k King)CheckStepAvailable(point utility.Point)(bool){
	stepsAvailable:=k.StepsAvailable()
	for _,element:=range stepsAvailable{
		if element.Equal(point){
			return true
		}
	}
	return false
}

func (k King)CheckAttackAvailable(point utility.Point)(bool){
	return k.CheckStepAvailable(point)
}

func (k King)CheckForCollision(destination, obstacle utility.Point)(bool){
	return destination.Equal(obstacle)
}
type Queen struct{
	Figure
}

type Bishop struct{
	Figure	//	slon
}

func(b Bishop)StepsAvailable()(Buffer []utility.Point){
	Buffer = make([]utility.Point, 16)
	for i:=0;i<8;i++{
		tmp1,tmp2:=utility.Point{b.X,i},utility.Point{i,b.Y}
		if !tmp1.Equal(b.Point){
			Buffer = append(Buffer,tmp1)
		}
		if !tmp2.Equal(b.Point){
			Buffer = append(Buffer,tmp2)
		}
	}
	return Buffer
}

func (b Bishop)AttacksAvailable()(Buffer []utility.Point){
	return b.StepsAvailable()
}

func (b Bishop)CheckStepAvailable(point utility.Point)(bool){
	for _,element:=range b.StepsAvailable(){
		if b.Point.Equal(element){
			return true
		}
	}
	return false
}

func (b Bishop)CheckAttackAvailable(point utility.Point)(bool){
	return b.CheckStepAvailable(point)
}

func (b Bishop)CheckForCollision(destination, obstacle utility.Point)(bool){
	if b.CheckStepAvailable(destination){
		way:=utility.Line{b.Point, destination}
		return way.Intersect(obstacle)
	}else{
		return false
	}
}

type Knight struct{
	Figure	//	kon
}

func(k Knight)StepsAvailable()(Buffer []utility.Point){
	storage:=[]utility.Point{Point{},Point{}}
	Buffer = make([]utility.Point,8)
	return Buffer
}

func (k Knight)AttacksAvailable()(Buffer []utility.Point){
	return b.StepsAvailable()
}

func (k Knight)CheckStepAvailable(point utility.Point)(bool){
	for _,element:=range b.StepsAvailable(){
		if b.Point.Equal(element){
			return true
		}
	}
	return false
}

func (k Knight)CheckAttackAvailable(point utility.Point)(bool){
	return b.CheckStepAvailable(point)
}

func (k Knight)CheckForCollision(destination, obstacle utility.Point)(bool){
	if b.CheckStepAvailable(destination){
		way:=utility.Line{b.Point, destination}
		return way.Intersect(obstacle)
	}else{
		return false
	}
}


type Rook struct{
	Figure	//	ladya
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



