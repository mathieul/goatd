package main

import (
    "goatd/app/tcp"
    "goatd/app/atd"
)

func main() {
    go atd.GetInstance().Run()
    tcp.ServeRequests(4242)
}