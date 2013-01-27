package model

import (
    "goatd/app/event"
    "goatd/app/sm"
)

/*
 * Task
 */

type Task struct {
    *event.Identity
    busManager *event.BusManager
    store *Store
    stateMachine *sm.StateMachine
    InternalStatus sm.Status
    AttrTitle string
    AttrTeamUid string
    AttrQueueUid string
}

func setupTaksStateMachine(task *Task, status sm.Status) *sm.StateMachine {
    stateMachine := sm.NewStateMachine(status, func (b sm.Builder) {
        b.Event(EventEnqueue, StatusCreated, StatusQueued, func (args []interface{}) bool {
            task, queueUid := args[0].(*Task), args[1].(string)
            task.Update("QueueUid", queueUid)
            return true
        })
        b.Event(EventDequeue, StatusQueued, StatusCreated, func (args []interface{}) bool {
            task, queueUid := args[0].(*Task), args[1].(string)
            if queueUid != task.QueueUid() { return false }
            task.Update("QueueUid", "")
            return true
        })
        b.Event(EventOffer, StatusQueued, StatusOffered, sm.NoAction)
        b.Event(EventAssign, StatusOffered, StatusAssigned, sm.NoAction)
        b.Event(EventComplete, StatusAssigned, StatusCompleted, sm.NoAction)
    })
    stateMachine.SetTriggerValidator(func (oldStatus, newStatus sm.Status, args ...interface{}) bool {
        task := args[0].(*Task)
        accepted := task.store.SetStatus(KindTask, task.Uid(), oldStatus, newStatus)
        return accepted
    })
    return stateMachine
}

func NewTask(attributes A) (task *Task) {
    task = newModel(&Task{}, &attributes).(*Task)
    if task.InternalStatus == StatusNone { task.InternalStatus = StatusCreated }
    task.Identity = event.NewIdentity("Task")
    return task
}

func (task *Task) Copy() Model {
    stateMachine := setupTaksStateMachine(task, task.InternalStatus)
    identity := task.Identity.Copy()
    return &Task{identity, nil, nil, stateMachine, task.InternalStatus,
        task.AttrTitle, task.AttrTeamUid, task.AttrQueueUid}
}

func (task *Task) SetupComs(busManager *event.BusManager, store *Store) {
    task.busManager = busManager
    task.store = store
}

func (task *Task) Update(name string, value interface{}) bool {
    setAttributeValue(task, name, value)
    return task.store.Update(KindTask, task.Uid(), name, value)
}

func (task Task) Title() string { return task.AttrTitle }

func (task Task) TeamUid() string { return task.AttrTeamUid }

func (task Task) QueueUid() string { return task.AttrQueueUid }

func (task *Task) Status(newStatus ...sm.Status) sm.Status {
    if len(newStatus) > 0 {
        task.InternalStatus = newStatus[0]
    }
    if task.IsCopy() && task.stateMachine != nil {
        return task.stateMachine.Status()
    }
    return task.InternalStatus    
}

func (task *Task) Enqueue(queueUid string) bool {
    return task.stateMachine.Trigger(EventEnqueue, task, queueUid)
}

func (task *Task) Dequeue(queueUid string) bool {
    return task.stateMachine.Trigger(EventDequeue, task, queueUid)
}


/*
 * TaskStoreProxy
 */

type TaskStoreProxy struct {
    store *Store
}

func toTaskSlice(source []Model) []*Task {
    tasks := make([]*Task, 0, len(source))
    for _, task := range source {
        tasks = append(tasks, task.(*Task))
    }
    return tasks
}

func (proxy *TaskStoreProxy) Create(attributes A, owners ...event.Identified) *Task {
    for _, owner := range owners { attributes = owner.AddToAttributes(attributes) }
    return proxy.store.Create(KindTask, attributes).(*Task)
}

func (proxy *TaskStoreProxy) Find(uid string) *Task {
    if value := proxy.store.Find(KindTask, uid); value != nil { return value.(*Task) }
    return nil
}

func (proxy *TaskStoreProxy) FindAll(uids []string) []*Task {
    values := proxy.store.FindAll(KindTask, uids)
    return toTaskSlice(values)
}

func (proxy *TaskStoreProxy) Select(tester func(interface{}) bool) []*Task {
    values := proxy.store.Select(KindTask, tester)
    return toTaskSlice(values)
}
