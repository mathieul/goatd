package main

import (
    "net/http"
    "github.com/gorilla/rpc"
    "github.com/gorilla/rpc/json"
)

func main() {
    s := rpc.NewServer()
    s.RegisterCodec(json.NewCodec(), "application/json")
    s.RegisterService(new(HelloService), "")
    println("before handle...")
    http.Handle("/rpc", s)
    http.ListenAndServe(":8080", nil)
}

type HelloArgs struct {
    Who string
}

type HelloReply struct {
    Message string
}

type HelloService struct {}

// curl -v http://localhost:8080/rpc -d '{"method":"HelloService.Say", "params": [{"Who": "blah"}], "id": 99}' -H "Content-Type: application/json"
// {"result":{"Message":"Hello, blah!"},"error":null,"id":99}
func (h *HelloService) Say(r *http.Request, args *HelloArgs, reply *HelloReply) error {
    reply.Message = "Hello, " + args.Who + "!"
    return nil
}
