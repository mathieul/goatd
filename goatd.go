package main

import (
    "goatd/app/atd"
    "goatd/app/web"
    "goatd/app/tcp"
)

func main() {
    go atd.Run()
    go web.ServeApplication(8080)
    tcp.ServeRequests(5000)
}