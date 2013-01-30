package model

import (
    "goatd/app/event"
    "goatd/app/sm"
)

/*
 * Teammate
 */

type Teammate struct {
    *event.Identity
    busManager *event.BusManager
    store *Store
    stateMachine *sm.StateMachine
    InternalStatus sm.Status
    AttrName string
    AttrTeamUid string
    AttrTaskUid string
}

func setupTeammateStateMachine(teammate *Teammate, status sm.Status) *sm.StateMachine {
    // trigger arguments: <teammate>, [taskUid]
    stateMachine := sm.NewStateMachine(status, func (b sm.Builder) {
        b.Event(EventSignIn, StatusSignedOut, StatusOnBreak, sm.NoAction)
        b.Event(EventGoOnBreak, func (b sm.Builder) {
            b.Transition(StatusWaiting, StatusOnBreak, sm.NoAction)
            b.Transition(StatusWrappingUp, StatusOnBreak, sm.NoAction)
            b.Transition(StatusOffered, StatusOnBreak, sm.NoAction)
            b.Transition(StatusOtherWork, StatusOnBreak, sm.NoAction)
        })
        b.Event(EventSignOut, func (b sm.Builder) {
            b.Transition(StatusOnBreak, StatusSignedOut, sm.NoAction)
            b.Transition(StatusWaiting, StatusSignedOut, sm.NoAction)
        })
        b.Event(EventMakeAvailable, StatusOnBreak, StatusWaiting, sm.NoAction)
        b.Event(EventOfferTask, StatusWaiting, StatusOffered, func (args []interface{}) bool {
            teammate, taskUid := args[0].(*Teammate), args[1].(string)
            teammate.Update("TaskUid", taskUid)
            return true
        })
        b.Event(EventAcceptTask, StatusOffered, StatusBusy, func (args []interface{}) bool {
            teammate, taskUid := args[0].(*Teammate), args[1].(string)
            if taskUid != teammate.TaskUid() { return false }
            return true
        })
        wrapupTask := func (args []interface{}) bool {
            teammate, taskUid := args[0].(*Teammate), args[1].(string)
            if taskUid != teammate.TaskUid() { return false }
            teammate.Update("TaskUid", "")
            return true
        }
        b.Event(EventRejectTask, StatusOffered, StatusWaiting, wrapupTask)
        b.Event(EventFinishTask, StatusBusy, StatusWrappingUp, wrapupTask)
        b.Event(EventStartOtherWork, func (b sm.Builder) {
            b.Transition(StatusOnBreak, StatusOtherWork, sm.NoAction)
            b.Transition(StatusWaiting, StatusOtherWork, sm.NoAction)
            b.Transition(StatusWrappingUp, StatusOtherWork, sm.NoAction)
        })
    })
    stateMachine.SetTriggerValidator(func (oldStatus, newStatus sm.Status, args ...interface{}) bool {
        teammate := args[0].(*Teammate)
        accepted := teammate.store.SetStatus(KindTeammate, teammate.Uid(), oldStatus, newStatus)
        return accepted
    })
    return stateMachine
}

func NewTeammate(attributes A) (teammate *Teammate) {
    teammate = newModel(&Teammate{}, &attributes).(*Teammate)
    if teammate.InternalStatus == StatusNone { teammate.InternalStatus = StatusSignedOut }
    teammate.Identity = event.NewIdentity("Teammate")
    return teammate
}

func (teammate *Teammate) Copy() Model {
    stateMachine := setupTeammateStateMachine(teammate, teammate.InternalStatus)
    identity := teammate.Identity.Copy()
    return &Teammate{identity, nil, nil, stateMachine, StatusNone,
        teammate.AttrName, teammate.AttrTeamUid, teammate.AttrTaskUid}
}

func (teammate *Teammate) SetupComs(busManager *event.BusManager, store *Store) {
    teammate.busManager = busManager
    teammate.store = store
}

func (teammate *Teammate) Update(name string, value interface{}) bool {
    setAttributeValue(teammate, name, value)
    return teammate.store.Update(KindTeammate, teammate.Uid(), name, value)
}

func (teammate Teammate) Reload() *Teammate {
    if found := teammate.store.Teammates.Find(teammate.Uid()); found != nil {
        return found
    }
    return nil
}

func (teammate Teammate) Name() string { return teammate.AttrName }

func (teammate Teammate) TeamUid() string { return teammate.AttrTeamUid }

func (teammate Teammate) TaskUid() string { return teammate.AttrTaskUid }

func (teammate *Teammate) Status(newStatus ...sm.Status) sm.Status {
    if len(newStatus) > 0 {
        teammate.InternalStatus = newStatus[0]
    }
    if teammate.IsCopy() && teammate.stateMachine != nil {
        return teammate.stateMachine.Status()
    }
    return teammate.InternalStatus
}

func (teammate *Teammate) SignIn() bool {
    return teammate.stateMachine.Trigger(EventSignIn, teammate)
}

func (teammate *Teammate) GoOnBreak() bool {
    return teammate.stateMachine.Trigger(EventGoOnBreak, teammate)
}

func (teammate *Teammate) MakeAvailable() bool {
    return teammate.stateMachine.Trigger(EventMakeAvailable, teammate)
}

func (teammate *Teammate) OfferTask(taskUid string) bool {
    return teammate.stateMachine.Trigger(EventOfferTask, teammate, taskUid)
}

func (teammate *Teammate) AcceptTask(taskUid string) bool {
    return teammate.stateMachine.Trigger(EventAcceptTask, teammate, taskUid)
}

func (teammate *Teammate) RejectTask(taskUid string) bool {
    return teammate.stateMachine.Trigger(EventRejectTask, teammate, taskUid)
}

func (teammate *Teammate) FinishTask(taskUid string) bool {
    return teammate.stateMachine.Trigger(EventFinishTask, teammate, taskUid)
}

func (teammate *Teammate) StartOtherWork() bool {
    return teammate.stateMachine.Trigger(EventStartOtherWork, teammate)
}

func (teammate *Teammate) SignOut() bool {
    return teammate.stateMachine.Trigger(EventSignOut, teammate)
}


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
