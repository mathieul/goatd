package web

import (
    "os"
    "fmt"
    "net/http"
    "github.com/gorilla/rpc"
    "github.com/gorilla/rpc/json"
)

func ServeApplication(port int) {
    server := rpc.NewServer()
    server.RegisterCodec(json.NewCodec(), "application/json")
    server.RegisterService(new(OverviewService), "Overview")
    http.Handle("/rpc", server)

    dir, _ := os.Getwd()
    http.Handle("/", http.FileServer(http.Dir(dir + "/webapp/public/")))
    address := fmt.Sprintf(":%d", port)
    http.ListenAndServe(address, nil)
}