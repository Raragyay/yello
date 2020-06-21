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
    let data = msg.data;
    if (data.startsWith("PONG QUEUE")) {
        document.getElementById("queue").innerText = data.split(" ")[2]
    } else if (data.startsWith("PONG GAME-INIT")) {
        //TODO take names of other players
        document.getElementById("mainui-play").style.display = 'none'
        gameActive = true
    } else if (data.startsWith("PONG GAME-UPDATE-POS")) {
        let split_data = data.split(' ');
        if (split_data[2] === 'p1') {
            //VECTOR TOSTRING SYNTAX WILL HAVE TO BE GIVEN. CURRENTLY ASSUME IT IS GIVEN AS {X}-{Y}
            let pvector = split_data[3].split('-');
            player1.xblock = parseInt(pvector[0])
            player1.yblock = parseInt(pvector[1])
        }
    }
    console.log("Server: " + msg.data);
}

document.getElementById("play").addEventListener("click", () => {
    let name = document.getElementById("nick").value;
    sendSocketMessage("PONG " + name)
    sendSocketMessage("PONG QUEUE")
})
