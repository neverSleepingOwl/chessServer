package logger

import "log"

const level int = 6 // max level of logging
//0 - for basic logging (main)
//1 - for server main parts
//2 and 3 for server classes
//4 for logging game session
//5 for logging figures
//6 for logging geometry
var levels  = map[int]string{
							0:"FROM MAIN",
							1:"FROM SERVER ",
							2:"FROM GAME BALANCER",
							3:"FROM GAME ROOM",
							4:"FROM GAME SESSION",
							5:"FROM FIGURES",
							6:"FROM GEOMETRY",
							}
func WriteLog(lvl int,v ...interface{}){
		if lvl <= level{
			log.Println("----------",levels[lvl],"----------")
			log.Println(v)
			log.Println("----------",levels[lvl],"----------")
		}
}

