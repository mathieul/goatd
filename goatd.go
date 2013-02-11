package main

import (
    "goatd/app/tcp"
    "goatd/app/atd"
)

func main() {
    go atd.GetInstance().Run()
    server := tcp.NewServer(4242)
    server.registerService(new(atd.OverviewService), "overview")
    server.registerService(new(atd.TeamsService), "teams")
    server.ListenAndReply()
}