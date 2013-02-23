package tcp

import (
    "fmt"
    "strings"
    "log"
    "reflect"
    "encoding/json"
    zmq "github.com/alecthomas/gozmq"
)

const (
    clientAddressTemplate  = "tcp://*:%d"
    serviceAddressTemplate = "ipc://%s.ipc"
    noTimeOut = -1
)

type Server struct {
    port int
    services map[string]interface{}
}

func NewServer(port int) (server *Server) {
    server = new(Server)
    server.port = port
    server.services = make(map[string]interface{})
    return server
}

func (server Server) ListenAndReply() {
    server.launchServices()
    server.runBroker()
}

func (server *Server) RegisterService(service interface{}, name string) {
    server.services[name] = service
}

func newBoundSocket(context zmq.Context, address string, kind zmq.SocketType) zmq.Socket {
    socket, _ := context.NewSocket(kind)
    socket.Bind(address)
    return socket
}

func (server *Server) runBroker() {
    context, _ := zmq.NewContext()
    defer context.Close()

    clientAddress := fmt.Sprintf(clientAddressTemplate, server.port)
    frontend := newBoundSocket(context, clientAddress, zmq.ROUTER)
    defer frontend.Close()

    toPoll := zmq.PollItems{
        zmq.PollItem{zmq.Socket: frontend, zmq.Events: zmq.POLLIN},
    }
    socketByName := make(map[string]zmq.Socket)

    for name, _ := range server.services {
        serviceAddress := fmt.Sprintf(serviceAddressTemplate, name)
        serviceSocket := newBoundSocket(context, serviceAddress, zmq.DEALER)
        defer serviceSocket.Close()
        socketByName[name] = serviceSocket
        toPoll = append(toPoll,
            zmq.PollItem{zmq.Socket: serviceSocket, zmq.Events: zmq.POLLIN},
        )
    }
    numSockets := len(toPoll)

    for {
        zmq.Poll(toPoll, noTimeOut)

        if toPoll[0].REvents & zmq.POLLIN != 0 {
            messages, _ := toPoll[0].Socket.RecvMultipart(0)
            serviceName := string(messages[len(messages) - 1])
            println("Request for service:", serviceName)
            if serviceSocket, found := socketByName[serviceName]; found {
                messages = messages[:len(messages) - 1]
                println("forwarding to service socket")
                serviceSocket.SendMultipart(messages, 0)
            }
        } else {
            for i := 1; i < numSockets; i++ {
                if toPoll[i].REvents & zmq.POLLIN != 0 {
                    messages, _ := toPoll[i].Socket.RecvMultipart(0)
                    frontend.SendMultipart(messages, 0)
                    break
                }
            }
        }
    }
}

func (server Server) launchServices() {
    for name, service := range server.services {
        go server.runService(service, name)
    }
}

func (server Server) runService(service interface{}, name string) {
    context, _ := zmq.NewContext()
    defer context.Close()

    receiver, _ := context.NewSocket(zmq.REP)
    defer receiver.Close()

    serviceAddress := fmt.Sprintf(serviceAddressTemplate, name)
    receiver.Connect(serviceAddress)
    println("Service", name, "is ready.")

    for {
        received, _ := receiver.RecvMultipart(0)
        action := strings.Title(string(received[1]))
        fmt.Println("runService[", name, "]: callServiceAction[", action, "] >>>", string(received[0]))
        response := callServiceAction(service, action, received[0])
        fmt.Println("runService[", name, "]: callServiceAction <<<", string(response))
        receiver.Send(response, 0)
    }
}

func callServiceAction(service interface{}, methodName string, param []byte) []byte {
    value := reflect.ValueOf(service)
    method := value.MethodByName(methodName)
    if !method.IsValid() {
        log.Print(fmt.Errorf("callServiceAction(): missing service method %q.", methodName))
        return nil
    }
    request := decodeRequestForMethod(method, param)
    response := method.Call([]reflect.Value{request})
    if len(response) != 1 {
        log.Fatal(fmt.Errorf("callServiceAction(): method %q should just return 1 value.", methodName))
    }
    return encodeResponseForMethod(response[0])
}

func decodeRequestForMethod(method reflect.Value, encoded []byte) reflect.Value {
    inType := method.Type().In(0)
    request := reflect.New(inType).Interface()
    if err := json.Unmarshal(encoded, &request); err != nil {
        log.Print(fmt.Errorf("decodeRequestForMethod(): received invalid request."))
        return reflect.Zero(inType)
    }
    pRequest := reflect.ValueOf(request)
    return reflect.Indirect(pRequest)
}

func encodeResponseForMethod(response reflect.Value) []byte {
    if encoded, err := json.Marshal(response.Interface()); err == nil {
        return encoded
    }
    return []byte{}
}
