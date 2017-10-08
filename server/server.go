package server

import(
	"fmt"
	"net"
	"github.com/gorilla/websocket"
	"chessServer/model"
	"encoding/json"
	"chessServer/utility/geometry"
)


type GameServer struct{
	Rooms []*GameRoom
}

type GameRoom struct{
	Conns map[ * websocket.Conn]int
	leave chan * websocket.Conn
	msg chan message
	session model.GameSession
	index int
	server * GameServer
}

type message struct{
	from * websocket.Conn
	msg []byte
}

func New(first, second * websocket.Conn,index int) * GameRoom {
	g :=  &GameRoom{}
	g.session = model.New()
	g.Conns = make(map[*websocket.Conn]int)
	g.Conns[first] = 1
	g.Conns[second] = 0
	g.msg = make(chan message)
	g.leave = make(chan *websocket.Conn)
	g.index = index
	go g.SendState(g.session.InitialToJsonRepr())
	return g
}

func(g * GameRoom) listen(conn * websocket.Conn){
	for{
		_, data, err := conn.ReadMessage()
		if err != nil{
			break
		}
		g.msg <- message{conn, data}
	}
	g.leave <- conn
	conn.Close()
}

func (g * GameRoom)SendState(msg model.GameSessionJsonRepr){
	for key,value := range g.Conns{
		tmp_msg := msg
		tmp_msg.Player = value
		if value != int(g.session.PlayingNow){
			tmp_msg.ProbSteps = []geometry.Point{}
		}
		msgToSend,_:=json.Marshal(&tmp_msg)
		err := key.WriteMessage(websocket.TextMessage,msgToSend)
		if err != nil{
			g.leave <- key
			key.Close()
		}
	}
}

func (g * GameRoom)run(){
	for key := range g.Conns{	//	start listen to websocket connections
		go g.listen(key)
	}
	for{
		select{
			case msg := <- g.msg:
				if g.Conns[msg.from] == int(g.session.PlayingNow){
					clicked := geometry.Point{}
					err := json.Unmarshal(msg.msg,&clicked)
					if err == nil {
						output := g.session.Act(clicked)
						g.SendState(output)
					}
				}
			case left := <- g.leave:	//	if player leaves
				g.SendState(model.GameSessionJsonRepr{GameOver: 1 + (^g.Conns[left])})	//over player wins
				for key := range g.Conns{	//	close all connections
					key.Close()
				}
				g.server.Rooms = append(g.server.Rooms[:g.index], g.server.Rooms[g.index+1:]...)	//	remove session and game room
		}
	}
}








