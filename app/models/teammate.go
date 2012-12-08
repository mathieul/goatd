package models

import (
    "strings"
    "github.com/sdegutis/fsm"
    "goatd/app/identification"
    "goatd/app/event"
)

/*
 * Teammate
 */

type Teammate struct {
    Storage
    identity *identification.Identity
    team *Team
    sm fsm.StateMachine
    AttrName string
    AttrTeamUid string
    AttrTaskUid string
}

func setupTeammateStateMachine(teammate *Teammate) fsm.StateMachine {
    rules := []fsm.Rule{
        {From: "signed-out", Event: "sign-in", To: "on-break"},
        {From: "waiting", Event: "go-on-break", To: "on-break"},
        {From: "wrapping-up", Event: "go-on-break", To: "on-break"},
        {From: "offered", Event: "go-on-break", To: "on-break"},
        {From: "other-work", Event: "go-on-break", To: "on-break"},

        {From: "on-break", Event: "sign-out", To: "signed-out"},
        {From: "waiting", Event: "sign-out", To: "signed-out"},

        {From: "on-break", Event: "make-available", To: "waiting", Action: "publishWaiting"},
        {From: "waiting", Event: "offer-task", To: "offered", Action: "setTaskUid"},
        {From: "offered", Event: "accept-task", To: "busy"},
        {From: "offered", Event: "reject-task", To: "waiting", Action: "resetTaskUid"},
        {From: "busy", Event: "finish-task", To: "wrapping-up", Action: "resetTaskUid"},

        {From: "on-break", Event: "start-other-work", To: "other-work"},
        {From: "waiting", Event: "start-other-work", To: "other-work"},
        {From: "wrapping-up", Event: "start-other-work", To: "other-work"},
    }
    sm := fsm.NewStateMachine(rules, teammate)
    return sm
}

func NewTeammate(attributes Attrs) *Teammate {
    teammate := newModel(&Teammate{}, &attributes).(*Teammate)
    teammate.identity = identification.NewIdentity("Teammate", teammate.Uid(), teammate)
    teammate.sm = setupTeammateStateMachine(teammate)
    return teammate
}

func CreateTeammate(attributes Attrs) *Teammate {
    teammate := NewTeammate(attributes)
    teammate.Save()
    return teammate
}

func (teammate *Teammate) StateMachineCallback(action string, args []interface{}) {
    switch action {
    case "setTaskUid":
        teammate.AttrTaskUid = args[0].(string)
    case "resetTaskUid":
        teammate.AttrTaskUid = ""
    case "publishWaiting":
        event.Manager().PublishEvent(event.KindTeammateAvailable, *teammate.identity, nil)
    }
}

func (teammate Teammate) Name() string { return teammate.AttrName }

func (teammate Teammate) TeamUid() string { return teammate.AttrTeamUid }

func (teammate *Teammate) SetTeam(team *Team) { teammate.team = team }

func (teammate Teammate) Team() (team *Team) { return teammate.team }

func (teammate Teammate) Status() Status { return statusFromString[teammate.sm.CurrentState] }

func (teammate *Teammate) SignIn() bool {
    if error := teammate.sm.Process("sign-in"); error != nil { return false }
    return true
}

func (teammate *Teammate) GoOnBreak() bool {
    if error := teammate.sm.Process("go-on-break"); error != nil { return false }
    return true
}

func (teammate *Teammate) MakeAvailable() bool {
    if error := teammate.sm.Process("make-available"); error != nil { return false }
    return true
}

func (teammate *Teammate) OfferTask(task *Task) bool {
    if error := teammate.sm.Process("offer-task", task.Uid()); error != nil { return false }
    task.Offer()
    return true
}

func (teammate *Teammate) AcceptTask(task *Task) bool {
    if task.Uid() != teammate.AttrTaskUid { return false }
    if error := teammate.sm.Process("accept-task"); error != nil { return false }
    return true
}

func (teammate *Teammate) RejectTask(task *Task) bool {
    if task.Uid() != teammate.AttrTaskUid { return false }
    if error := teammate.sm.Process("reject-task"); error != nil { return false }
    return true
}

func (teammate *Teammate) FinishTask(task *Task) bool {
    if task.Uid() != teammate.AttrTaskUid { return false }
    if error := teammate.sm.Process("finish-task"); error != nil { return false }
    return true
}

func (teammate *Teammate) StartOtherWork() bool {
    if error := teammate.sm.Process("start-other-work"); error != nil { return false }
    return true
}

func (teammate *Teammate) SignOut() bool {
    if error := teammate.sm.Process("sign-out"); error != nil { return false }
    return true
}

func (teammate Teammate) CurrentTask() *Task {
    if teammate.AttrTaskUid == "" { return nil }
    found := teammate.team.Tasks.Select(func (item interface{}) bool {
        task := item.(*Task)
        return strings.Contains(task.Uid(), teammate.AttrTaskUid)
    })
    if len(found) == 0 { return nil }
    return found[0]
}


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
