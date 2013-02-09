package tcp

import (
    "fmt"
    zmq "github.com/alecthomas/gozmq"
)

func ServeRequests(port int) {
    context, _ := zmq.NewContext()
    socket, _ := context.NewSocket(zmq.REP)
    address := fmt.Sprintf("tcp://127.0.0.1:%d", port)
    socket.Bind(address)
    fmt.Println("listening for zeromq REQ requests on", address)

    for {
        msg, _ := socket.Recv(0)
        fmt.Println("Received:", string(msg))
        socket.Send(msg, 0)
        fmt.Println("sent:", string(msg))
    }
}