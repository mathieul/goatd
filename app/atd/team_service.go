package atd

import (
    "goatd/app/model"
)

type TeamService struct {}

/*
 * Team.List
 */

type TeamIndexReq struct {}

type TeamAttributes struct {
    Uid string `json:"uid"`
    Name string `json:"name"`
}

type TeamIndexRes struct {
    Teams []TeamAttributes `json:"teams"`
}

func (service TeamService) List(req TeamIndexReq) TeamIndexRes {
    res := new(TeamIndexRes)

    teams := make([]TeamAttributes, 0)
    store().Teams.Select(func(item interface{}) bool {
        team := item.(*model.Team)
        teams = append(teams, TeamAttributes{team.Uid(), team.Name()})
        return false
    })
    res.Teams = teams

    return *res
}

/*
 * Team.Create
 */
type TeamCreateReq struct {
    Name string
}

func (service TeamService) Create(req TeamCreateReq) TeamAttributes {
    res := new(TeamAttributes)

    team  := store().Teams.Create(model.A{"Name": req.Name})
    res.Uid = team.Uid()
    res.Name = team.Name()

    return *res
}

/*
 * Team.Update
 */

type TeamUpdateableAttributes struct {
    Name string `json:"name"`
}

type TeamUpdateReq struct {
    Uid string
    Team TeamUpdateableAttributes
}

type TeamUpdateRes struct {}

func (service TeamService) Update(req TeamUpdateReq) TeamUpdateRes {
    if team := store().Teams.Find(req.Uid); team != nil {
        if len(req.Team.Name) > 0 {
            team.Update("Name", req.Team.Name)
        }
    }

    return TeamUpdateRes{}
}

/*
 * Team.Destroy
 */

type TeamDestroyReq struct {
    Uid string
}

type TeamDestroyRes struct {}

func (service TeamService) Destroy(req TeamDestroyReq) TeamDestroyRes {
    if team := store().Teams.Find(req.Uid); team != nil {
        team.Destroy()
    }

    return TeamDestroyRes{}
}
