package model

import (
    "github.com/sdegutis/fsm"
    "goatd/app/event"
)

/*
 * Task
 */

type Task struct {
    *event.Identity
    sm fsm.StateMachine
    AttrTitle string
    AttrTeamUid string
    AttrQueueUid string
}

func setupTaksStateMachine(task *Task) fsm.StateMachine {
    rules := []fsm.Rule{
        {From: "created", Event: "enqueue", To: "queued", Action: "setQueueUid"},
        {From: "queued", Event: "dequeue", To: "created", Action: "resetQueueUid"},
        {From: "queued", Event: "offer", To: "offered"},
        {From: "offered", Event: "assign", To: "assigned"},
        {From: "assigned", Event: "complete", To: "completed"},
    }
    sm := fsm.NewStateMachine(rules, task)
    return sm
}

func NewTask(attributes A) (task *Task) {
    task = newModel(&Task{}, &attributes).(*Task)
    task.Identity = event.NewIdentity("Task")
    task.sm = setupTaksStateMachine(task)
    return task
}

func (task *Task) Copy() Model {
    return &Task{task.Identity, task.sm, task.AttrTitle, task.AttrTeamUid, task.AttrQueueUid}
}

func (task *Task) StateMachineCallback(action string, args []interface{}) {
    switch action {
    case "setQueueUid":
        task.AttrQueueUid = args[0].(string)
    case "resetQueueUid":
        task.AttrQueueUid = ""
    }
}

func (task Task) Title() string { return task.AttrTitle }

func (task Task) TeamUid() string { return task.AttrTeamUid }

func (task Task) QueueUid() string { return task.AttrQueueUid }

// func (task Task) Status() Status { return statusFromString[task.sm.CurrentState] }

// func (task *Task) Enqueue(queue *Queue) bool {
//     if error := task.sm.Process("enqueue", queue.Uid()); error != nil { return false }
//     if !queue.InsertTask(task) {
//         task.sm.Process("dequeue", queue.Uid())
//         return false
//     }
//     return true
// }

// func (task *Task) Offer() bool {
//     if error := task.sm.Process("offer"); error != nil { return false }
//     return true
// }

// func (task *Task) Assign() bool {
//     if error := task.sm.Process("assign"); error != nil { return false }
//     return true
// }

// func (task *Task) Complete() bool {
//     if error := task.sm.Process("complete"); error != nil { return false }
//     return true
// }


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
