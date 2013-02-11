package atd

import (
    "time"
)

type TeamsService struct {}

func (service TeamsService) Index(request []byte) []byte {
    message := string(request)
    response := "[TeamsService] " + message + "[" + time.Now().String() + "]"
    println("TeamsService received:", message)
    println("TeamsService responds:", response)
    return []byte(response)
}
