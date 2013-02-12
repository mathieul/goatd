package atd

import (
    "time"
)

type TeamsService struct {}

type TeamsIndexReq struct {
    Message string `json:"message"`
}

type TeamsIndexRes struct {
    Response string `json:"response"`
}

func (service TeamsService) Index(req TeamsIndexReq) TeamsIndexRes {
    response := "[TeamsService] " + req.Message + "[" + time.Now().String() + "]"
    println("TeamsService received:", req.Message)
    println("TeamsService responds:", response)
    return TeamsIndexRes{response}
}
