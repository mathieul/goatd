package atd

import (
    "goatd/app/model"
)

type TeammateService struct {}

/*
 * Teammate.List
 */

type TeammateIndexReq struct {}

type TeammateAttributes struct {
    Uid string `json:"uid"`
    Name string `json:"name"`
    TeamUid string `json:"team_uid"`
}

type TeammateIndexRes struct {
    Teammates []TeammateAttributes `json:"teammates"`
}

func (service TeammateService) List(req TeammateIndexReq) TeammateIndexRes {
    res := new(TeammateIndexRes)

    teammates := make([]TeammateAttributes, 0)
    store().Teammates.Select(func(item interface{}) bool {
        teammate := item.(*model.Teammate)
        teammates = append(teammates, TeammateAttributes{
            teammate.Uid(),
            teammate.Name(),
            teammate.TeamUid(),
        })
        return false
    })
    res.Teammates = teammates

    return *res
}

/*
 * Teammate.Create
 */
type TeammateCreateReq struct {
    TeamUid string `json:"team_uid"`
    Name string `json:"name"`
}

func (service TeammateService) Create(req TeammateCreateReq) TeammateAttributes {
    res := new(TeammateAttributes)

    teammate := store().Teammates.Create(model.A{"Name": req.Name, "TeamUid": req.TeamUid})
    res.Uid = teammate.Uid()
    res.Name = teammate.Name()
    res.TeamUid = teammate.TeamUid()

    return *res
}

/*
 * Teammate.Update
 */

type TeammateUpdateableAttributes struct {
    Name string `json:"name"`
}

type TeammateUpdateReq struct {
    Uid string `json:"uid"`
    Teammate TeammateUpdateableAttributes
}

type TeammateUpdateRes struct {}

func (service TeammateService) Update(req TeammateUpdateReq) TeammateUpdateRes {
    if teammate := store().Teammates.Find(req.Uid); teammate != nil {
        if len(req.Teammate.Name) > 0 {
            teammate.Update("Name", req.Teammate.Name)
        }
    }

    return TeammateUpdateRes{}
}

/*
 * Teammate.Destroy
 */

type TeammateDestroyReq struct {
    Uid string `json:"uid"`
}

type TeammateDestroyRes struct {}

func (service TeammateService) Destroy(req TeammateDestroyReq) TeammateDestroyRes {
    if teammate := store().Teammates.Find(req.Uid); teammate != nil {
        teammate.Destroy()
    }

    return TeammateDestroyRes{}
}
