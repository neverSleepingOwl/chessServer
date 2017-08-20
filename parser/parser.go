package parser

import(
	"strings"
	"regexp"
	"chessServer/model"
	"strconv"
)


//Parse string, stored in database into fields of struct GameSession

func split(s string)(tokens []string, ok bool){
	tokens = strings.Split(s, ";")
	reg, err:=regexp.Compile(`\w+:[a-z]+|[1-9]?\d*|\{[0-7]\,[0-7]\,[0-1]\}`)
	if err != nil{
		return []string{}, false
	}
	for _,element := range tokens{
		if !reg.MatchString(element){
			return []string{}, false
		}
	}
	return tokens, true
}


//Parse incoming JSON-style format strings, stored in Database
//format example: AuthToken:1234;PlayingNow:1;Pawn:{1,2,1};etc...
//records are divided by semicolons, record example: AuthToken:12312312
// return all data, required for session constructor
func GenerateGameSessionData(s string)(figures []model.StepMaker, authToken string, playingNow model.Colour, stepDone bool, ok bool){
	var(
		tokenMatch = regexp.MustCompile(`^AuthToken:[1-9]\d+|\d$`)	//	regexp to detect auth authToken
		playingNowMatch = regexp.MustCompile(`^PlayingNow:[0-1]$`)	//	reqexp to detect integer representation of colour
		figureMatch = regexp.MustCompile(`^(Pawn|Rook|Bishop|Knight|King|Queen):\{[0-7]\,[0-7]\,[0-1]\}$`)	//	regexp to detect figure data
		stepDoneMatch = regexp.MustCompile(`^StepDone:(true|false)$`)	//	regexp to detect bool value of stepDone
		figureDelimeter = regexp.MustCompile(`[:\{},]`)	//	set of delimeters to split figure data
	)

	if tokens,ok:=split(s);!ok{	//	split
		return []model.StepMaker{}, "", 0, false, false	//	if split fails return ok = false
	}else{
		figures = make([]model.StepMaker,0,32)	//	allocate memory

		for _,element := range tokens{	//	parse each record
			switch {
			case tokenMatch.MatchString(element):
				authToken = strings.Split(element,":")[1]

			case playingNowMatch.MatchString(element):
				tmp,_ := strconv.Atoi(strings.Split(element,":")[1])
				playingNow=model.Colour(tmp)

			case stepDoneMatch.MatchString(element):
				if tmp:=strings.Split(element,":");tmp[0] == "true"{
					stepDone = true
				}else{
					stepDone = false
				}

			case figureMatch.MatchString(element):
				tmp := figureDelimeter.Split(element, -1)

				x,_ := strconv.Atoi(tmp[1])
				y,_ := strconv.Atoi(tmp[2])

				col,_ := strconv.Atoi(tmp[3])
				colour:=model.Colour(col)

				switch tmp[0]{
				case "Pawn":
					colour:=model.Colour(col)
					figures = append(figures, model.ConstructPawn(x,y,colour))

				case "King":
					colour:=model.Colour(col)
					figures = append(figures, model.ConstructKing(x,y,colour))

				case "Queen":
					colour:=model.Colour(col)
					figures = append(figures, model.ConstructQueen(x,y,colour))

				case "Rook":
					figures = append(figures, model.ConstructRook(x,y,colour))

				case "Knight":
					figures = append(figures, model.ConstructKnight(x,y,colour))

				case "Bishop":
					figures = append(figures, model.ConstructBishop(x,y,colour))

				default:
					return []model.StepMaker{}, "", 0, false, false	//	if found unexpected figure record return ok = false
				}
			default:
				return []model.StepMaker{}, "", 0, false, false	//	if found unexpected record
			}
		}
		ok = true	//	return parsed data and ok code = true
		return
	}
}