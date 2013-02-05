package web

import (
    "net/http"
    "goatd/app/model"
    "goatd/app/atd"
)

type TeamService struct {}

type TeamAttributes struct {
    Uid string
    Name string
}

/*
 * Team.List
 */
type TeamListReply struct {
    Rows []TeamAttributes
}

type TeamListArgs struct {}

// curl -v http://localhost:8080/rpc -d '{"method":"Team.List", "params": [], "id": 42}' -H "Content-Type: application/json"
func (service *TeamService) List(r *http.Request, args *EmptyStruct, reply *TeamListReply) error {
    store := atd.GetInstance().Store()

    rows := make([]TeamAttributes, 0)
    store.Teams.Select(func(item interface{}) bool {
        team := item.(*model.Team)
        rows = append(rows, TeamAttributes{team.Uid(), team.Name()})
        return false
    })
    reply.Rows = rows
    return nil
}


/*
 * Team.Create
 */
type TeamCreateArgs struct {
    Name string
}

// curl -v http://localhost:8080/rpc -d '{"method":"Team.Create", "params": [{"Name": "Toto"}], "id": 42}' -H "Content-Type: application/json"
func (service *TeamService) Create(r *http.Request, args *TeamCreateArgs, reply *TeamAttributes) error {
    store := atd.GetInstance().Store()
    team  := store.Teams.Create(model.A{"Name": args.Name})
    reply = new(TeamAttributes)
    reply.Uid = team.Uid()
    reply.Name = team.Name()
    // reply = &TeamAttributes{team.Uid(), team.Name()}
    return nil
}
