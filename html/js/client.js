var canvas = document.querySelector("canvas");
var context = canvas.getContext("2d");

var figImages = {'king':'', 'queen':'','rook':'','knight':'','bishop':'','pawn':''};
var  images = {'bg':'','vic1':'','vic2':'','draw':''};
var probStep = document.createElement("img");
probStep.src = 'images/prob.png';

var preview_pic = document.createElement("img");
preview_pic.src = 'images/game.jpg';

var gameData = {
    player:1
};

for (var key in images){    //  cache images
    images[key] = document.createElement("img");
    images[key].src = 'images/' + key + '.png'
}

for (var figKey in figImages){
    figImages[figKey] = [document.createElement("img"),document.createElement("img")];
    for (var i = 0; i < figImages[figKey].length; i++){
        figImages[figKey][i].src = 'images/'+ figKey + i + '.png';
    }
}

context.clearRect(0, 0, canvas.width, canvas.height);
context.drawImage(preview_pic,0,0,canvas.width, canvas.height);

var socket = new WebSocket("ws://chessserver.herokuapp.com:8080/ws");    //  init websocket connection

socket.onopen = function() {
    console.log("Connection established");
};

socket.onclose = function(event) {  //  set callbacks
    if (event.wasClean) {
        console.log('Connection closed clear');
    } else {
        console.log('Connection refused by server');
    }
    console.log('Close code: ' + event.code + ' reason : ' + event.reason);
};

socket.onmessage = function(event) {
    console.log("receivedData " + event.data);
    gameData = JSON.parse(event.data);
    repaint(gameData)
};

socket.onerror = function(error) {
    console.log("Error:  " + error.message);
};
var coords = {
    x:0,
    y:0
};

canvas.addEventListener('mouseup', function (e) {
    if (gameData.player === 1){
        coords.x = Math.floor((e.pageX - e.target.offsetLeft)/50);
        coords.y = Math.floor((e.pageY - e.target.offsetTop)/50);
    }else{
        coords.x = 7 - Math.floor((e.pageX - e.target.offsetLeft)/50);
        coords.y = 7 - Math.floor((e.pageY - e.target.offsetTop)/50);
    }
    console.log(coords);
    socket.send(JSON.stringify(coords))
});
/*
type GameSessionJsonRepr struct{
	Figs []FigJsonRepr `json:"figs"`
	GameOver int 		`json:"game_over"` // 0 for normal game, 1 for black, 2 for white, 3 for draw
	ProbSteps []geometry.Point `json:"list_steps"`
	Player int `json:"player"`	//	player, for clients to know what colour does he has << is that a fucking psycho pass reference??
}

type FigJsonRepr struct{
	Name string `json:"name"`
	X   int `json:"x"`
	Y   int `json:"y"`
	Col int `json:"colour"`
}
 */
function repaint(gameData){
    context.clearRect(0, 0, canvas.width, canvas.height);
    context.drawImage(images['bg'],0,0,canvas.width, canvas.height);
    if (gameData.game_over === 0){
        if (gameData.list_steps !== null){
            for (var j = 0; j < gameData.list_steps.length;j++){
                if(gameData.player === 1){
                    context.drawImage(probStep,gameData.list_steps[j].x * 50, gameData.list_steps[j].y * 50, 50, 50);
                }else{
                    context.drawImage(probStep,(7 - gameData.list_steps[j].x) * 50, (7-gameData.list_steps[j].y) * 50, 50, 50);
                }
            }
        }
        for (var i = 0; i < gameData.figs.length;i++){
            if(gameData.player === 1){
                context.drawImage(figImages[gameData.figs[i].name][gameData.figs[i].colour],
                    gameData.figs[i].x * 50, gameData.figs[i].y * 50, 50, 50);
            }else{
                context.drawImage(figImages[gameData.figs[i].name][gameData.figs[i].colour],
                    (7 - gameData.figs[i].x) * 50, (7 - gameData.figs[i].y)*50, 50, 50);
            }
        }
        return
    }else if(gameData.game_over > 0 && gameData.game_over < 3){
        context.drawImage(images['vic'+gameData.game_over],0,0,canvas.width, canvas.height);
    }else{
        context.drawImage(images['vic'+gameData.game_over],0,0,canvas.width, canvas.height);
    }
    socket.close()
}
