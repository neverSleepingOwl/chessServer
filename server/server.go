package server

import(
	"fmt"
	"net"
	"os"
)
//TODO add DDOS protection

const(
	maxRequests int=1000000
	CONN_HOST = "localhost"
	CONN_PORT = "3333"
	CONN_TYPE = "tcp"
)

type PostCard struct{	//	PostCard struct is analogy of real postcard,
	message string		// to deal with multiple receivers we need address of receiver(connection)
	address net.Conn	//	and message to send
}
type TcpServer struct{//TODO add slice of connections
	Balancer GameBalancer
	Input chan PostCard	//	chanel for incoming data
	Output chan PostCard	//	chanel for sent data
}

func ConstructTcpServer()TcpServer{
	return TcpServer{ConstructGameBalancer(),make(chan PostCard,maxRequests),make(chan PostCard, maxRequests)}
}

func (t * TcpServer)waitForConnection(){
	l,err:=net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
	if err != nil{
		fmt.Println("Error while listening to socket", err.Error())
	}
	defer l.Close()
	for{
		connection,err:=l.Accept()
		if err != nil{
			fmt.Print("Error while accepting connection", err.Error())
		}else{
			go func{
				buffer:=make([]byte,1024)
				_,err := connection.Read(buffer)

				if err != nil{
					fmt.Println("Error while receiving data:",err.Error())
				}else{
					t.Input<-PostCard{string(buffer),connection}
				}

			}()
		}

	}
}

func (t * TcpServer)handleRequests{
	go func{
		for{
			received:=<-t.Input
			if received.message == "wannaplay" || received.message == "wannaplay\n"{
				t.Balancer.BalanceRequestsEnqueue(received.address)
			}
		}
	}()
}


type GameBalancer struct{	//	class which builds pairs of players from incoming new game requests
	putToFirst bool		//	put incoming client in first or second queue? this condition need to divide
	//	incoming queue into two queues with close amount of clients
	FirstQueue chan net.Conn	//	queues themselves are implemented as buffered channels
	SecondQueue chan net.Conn
	SessionCounter uint64	//	session counter to create sessions with unique keys in DataBase
}

func ConstructGameBalancer()GameBalancer{	//	constructor
	return GameBalancer{true,make(chan net.Conn, maxRequests/2), make(chan net.Conn, maxRequests/2),0}
}

func (g * GameBalancer) BalanceRequestsEnqueue(client net.Conn){	//	split incoming requests in two queues
	if g.putToFirst{
		g.FirstQueue<-client
	}else{
		g.SecondQueue<-client
	}
	g.putToFirst=!g.putToFirst
}

func (g * GameBalancer)BalanceRequestsDeque(){
	first:=<-g.FirstQueue
	second:=<-g.SecondQueue
	//TODO add game creation
	//TODO add interaction with DB
}







