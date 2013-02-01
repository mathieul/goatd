package web

import (
    "os"
    "net/http"
    "github.com/gorilla/rpc"
    "github.com/gorilla/rpc/json"
)

func ServeApplication() {
    server := rpc.NewServer()
    server.RegisterCodec(json.NewCodec(), "application/json")
    server.RegisterService(new(TestService), "Test")
    http.Handle("/rpc", server)

    dir, _ := os.Getwd()
    http.Handle("/", http.FileServer(http.Dir(dir + "/public/")))
    http.ListenAndServe(":8080", nil)
}