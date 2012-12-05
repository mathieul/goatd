package models

import (
    "goatd/app/identification"
    "github.com/sdegutis/fsm"
)

/*
 * Teammate
 */

type Teammate struct {
    Storage
    team *Team
    sm fsm.StateMachine
    AttrName string
    AttrTeamUid string
}

func setupTeammateStateMachine(teammate *Teammate) fsm.StateMachine {
    rules := []fsm.Rule{
        {From: "signed-out", Event: "sign-in", To: "on-break"},
        {From: "on-break", Event: "sign-out", To: "signed-out"},
    }
    sm := fsm.NewStateMachine(rules, teammate)
    return sm
}

func NewTeammate(attributes Attrs) *Teammate {
    teammate := newModel(&Teammate{}, &attributes).(*Teammate)
    teammate.sm = setupTeammateStateMachine(teammate)
    return teammate
}

func CreateTeammate(attributes Attrs) *Teammate {
    teammate := NewTeammate(attributes)
    teammate.Save()
    return teammate
}

func (teammate *Teammate) StateMachineCallback(action string, args []interface{}) {
    // TODO
}

func (teammate Teammate) Name() string { return teammate.AttrName }

func (teammate Teammate) TeamUid() string { return teammate.AttrTeamUid }

func (teammate *Teammate) SetTeam(team *Team) { teammate.team = team }

func (teammate Teammate) Team() (team *Team) { return teammate.team }

func (teammate Teammate) Status() Status { return statusFromString[teammate.sm.CurrentState] }

func (teammate *Teammate) SignIn() { teammate.sm.Process("sign-in") }

func (teammate *Teammate) SignOut() { teammate.sm.Process("sign-out") }


/*
 * Teammates
 */

type Teammates struct {
    Collection
}

func toTeammateSlice(source []interface{}) []*Teammate {
    teammates := make([]*Teammate, 0, len(source))
    for _, teammate := range source {
        teammates = append(teammates, teammate.(*Teammate))
    }
    return teammates
}

func NewTeammates(owner identification.Identity) (teammates *Teammates) {
    teammates = new(Teammates)
    teammates.Collection = NewCollection(func(attributes Attrs, owner interface{}) interface{} {
        teammate := CreateTeammate(attributes)
        teammate.SetTeam(owner.(*Team))
        return teammate
    }, owner)
    return teammates
}

func (teammates *Teammates) Create(attributes Attrs) (teammate *Teammate) {
    return teammates.Collection.Create(attributes).(*Teammate)
}

func (teammates Teammates) Slice() []*Teammate {
    return toTeammateSlice(teammates.Collection.Slice())
}

func (teammates Teammates) Find(uid string) *Teammate {
    if found := teammates.Collection.Find(uid); found != nil {
        return found.(*Teammate)
    }
    return nil
}

func (teammates Teammates) FindAll(uids []string) []*Teammate {
    return toTeammateSlice(teammates.Collection.FindAll(uids))
}

func (teammates Teammates) Select(tester func(interface{}) bool) (result []*Teammate) {
    return toTeammateSlice(teammates.Collection.Select(tester))
}
