package model

import (
    // "strings"
    // "fmt"
    "github.com/sdegutis/fsm"
    "goatd/app/event"
)

/*
 * Teammate
 */

type Teammate struct {
    *event.Identity
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
        {From: "offered", Event: "accept-task", To: "busy", Action: "publishAcceptTask"},
        {From: "offered", Event: "reject-task", To: "waiting", Action: "resetTaskUid"},
        {From: "busy", Event: "finish-task", To: "wrapping-up", Action: "resetTaskUid"},

        {From: "on-break", Event: "start-other-work", To: "other-work"},
        {From: "waiting", Event: "start-other-work", To: "other-work"},
        {From: "wrapping-up", Event: "start-other-work", To: "other-work"},
    }
    sm := fsm.NewStateMachine(rules, teammate)
    return sm
}

func NewTeammate(attributes A) (teammate *Teammate) {
    teammate = newModel(&Teammate{}, &attributes).(*Teammate)
    teammate.Identity = event.NewIdentity("Teammate")
    teammate.sm = setupTeammateStateMachine(teammate)
    return teammate
}

func (teammate *Teammate) Copy() Model {
    return &Teammate{teammate.Identity, teammate.sm, teammate.AttrName,
        teammate.AttrTeamUid, teammate.AttrTaskUid}
}

func (teammate *Teammate) StateMachineCallback(action string, args []interface{}) {
    // switch action {
    // case "setTaskUid":
    //     teammate.AttrTaskUid = args[0].(string)
    // case "resetTaskUid":
    //     teammate.AttrTaskUid = ""
    //     event.Manager().PublishEvent(event.KindCompleteTask, *teammate.identity,
    //         []string{teammate.Uid(), args[0].(*Task).Uid()})
    // case "publishWaiting":
    //     event.Manager().PublishEvent(event.KindTeammateAvailable, *teammate.identity, nil)
    // case "publishAcceptTask":
    //     event.Manager().PublishEvent(event.KindAcceptTask, *teammate.identity,
    //         []string{teammate.Uid(), args[0].(*Task).Uid()})

    // }
}

func (teammate Teammate) Name() string { return teammate.AttrName }

func (teammate Teammate) TeamUid() string { return teammate.AttrTeamUid }

// func (teammate Teammate) Status() Status { return statusFromString[teammate.sm.CurrentState] }

// func (teammate *Teammate) SignIn() bool {
//     if error := teammate.sm.Process("sign-in"); error != nil { return false }
//     return true
// }

// func (teammate *Teammate) GoOnBreak() bool {
//     if error := teammate.sm.Process("go-on-break"); error != nil { return false }
//     return true
// }

// func (teammate *Teammate) MakeAvailable() bool {
//     if error := teammate.sm.Process("make-available"); error != nil { return false }
//     return true
// }

// func (teammate *Teammate) OfferTask(task *Task) bool {
//     if error := teammate.sm.Process("offer-task", task.Uid()); error != nil { return false }
//     return true
// }

// func (teammate *Teammate) AcceptTask(task *Task) bool {
//     if task.Uid() != teammate.AttrTaskUid { return false }
//     if error := teammate.sm.Process("accept-task", task); error != nil { return false }
//     return true
// }

// func (teammate *Teammate) RejectTask(task *Task) bool {
//     if task.Uid() != teammate.AttrTaskUid { return false }
//     if error := teammate.sm.Process("reject-task"); error != nil { return false }
//     return true
// }

// func (teammate *Teammate) FinishTask(task *Task) bool {
//     if task.Uid() != teammate.AttrTaskUid { return false }
//     if error := teammate.sm.Process("finish-task", task); error != nil { return false }
//     return true
// }

// func (teammate *Teammate) StartOtherWork() bool {
//     if error := teammate.sm.Process("start-other-work"); error != nil { return false }
//     return true
// }

// func (teammate *Teammate) SignOut() bool {
//     if error := teammate.sm.Process("sign-out"); error != nil { return false }
//     return true
// }

// func (teammate Teammate) CurrentTask() *Task {
//     if teammate.AttrTaskUid == "" { return nil }
//     found := teammate.team.Tasks.Select(func (item interface{}) bool {
//         task := item.(*Task)
//         return strings.Contains(task.Uid(), teammate.AttrTaskUid)
//     })
//     if len(found) == 0 { return nil }
//     return found[0]
// }


/*
 * TeammateStoreProxy
 */

type TeammateStoreProxy struct {
    store *Store
}

func toTeammateSlice(source []Model) []*Teammate {
    teammates := make([]*Teammate, 0, len(source))
    for _, teammate := range source {
        teammates = append(teammates, teammate.(*Teammate))
    }
    return teammates
}

func (proxy *TeammateStoreProxy) Create(attributes A, owners ...event.Identified) *Teammate {
    for _, owner := range owners { attributes = owner.AddToAttributes(attributes) }
    return proxy.store.Create(KindTeammate, attributes).(*Teammate)
}

func (proxy *TeammateStoreProxy) Find(uid string) *Teammate {
    if value := proxy.store.Find(KindTeammate, uid); value != nil { return value.(*Teammate) }
    return nil
}

func (proxy *TeammateStoreProxy) FindAll(uids []string) []*Teammate {
    values := proxy.store.FindAll(KindTeammate, uids)
    return toTeammateSlice(values)
}

func (proxy *TeammateStoreProxy) Select(tester func(interface{}) bool) []*Teammate {
    values := proxy.store.Select(KindTeammate, tester)
    return toTeammateSlice(values)
}
