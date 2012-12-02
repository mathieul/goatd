package models

import (
    "goatd/app/event"
)

/*
 * Teammate
 */

type Teammate struct {
    Storage
    team *Team
    AttrName string
    AttrTeamUid string
}

func NewTeammate(attributes Attrs) *Teammate {
    return newModel(&Teammate{}, &attributes).(*Teammate)
}

func CreateTeammate(attributes Attrs, team *Team) *Teammate {
    teammate := NewTeammate(attributes)
    teammate.SetTeam(team)
    teammate.Save()
    return teammate
}

func (teammate Teammate) Name() string {
    return teammate.AttrName
}

func (teammate Teammate) TeamUid() string {
    return teammate.AttrTeamUid
}

func (teammate *Teammate) SetTeam(team *Team) {
    teammate.team = team
}

func (teammate Teammate) Team() (team *Team) {
    return teammate.team
}


/*
 * Teammates
 */

type Teammates struct {
    owner event.Identity
    items []*Teammate
}

func NewTeammates(owner event.Identity) (teammates *Teammates) {
    teammates = new(Teammates)
    teammates.owner = owner
    return teammates
}

func (teammates *Teammates) Create(attributes Attrs) (teammate *Teammate) {
    attributes = teammates.owner.AddToAttributes(attributes)
    teammate = CreateTeammate(attributes, teammates.owner.Value().(*Team))
    teammates.items = append(teammates.items, teammate)
    return teammate
}

func (teammates *Teammates) Slice() (slice []*Teammate) {
    slice = make([]*Teammate, len(teammates.items))
    copy(slice, teammates.items)
    return slice
}
