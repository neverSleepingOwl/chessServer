### chessServer

## Breif Description

Websocket chess game server, written in golang. All logic is provided by server.

## Usage
Clone to your $GOPATH/src/github.com/ .
Install gorilla/websocket (go get)
run : 
```
go build ./main
```
Change line:
```
var socket = new WebSocket("ws://chessserver.herokuapp.com:8080/ws");
```

with 
```
var socket = new WebSocket("ws://your_ip:your_docker_port_binded_to_8080/ws");
```
Run docker container and enjoy your server.

## Requirements:
docker-cli, docker daemon running, go installed.
