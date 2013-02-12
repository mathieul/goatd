package atd

import (
    "time"
)

type OverviewService struct {}

type OverviewIndexReq struct {
    Message string `json:"message"`
}

type OverviewIndexRes struct {
    Response string `json:"response"`
}

func (service OverviewService) Index(req OverviewIndexReq) OverviewIndexRes {
    response := "[OverviewService] " + req.Message + "[" + time.Now().String() + "]"
    println("OverviewService received:", req.Message)
    println("OverviewService responds:", response)
    return OverviewIndexRes{response}
}
