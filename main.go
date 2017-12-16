package main

import (
	"net/http"
	"github.com/gorilla/websocket"
	"github.com/chessServer/server"
	"log"
	"os"
	"github.com/chessServer/utility/logger"
	"fmt"
)

const ADDR = ":8080"

func init(){
	f, err := os.OpenFile("testlogfile", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	if err != nil {
		println("Error, can't open file: ", err.Error())
		os.Exit(1)
	}
	//defer f.Close()

	log.SetOutput(f)
}

func main(){
	fmt.Print("DAEMON STARTED")
        fs := http.FileServer(http.Dir("html"))
        http.Handle("/", fs)
	serv := server.NewServer()
	go serv.SchedGames()
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request){
		var upgrader = websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 16384,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		}
		logger.WriteLog(0,"run callback")
		ws, err := upgrader.Upgrade(w,r,nil)
		if _, ok := err.(websocket.HandshakeError); ok {
			http.Error(w, "Not a websocket handshake", 400)
			logger.WriteLog(0,"Error while receiving connection", err.Error())
			return
		} else if err != nil {
			logger.WriteLog(0,"Error while receiving connection", err.Error())
			return
		}
		logger.WriteLog(0,"WS OK")
		serv.Incoming <- ws
	})
	if err := http.ListenAndServe(ADDR, nil); err != nil {
		logger.WriteLog(0,"ListenAndServe:", err)
	}
}


