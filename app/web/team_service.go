package web

import (
    "net/http"
    "goatd/app/model"
    "goatd/app/atd"
)

type TeamArgs struct {
}

type row struct {
    Uid string
    Name string
}

type TeamReply struct {
    Rows []row
}

type TeamService struct {}

func (service *TeamService) List(r *http.Request, args *TeamArgs, reply *TeamReply) error {
    store := atd.GetInstance().Store()

    rows := make([]row)
    store.Teams.Select(func(item interface{}) bool {
        team := item.(*model.Team)
        rows = append(rows, row{team.Uid(), team.Name()})
        return false
    })
    reply.Rows = rows
    return nil
}
