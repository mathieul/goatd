package tcp

import (
    "fmt"
    "time"
    zmq "github.com/alecthomas/gozmq"
)

const (
    overviewAddress = "ipc://overview.ipc"
    teamsAddress = "ipc://teams.ipc"
    addressPrefix  = "tcp://127.0.0.1:%d"
    noTimeOut = -1
)

func ServeRequests(port int) {
    launchServices()

    clientAddress := fmt.Sprintf(addressPrefix, port)
    context, _ := zmq.NewContext()
    defer context.Close()

    frontend, _ := context.NewSocket(zmq.ROUTER)
    defer frontend.Close()
    frontend.Bind(clientAddress)

    overviewSocket, _ := context.NewSocket(zmq.DEALER)
    defer overviewSocket.Close()
    overviewSocket.Bind(overviewAddress)

    teamsSocket, _ := context.NewSocket(zmq.DEALER)
    defer teamsSocket.Close()
    teamsSocket.Bind(teamsAddress)

    toPoll := zmq.PollItems{
        zmq.PollItem{zmq.Socket: frontend,       zmq.Events: zmq.POLLIN},
        zmq.PollItem{zmq.Socket: overviewSocket, zmq.Events: zmq.POLLIN},
        zmq.PollItem{zmq.Socket: teamsSocket,    zmq.Events: zmq.POLLIN},
    }
    
    for {
        _, _ = zmq.Poll(toPoll, noTimeOut)

        switch {
            case toPoll[0].REvents & zmq.POLLIN != 0:
                messages, _ := toPoll[0].Socket.RecvMultipart(0)
                serviceName := string(messages[len(messages) - 1])
                println("Request for service:", serviceName)
                messages = messages[:len(messages) - 1]
                switch serviceName {
                case "overview":
                    overviewSocket.SendMultipart(messages, 0)
                case "teams":
                    teamsSocket.SendMultipart(messages, 0)
                }
                
            case toPoll[1].REvents & zmq.POLLIN != 0:
                messages, _ := toPoll[1].Socket.RecvMultipart(0)
                frontend.SendMultipart(messages, 0)

            case toPoll[2].REvents & zmq.POLLIN != 0:
                messages, _ := toPoll[2].Socket.RecvMultipart(0)
                frontend.SendMultipart(messages, 0)
        }
    }
}

func launchServices() {
    go overviewService()
    go teamsService()
}

func timestampedMessage(message string) string {
    return message + "[" + time.Now().String() + "]"
} 

func overviewService() {
    context, _ := zmq.NewContext()
    defer context.Close()

    receiver, _ := context.NewSocket(zmq.REP)
    defer receiver.Close()
    receiver.Connect(overviewAddress)

    for {
        received, _ := receiver.Recv(0)
        println("           Overview:", string(received))
        receiver.Send([]byte(timestampedMessage("Overview Service")), 0)
    }
}

func teamsService() {
    context, _ := zmq.NewContext()
    defer context.Close()

    receiver, _ := context.NewSocket(zmq.REP)
    defer receiver.Close()
    receiver.Connect(teamsAddress)

    for {
        received, _ := receiver.Recv(0)
        println("              Teams:", string(received))
        receiver.Send([]byte(timestampedMessage("Teams Service")), 0)
    }
}

