package atd

import (
    "time"
    "log"
    "encoding/json"
)

type TeamsService struct {}

type teamsIndexReq struct {
    Message string `json:"message"`
}

type teamsIndexRes struct {
    Response string `json:"response"`
}

func (service TeamsService) Index(in []byte) []byte {
    var req teamsIndexReq
    if err := json.Unmarshal(in, &req); err != nil {
        log.Print("TeamsService.Index(): received invalid request.")
        return nil
    }
    response := "[TeamsService] " + req.Message + "[" + time.Now().String() + "]"
    println("TeamsService received:", req.Message)
    println("TeamsService responds:", response)
    res := teamsIndexRes{response}
    if out, err := json.Marshal(res); err == nil {
        return out
    }
    log.Print("TeamsService.Index(): couldn't marshal response.")
    return nil
}
