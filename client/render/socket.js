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
    console.log('Client: CONNECTION ESTABLISHED')
    socketOpen = true
}

socket.onclose = event => {
    console.log("Socket Closed Connection: ", event);
    sendSocketMessage("PONG CLOSE")
    socketOpen = false;
}
socket.onerror = error => {
    console.error("Socket Error: ", error);
};

socket.onmessage = (msg) => {
    var data = msg.data
    if (data.startsWith("PONG QUEUE")) {
        document.getElementById("queue").innerText = data.split(" ")[2]
    }
    console.log("Server: " + msg.data);
}

document.getElementById("play").addEventListener("click", () => {
    var name = document.getElementById("nick").value;
    sendSocketMessage("PONG " + name)
    sendSocketMessage("PONG QUEUE")
})
