package main

import (
    "goatd/app/tcp"
    "goatd/app/atd"
)

func main() {
    go atd.GetInstance().Run()
    server := tcp.NewServer(4242)
    server.RegisterService(new(atd.OverviewService), "overview")
    server.RegisterService(new(atd.TeamService), "team")
    server.RegisterService(new(atd.TeamService), "teammate")
    server.ListenAndReply()
}