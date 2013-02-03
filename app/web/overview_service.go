package web

import (
    "net/http"
)

type OverviewArgs struct {
}

type row struct {
    Name string
    Count int
}

type OverviewReply struct {
    Rows []row
}

type OverviewService struct {}

// curl -v http://localhost:8080/rpc -d '{"method":"Overview.List", "params": [], "id": 42}' -H "Content-Type: application/json"
// {"result":{"Rows": [{"Name": "Team", "Count": 12}]},"error":null,"id":42}
func (service *OverviewService) List(r *http.Request, args *OverviewArgs, reply *OverviewReply) error {
    reply.Rows = []row{{"Team", 2}, {"Teammate", 5}, {"Queue", 3}, {"Skill", 10}, {"Task", 42}}
    return nil
}
