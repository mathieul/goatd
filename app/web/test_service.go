package web

import (
    "fmt"
    "net/http"
)

type TestArgs struct {
    Name string
    Number int
}

type TestReply struct {
    Square int
    Message string
}

type TestService struct {}

// curl -v http://localhost:8080/rpc -d '{"method":"Test.Run", "params": [{"Name": "Mathieu", "Number": 3}], "id": 42}' -H "Content-Type: application/json"
// {"result":{"Square":9,"Message":"As requested: square(3) = 9"},"error":null,"id":42}
func (h *TestService) Run(r *http.Request, args *TestArgs, reply *TestReply) error {
    reply.Square = args.Number * args.Number
    reply.Message = fmt.Sprintf("As requested: square(%d) = %d", args.Number, reply.Square)
    return nil
}
