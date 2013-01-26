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
    status sm.Status
    AttrName string
    AttrTeamUid string
    AttrTaskUid string
}

func setupTeammateStateMachine(teammate *Teammate, status sm.Status) *sm.StateMachine {
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
        b.Event(EventMakeAvailable, StatusOnBreak, StatusWaiting, func (args []interface{}) bool {
            teammate := args[0].(*Teammate)
            teammate.busManager.PublishEvent(event.TeammateAvailable,
                teammate.Identity, []interface{}{})
            return true
        })
        b.Event(EventOfferTask, StatusWaiting, StatusOffered, func (args []interface{}) bool {
            teammate, taskUid := args[0].(*Teammate), args[1].(string)
            teammate.AttrTaskUid = taskUid
            return true
        })
        b.Event(EventAcceptTask, StatusOffered, StatusBusy, func (args []interface{}) bool {
            teammate, task := args[0].(*Teammate), args[1].(*Task)
            if task.Uid() != teammate.AttrTaskUid { return false }
            teammate.busManager.PublishEvent(event.AcceptTask, teammate.Identity,
                []interface{}{teammate.Uid(), task.Uid()})
            return true
        })
        b.Event(EventRejectTask, StatusOffered, StatusWaiting, func (args []interface{}) bool {
            teammate, task := args[0].(*Teammate), args[1].(*Task)
            if task.Uid() != teammate.AttrTaskUid { return false }
            teammate.AttrTaskUid = ""
            teammate.busManager.PublishEvent(event.CompleteTask, teammate.Identity,
                []interface{}{teammate.Uid(), task.Uid()})
            return true
        })
        b.Event(EventFinishTask, StatusBusy, StatusWrappingUp, func (args []interface{}) bool {
            teammate, task := args[0].(*Teammate), args[1].(*Task)
            if task.Uid() != teammate.AttrTaskUid { return false }
            teammate.AttrTaskUid = ""
            teammate.busManager.PublishEvent(event.CompleteTask, teammate.Identity,
                []interface{}{teammate.Uid(), task.Uid()})
            return true
        })
        b.Event(EventStartOtherWork, func (b sm.Builder) {
            b.Transition(StatusOnBreak, StatusOtherWork, sm.NoAction)
            b.Transition(StatusWaiting, StatusOtherWork, sm.NoAction)
            b.Transition(StatusWrappingUp, StatusOtherWork, sm.NoAction)
        })
    })
    return stateMachine
}

func NewTeammate(attributes A) (teammate *Teammate) {
    teammate = newModel(&Teammate{}, &attributes).(*Teammate)
    if teammate.status == StatusNone { teammate.status = StatusSignedOut }
    teammate.Identity = event.NewIdentity("Teammate")
    return teammate
}

func (teammate *Teammate) Copy() Model {
    stateMachine := setupTeammateStateMachine(teammate, teammate.status)
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

func (teammate Teammate) Name() string { return teammate.AttrName }

func (teammate Teammate) TeamUid() string { return teammate.AttrTeamUid }

func (teammate Teammate) TaskUid() string { return teammate.AttrTaskUid }

func (teammate Teammate) Status() sm.Status {
    if teammate.IsCopy() && teammate.stateMachine != nil {
        return teammate.stateMachine.Status()
    }
    return teammate.status    
}

func (teammate *Teammate) SignIn() bool {
    return teammate.stateMachine.Trigger(EventSignIn)
}

func (teammate *Teammate) GoOnBreak() bool {
    return teammate.stateMachine.Trigger(EventGoOnBreak)
}

func (teammate *Teammate) MakeAvailable() bool {
    return teammate.stateMachine.Trigger(EventMakeAvailable, teammate)
}

func (teammate *Teammate) OfferTask(task *Task) bool {
    return teammate.stateMachine.Trigger(EventOfferTask, teammate, task.Uid())
}

func (teammate *Teammate) AcceptTask(task *Task) bool {
    return teammate.stateMachine.Trigger(EventAcceptTask, teammate, task)
}

func (teammate *Teammate) RejectTask(task *Task) bool {
    return teammate.stateMachine.Trigger(EventRejectTask, teammate, task)
}

func (teammate *Teammate) FinishTask(task *Task) bool {
    return teammate.stateMachine.Trigger(EventFinishTask, teammate, task)
}

func (teammate *Teammate) StartOtherWork() bool {
    return teammate.stateMachine.Trigger(EventStartOtherWork)
}

func (teammate *Teammate) SignOut() bool {
    return teammate.stateMachine.Trigger(EventSignOut)
}

func (teammate Teammate) CurrentTask() *Task {
    if teammate.AttrTaskUid == "" { return nil }
    found := teammate.store.Tasks.Select(func (item interface{}) bool {
        return item.(*Task).Uid() == teammate.AttrTaskUid
    })
    if len(found) == 0 { return nil }
    return found[0]
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
