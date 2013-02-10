package tcp

import (
    "fmt"
    "time"
    zmq "github.com/alecthomas/gozmq"
)

const (
    serviceAddress = "ipc://services.ipc"
    addressPrefix  = "tcp://127.0.0.1:%d"
)

func ServeRequests(port int) {
    launchServices()

    clientAddress := fmt.Sprintf(addressPrefix, port)
    context, _ := zmq.NewContext()
    defer context.Close()

    frontend, _ := context.NewSocket(zmq.ROUTER)
    defer frontend.Close()
    frontend.Bind(clientAddress)

    backend, _ := context.NewSocket(zmq.DEALER)
    defer backend.Close()
    backend.Bind(serviceAddress)

    zmq.Device(zmq.QUEUE, frontend, backend)
}

func launchServices() {
    go overviewService()
}

func overviewService() {
    context, _ := zmq.NewContext()
    defer context.Close()

    receiver, _ := context.NewSocket(zmq.REP)
    defer receiver.Close()
    receiver.Connect(serviceAddress)

    for {
        received, _ := receiver.Recv(0)
        fmt.Println("Received:", string(received))
        receiver.Send([]byte(time.Now().String()), 0)
    }
}
