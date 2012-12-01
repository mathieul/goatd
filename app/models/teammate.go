package models

import (
	"goatd/app/event"
)

/*
 * Teammate
 */

type Teammate struct {
	Storage
    AttrName string
    AttrTeamUid string
}

func NewTeammate(attributes Attrs) *Teammate {
	return newModel(&Teammate{}, &attributes).(*Teammate)
}

func CreateTeammate(attributes Attrs) (teammate *Teammate) {
	teammate = NewTeammate(attributes)
	teammate.Save()
	return teammate
}

func (team *Teammate) Name() string {
    return team.AttrName
}

func (team *Teammate) TeamUid() string {
    return team.AttrTeamUid
}

/*
 * Teammates
 */

type Teammates struct {
	owner *event.Identity
	items []*Teammate
}

func NewTeammates(kind, uid string) (teammates *Teammates) {
	teammates = new(Teammates)
	teammates.owner = event.NewIdentity(kind, uid)
	return teammates
}

func (teammates *Teammates) Create(attributes Attrs) (teammate *Teammate) {
	teammate = CreateTeammate(teammates.owner.AddToAttributes(attributes))
	teammates.items = append(teammates.items, teammate)
	return teammate
}
