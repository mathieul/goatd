package main

import (
    zmq "github.com/alecthomas/gozmq"
    // "fmt"
    // "time"
    // "goatd/app/event"
    // "goatd/app/model"
    // "goatd/app/dispatch"
)

func main() {
    context, _ := zmq.NewContext()
    socket, _ := context.NewSocket(zmq.REP)
    socket.Bind("tcp://127.0.0.1:5000")

    for {
        msg, _ := socket.Recv(0)
        println("Got", string(msg))
        socket.Send(msg, 0)
    }
}