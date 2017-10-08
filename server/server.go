package server

import(
	"github.com/gorilla/websocket"
	"chessServer/model"
	"encoding/json"
	"chessServer/utility/geometry"
	"runtime"
	"log"
)


type GameServer struct{
	Rooms []*GameRoom
	GameBalancer
}

func NewServer()*GameServer{
	return &GameServer{Rooms:make([]*GameRoom,0,100),
		GameBalancer:GameBalancer{Incoming:make(chan * websocket.Conn, 1000),
			out1:make(chan * websocket.Conn, 500), out2:make(chan * websocket.Conn, 500),counter:0}}
}

func (g * GameServer)SchedGames(){
	go g.GameBalancer.splitToPairs()
	for{
		if runtime.NumGoroutine() < 1000{
			log.Println("scheduling available")
			first,second := <-g.out1, <- g.out2
			log.Println("received both messages")
			if first != nil && second != nil{
				room := New(first,second, len(g.Rooms),g)
				g.Rooms = append(g.Rooms, room)
				go room.run()
			}
		}
	}
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

func New(first, second * websocket.Conn,index int, server * GameServer) * GameRoom {
	g :=  &GameRoom{}
	g.session = model.New()
	g.Conns = make(map[*websocket.Conn]int)
	g.Conns[first] = 1
	g.Conns[second] = 0
	g.msg = make(chan message)
	g.leave = make(chan *websocket.Conn)
	g.index = index
	g.server = server
	g.SendState(g.session.InitialToJsonRepr())
	return g
}

func(g * GameRoom) listen(conn * websocket.Conn){
	for{
		_, data, err := conn.ReadMessage()
		if err != nil{
			log.Println("Error, server.go line 69")
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
			log.Println("error server.go line:88")
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
					}else{
						log.Println("Error, server.go line:109")
					}
				}
			case left := <- g.leave:	//	if player leaves
				g.SendState(model.GameSessionJsonRepr{GameOver: 1 + (^g.Conns[left])})	//over player wins
				for key := range g.Conns{	//	close all connections
					log.Println("close connection")
					key.Close()
				}
				g.server.Rooms = append(g.server.Rooms[:g.index], g.server.Rooms[g.index+1:]...)	//	remove session and game room
		}
	}
}



type GameBalancer struct{
	Incoming chan *websocket.Conn
	out1 chan * websocket.Conn
	out2 chan * websocket.Conn
	counter uint64
}

func (g * GameBalancer)splitToPairs(){
	for{
		select {
		case conn := <- g.Incoming:
			if g.counter % 2 == 0{
				g.out1 <- conn
				log.Println("received first conn")
			}else{
				g.out2 <- conn
				log.Println("received second conn")
			}

			if g.counter + 1 < g.counter{
				log.Println("overflow")
				g.counter = 0
			}else{
				g.counter += 1
			}
		}
	}
}




