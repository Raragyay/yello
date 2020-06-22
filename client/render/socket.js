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

function load_level(levelData) {
    level = []
    pellets = []
    let split_level_data = levelData.split('-')
    let cols = parseInt(split_level_data[0])
    let rows = parseInt(split_level_data[1])
    let idx = 2
    for (let i = 0; i < rows; i++) {
        let row_data = []
        for (let j = 0; j < cols; j++) {
            row_data.push(split_level_data[idx])
            if (split_level_data[idx] === '002') {
                pellets.push(new Pellet(j, i))
            } else if (split_level_data[idx] === '003') {
                pellets.push(new BigPellet(j, i))
            } else if (split_level_data[idx] === '010') {
                blinky.setPosition(j, i)
            } else if (split_level_data[idx] === '011') {
                pinky.setPosition(j, i)
            } else if (split_level_data[idx] === '012') {
                inky.setPosition(j, i)
            } else if (split_level_data[idx] === '013') {
                clyde.setPosition(j, i)
            } else if (split_level_data[idx] === '004') {
                pacman.setPosition(j, i)
            }
            idx++
        }
        level.push(row_data)
    }
    isWall = level.map(x => x.map(tile => tile === '100'))
    calc_block_size()
}

function updatePositionOfEntity(data) {
    let split_data = data.split(' ')[2].split('-');
    let x = parseInt(split_data[1])
    let y = parseInt(split_data[2])
    if (split_data[0] === 'P1') {
        //VECTOR TOSTRING SYNTAX WILL HAVE TO BE GIVEN. CURRENTLY ASSUME IT IS GIVEN AS {X}-{Y}
        pacman.setPosition(x, y)
    } else if (split_data[0] === 'P2') {
        blinky.setPosition(x, y)
    } else if (split_data[0] === 'P3') {
        inky.setPosition(x, y)
    } else if (split_data[0] === 'P4') {
        pinky.setPosition(x, y)
    } else if (split_data[0] === 'P5') {
        clyde.setPosition(x, y)
    }
}

function nullifyPellet(data) {
    let split_data = data.split(' ')[2].split('-');
    let x = parseInt(split_data[0])
    let y = parseInt(split_data[1])
    let pellet = pellets.find(element => element.x === x && element.y === y)
    pellet.draw = () => {
    }
}

isQueueing = false

socket.onmessage = (msg) => {
    let data = msg.data;
    if (data.startsWith("PONG QUEUE")) {
        document.getElementById("queue").innerText = data.split(" ")[2]
        isQueueing = true
    } else if (data.startsWith("PONG GAME-INIT")) {
        document.getElementById("mainui-play").style.display = 'none'
        gameActive = true
        names = data.split(' ')[2].split('-')
        print(names)
        for (var i = 1; i < names.length; i++) {
            print("player" + toString(i) + "name")
            document.getElementById("player" + i.toString() + "name").innerText = names[i]
        }
    } else if (data.startsWith("PONG SET-LEVEL")) {
        load_level(data.split(' ')[2])
    } else if (data.startsWith("PONG GAME-INVALID")) {
        console.log("SERVER: INVALID GAME")
    } else if (data.startsWith("PONG GAME-ENTITY-POS")) { //PONG GAME-ENTITY-POS p1-3-5
        updatePositionOfEntity(data);
    } else if (data.startsWith("PONG GAME-PELLET-HOM")) {

        console.log("Server: " + msg.data);
        nullifyPellet(data);
    } else if (data.startsWith("PONG GAME-SCARED")) {

        console.log("Server: " + msg.data);
        //TODO change to reflect actual variable name
        isScared = data.split(' ')[2] === '1'
    }
}

document.getElementById("play").addEventListener("click", () => {
    if (isQueueing) {
        return
    }
    let name = document.getElementById("nick").value;
    sendSocketMessage("PONG " + name)
    sendSocketMessage("PONG QUEUE")
})
