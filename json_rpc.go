package main

import (
    "os"
    "net/http"
    "github.com/gorilla/rpc"
    "github.com/gorilla/rpc/json"
)

func init() {
    s := rpc.NewServer()
    s.RegisterCodec(json.NewCodec(), "application/json")
    s.RegisterService(new(HelloService), "")
    http.Handle("/rpc", s)
}

func main() {
    dir, _ := os.Getwd()
    dir += "/public/"
    http.Handle("/public/", http.FileServer(http.Dir(dir)))
    println("Serving ", dir)
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
