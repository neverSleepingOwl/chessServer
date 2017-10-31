package model

import (
	"github.com/chessServer/utility/geometry"
	"github.com/chessServer/utility/logger"
)

// NonlinearFigure represents figures that can go and attack to the same cells:Kings and Knights, relative to it's coordinates,
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

	logger.WriteLog(5, "From line 30: ", "Nonlinear figure available steps: ", Buffer)
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
