let socket = new WebSocket("ws://185.163.47.170:5000/ws");

function sendSocketMessage(msg) {
    if (!socketOpen) {
        console.error("tried to write to socket despite socket being closed");
        return;
    }
    socket.send(msg);
    console.log("SENT: " + msg);
}

let socketOpen = false
socket.onopen = () => {
    console.log('CONNECTION ESTABLISHED')
    socketOpen = true

    sendSocketMessage("PONG " + "TESTER");
    sendSocketMessage("PONG QUEUE");
}
socket.onmessage = (msg) => {
    console.log("Server: " + msg.data);
}
