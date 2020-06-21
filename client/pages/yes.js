let socket = new WebSocket("ws://185.163.47.170:5000/ws");
console.log("Attempting Connection...");

var socketOpen = false;

socket.onopen = () => {
    console.log("Successfully Connected");
    socketOpen = true;
};
        
socket.onclose = event => {
    console.log("Socket Closed Connection: ", event);
 	socket.send("PONG CLOSE")
 	socketOpen = false;
};

socket.onerror = error => {
    console.error("Socket Error: ", error);
};

socket.onmessage = (msg) =>{
    console.log("Server: " + msg);
}


function sendSocketMessage(msg){
	if(!socketOpen){
		console.error("tried to write to socket despite socket being closed");
		return;
	}
	socket.send(msg);
	console.log("SENT: " + msg);
}

function onSubmitName(){
	var name = document.getElementById("onQueueButton").textContent();
	sendSocketMessage("PONG " + name);
}