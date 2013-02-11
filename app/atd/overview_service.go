package atd

import (
    "time"
    "log"
    "encoding/json"
)

type OverviewService struct {}

type overviewIndexReq struct {
    Message string `json:"message"`
}

type overviewIndexRes struct {
    Response string `json:"response"`
}

func (service OverviewService) Index(in []byte) []byte {
    var req overviewIndexReq
    if err := json.Unmarshal(in, &req); err != nil {
        log.Print("OverviewService.Index(): received invalid request.")
        return nil
    }
    response := "[OverviewService] " + req.Message + "[" + time.Now().String() + "]"
    println("OverviewService received:", req.Message)
    println("OverviewService responds:", response)
    res := overviewIndexRes{response}
    if out, err := json.Marshal(res); err == nil {
        return out
    }
    log.Print("OverviewService.Index(): couldn't marshal response.")
    return nil
}
