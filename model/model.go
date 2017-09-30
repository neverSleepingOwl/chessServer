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
	switch g.Figures[collision].(type) {	//	TODO Probably causes error
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


