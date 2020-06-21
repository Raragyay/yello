var socket=io('http://localhost:5000')
socket.on('connect',()=>{
    console.log('CONNECTION ESTABLISHED')
})
