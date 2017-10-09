package model

import "chessServer/utility/geometry"

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

func (f * Figure)Step(p geometry.Point){
	f.Point = p
}