package main

import (
    "goatd/app/web"
    "goatd/app/tcp"
    "goatd/app/atd"
)

func main() {
    go atd.GetInstance().Run()
    go web.ServeApplication(8080)
    tcp.ServeRequests(5000)
}