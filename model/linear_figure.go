package model

import (
	"chessServer/utility/geometry"
	"chessServer/utility/logger"
)

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
	logger.WriteLog(5, "From line 34: All linear figure available steps: ", Buffer)
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
	logger.WriteLog(5, "Linear figure collision check. From line 54: ", "Destination: ", destination, "Obstacle: ", obstacle)
	return way.Intersect(obstacle)                      // so we can skip checking for validity
}