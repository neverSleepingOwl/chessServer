package main

import (
	"net/http"
	"github.com/gorilla/websocket"
	"chessServer/server"
	"log"
	"runtime/debug"
)

const ADDR = ":8080"

func main(){
	debug.SetGCPercent(-1)
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
		log.Println("run callback")
		ws, err := upgrader.Upgrade(w,r,nil)
		if _, ok := err.(websocket.HandshakeError); ok {
			http.Error(w, "Not a websocket handshake", 400)
			log.Println("Error while receiving connection", err.Error())
			return
		} else if err != nil {
			log.Println("Error while receiving connection", err.Error())
			return
		}
		log.Println("WS OK")
		serv.Incoming <- ws
	})
	if err := http.ListenAndServe(ADDR, nil); err != nil {
		log.Println("ListenAndServe:", err)
	}
}


