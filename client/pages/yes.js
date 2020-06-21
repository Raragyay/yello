let socket = new WebSocket("ws:http://185.163.47.170:5000/ws");
console.log("Attempting Connection...");

socket.onopen = () => {
        console.log("Successfully Connected");
        socket.send("Hi From the Client!")
    };
        
    socket.onclose = event => {
    console.log("Socket Closed Connection: ", event);
    socket.send("Client Closed!")
    };

    socket.onmessage = (msg) =>{
    	console.log(msg);
    }

    socket.onerror = error => {
        console.log("Socket Error: ", error);
    };