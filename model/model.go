package model

import (
	"chessServer/utility/geometry"
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

type Colour int
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
	PlayingNow Colour	// which players turn
	StepDone bool		// is someone have finished step, or just picked a figure
	chosenFigure StepMaker	//	if someone has picked figure
}

type FigJsonRepr struct{
	Name string `json:"name"`
	X   int `json:"x"`
	Y   int `json:"y"`
	Col int `json:"colour"`
}

type GameSessionJsonRepr struct{
	Figs []FigJsonRepr `json:"figs"`
	GameOver int 		`json:"game_over"` // 0 for normal game, 1 for black, 2 for white, 3 for draw
	ProbSteps []geometry.Point `json:"list_steps"`
	Player int `json:"player"`	//	player, for clients to know what colour does he has << is that a fucking psycho pass reference??
}



func (g GameSession)ToJsonRepr()[]FigJsonRepr{
	tmp := make([]FigJsonRepr,0,32)
	for _,element := range g.Figures{
		p := element.RetCoords()
		figType := ""
		switch element.(type) {
		case *King:
			figType = "king"
		case *Queen:
			figType = "queen"
		case *Rook:
			figType = "rook"
		case *Knight:
			figType = "knight"
		case *Pawn:
			figType = "pawn"
		case *Bishop:
			figType = "bishop"
		}
		tmp = append(tmp, FigJsonRepr{figType, p.X,p.Y, int(element.RetColour())})
	}
	return tmp
}

func New()(GameSession){
	g := GameSession{PlayingNow:WHITE,chosenFigure:nil,Figures:make([]StepMaker,0,32)}

	g.Figures = append(g.Figures, ConstructKing(4,7,WHITE))
	g.Figures = append(g.Figures, ConstructKing(4,0,BLACK))
	g.Figures = append(g.Figures, ConstructQueen(3,7,WHITE))
	g.Figures = append(g.Figures, ConstructQueen(3,0,BLACK))
	g.Figures = append(g.Figures, ConstructBishop(0,7,WHITE))
	g.Figures = append(g.Figures, ConstructBishop(7,7,WHITE))
	g.Figures = append(g.Figures, ConstructBishop(0,0,BLACK))
	g.Figures = append(g.Figures, ConstructBishop(7,0,BLACK))
	g.Figures = append(g.Figures, ConstructRook(2,7,WHITE))
	g.Figures = append(g.Figures, ConstructRook(5,7,WHITE))
	g.Figures = append(g.Figures, ConstructRook(2,0,BLACK))
	g.Figures = append(g.Figures, ConstructRook(5,0,BLACK))
	g.Figures = append(g.Figures, ConstructKnight(1,7,WHITE))
	g.Figures = append(g.Figures, ConstructKnight(6,7,WHITE))
	g.Figures = append(g.Figures, ConstructKnight(1,0,BLACK))
	g.Figures = append(g.Figures, ConstructKnight(6,0,BLACK))

	for i:=0; i < 8;i++{
		g.Figures = append(g.Figures, ConstructPawn(i,6,WHITE))
		g.Figures = append(g.Figures, ConstructPawn(i,1,BLACK))
	}
	return g
}

func (g GameSession)InitialToJsonRepr()(GameSessionJsonRepr){
	jsonRepr := GameSessionJsonRepr{GameOver:0, ProbSteps:make([]geometry.Point,0,0)}
	jsonRepr.Figs = g.ToJsonRepr()
	return jsonRepr
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
	case *King:
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

func(g * GameSession)CheckGameOver()(bool){
	for _,fig := range g.Figures{
		if fig.RetColour() == g.PlayingNow{
			for _,step := range fig.AttacksAvailable(){
				if ok,_ := g.CanAct(step,fig);ok{
					return false
				}
			}
			for _,step := range fig.ListStepsAvailable(){
				if ok,_ := g.CanAct(step,fig);ok{
					return false
				}
			}
		}
	}
	return true
}

//Simple check if king is under attack or not
func (g GameSession)CheckForCheck()bool{
	var (
		flag bool = false
		king *King
	)
	for _,element:= range g.Figures{
		switch t := element.(type) {
		case *King:
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


func (g * GameSession)Act(clicked geometry.Point)GameSessionJsonRepr{
	var repr = GameSessionJsonRepr{}
	if g.chosenFigure != nil{
		if can, collision := g.CanAct(clicked, g.chosenFigure);can{
			if collision >= 0{
				g.Figures = append(g.Figures[:collision], g.Figures[collision+1:]...)
			}
			g.chosenFigure.SetCoords(clicked)
			g.chosenFigure = nil
		}
		if g.CheckGameOver(){
			if g.CheckForCheck(){
				repr.GameOver = 1 + int(g.PlayingNow)
			}
		}
		g.PlayingNow = (g.PlayingNow + 1) & 1
		if g.CheckGameOver(){
			if g.CheckForCheck(){
				repr.GameOver = 1 + int(g.PlayingNow)
			}else{
				repr.GameOver = 3
			}
		}else{
			repr.GameOver = 0
		}
	}else if fig,ok := g.At(clicked);ok{
		g.chosenFigure = fig
		repr.GameOver = 0
		tmpProbSteps := fig.ListStepsAvailable()
		tmpProbSteps = append(repr.ProbSteps,fig.AttacksAvailable()...)
		repr.ProbSteps = make([]geometry.Point,0,32)
		for _,element := range tmpProbSteps{
			if ok, _ :=g.CanAct(element,fig);ok{
				repr.ProbSteps = append(repr.ProbSteps, element)
			}
		}
	}
	repr.Figs = g.ToJsonRepr()
	return repr
}


