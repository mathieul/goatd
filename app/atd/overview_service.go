package atd

import (
    "time"
)

type OverviewService struct {}

func (service OverviewService) Index(request []byte) []byte {
    message := string(request)
    response := "[OverviewService] " + message + "[" + time.Now().String() + "]"
    println("OverviewService received:", message)
    println("OverviewService responds:", response)
    return []byte(response)
}
