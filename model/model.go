package model

import (
	"chessServer/utility/geometry"
	"log"
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

	Step(p geometry.Point)
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
	Figures      []StepMaker
	PlayingNow   Colour	// which players turn
	StepDone     bool		// is someone have finished step, or just picked a figure
	chosenFigure StepMaker	//	if someone has picked figure
	kings        []*King
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
	kw := ConstructKing(4,7,WHITE)
	kb := ConstructKing(4,0,BLACK)
	g.Figures = append(g.Figures, kw)
	g.Figures = append(g.Figures, kb)
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
	g.kings = []*King{kb,kw}
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
			log.Println("Clicked figure at: ",element.RetCoords(), " ", element.RetColour())
			return element,true
		}
	}
	return nil,false
}


// Check if figure collides with other figure while making step
//return index of collided figure if collision happens
func (g GameSession) FindAllCollisions(destination geometry.Point, fig StepMaker)([]int){
	buffer := make([]int,0,16)

	for i,element := range g.Figures{
		if element.RetCoords().Equal(fig.RetCoords()){
			//	ignore collision with self
		}else if fig.CheckForCollision(destination,element.RetCoords()){
			buffer = append(buffer, i)
		}
	}
	return buffer
}


// Recognises, whereas we can perform the attack/step to a given position
// performs standart step/attack action with collision check, removes attacked figure
// in case of attack, but saves it to a temporary variable
// if action doesn't case check to attacking side, then return true and index of attacked figure
// if no figures has been attacked returns negative index
// returns result of attack (true if success) and index of attacked figure
//WARNING, can't attack kings
//TODO add castling
func (g * GameSession)CanAct(destination geometry.Point, fig StepMaker)(bool, int){
	var (
		temporaryDeletedFig StepMaker
		prevCoords geometry.Point
		deleted bool = false	//	flag to measure if we eated ( deleted figure from main array
		output bool = true
		num = -1
	)
	if g.CanStepWithNoCollisions(destination,fig){
		prevCoords = g.StepVirtually(destination,fig)
	}else if attacked,can := g.CanAttack(destination, fig);can{
		deleted = true
		prevCoords  = g.StepVirtually(destination,fig)
		temporaryDeletedFig = g.Figures[attacked]
		g.Figures = append(g.Figures[:attacked], g.Figures[:attacked+1]...)
	}else{
		return false, -1
	}
	if g.CheckForCheckColour(fig.RetColour()){
		output = false
	}
	if deleted{
		num = len(g.Figures)
		g.Figures = append(g.Figures, temporaryDeletedFig)
	}
	fig.SetCoords(prevCoords)
	return output, num
}

//Check if we can perform step without collisions
func (g * GameSession)CanStepWithNoCollisions(destination geometry.Point, fig StepMaker)bool{
	if fig.CheckStepAvailable(destination){
		if collided := g.FindAllCollisions(destination,fig);len(collided) == 0 {
			return true
		}
	}
	return false
}

///Check if we can attack without collisions
func (g * GameSession)CanAttack(destination geometry.Point, fig StepMaker)(int,bool){
	collisionFigs := g.FindAllCollisions(destination,fig)
	var canAttack bool = false
	collision := -1
	if len(collisionFigs) == 1{
		//we can attack only enemies
		canAttack = g.Figures[collisionFigs[0]].isEnemy(fig)
		//we can't attack figures if collision is before destination
		canAttack = canAttack && g.Figures[collisionFigs[0]].RetCoords().Equal(destination)
		//we can't attack king
		canAttack = canAttack && !g.isKing(g.Figures[collisionFigs[0]])
		collision = collisionFigs[0]
	}
	return collision,canAttack
}

func (g * GameSession)isKing(fig StepMaker)bool{
	for _,element := range g.kings{
		if element == fig{
			return true
		}
	}
	return false
}

func (g * GameSession)StepVirtually(destination geometry.Point, fig StepMaker)(prevCoord geometry.Point){
	prevCoord  = fig.RetCoords()
	fig.SetCoords(destination)
	return
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

//Simple check if kings is under attack or not
func (g GameSession)CheckForCheck()int{
	for i,king := range g.kings{
		if g.CheckKing(king){
			return i
		}
	}
	return 0
}

func (g GameSession)CheckForCheckColour(c Colour)bool{
	return g.CheckKing(g.kings[c])
}

func (g GameSession)CheckKing(king * King)bool{
	for _,fig := range g.Figures{
		if _,check :=g.CanAttack(king.Point,fig);check{
			return true
		}
	}
	return false
}

func (g * GameSession)Act(clicked geometry.Point)GameSessionJsonRepr{
	var repr = GameSessionJsonRepr{}
	if g.chosenFigure != nil{
		if can, collision := g.CanAct(clicked, g.chosenFigure);can{
			if collision >= 0{
				g.Figures = append(g.Figures[:collision], g.Figures[collision+1:]...)
			}
			g.chosenFigure.Step(clicked)
			g.PlayingNow = (g.PlayingNow + 1) & 1
		}
		g.chosenFigure = nil
	}else if fig,ok := g.At(clicked);ok{
		g.chosenFigure = fig
		repr.GameOver = 0
		tmpProbSteps := fig.ListStepsAvailable()
		tmpProbSteps = append(tmpProbSteps,fig.AttacksAvailable()...)
		repr.ProbSteps = make([]geometry.Point,0,32)
		for _,element := range tmpProbSteps{
			if ok, _ :=g.CanAct(element,fig);ok{
				repr.ProbSteps = append(repr.ProbSteps, element)
			}
		}
	}
	repr.Figs = g.ToJsonRepr()
	log.Println(repr)
	return repr
}