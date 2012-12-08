package models

import (
    "goatd/app/identification"
    "github.com/sdegutis/fsm"
)

/*
 * Task
 */

type Task struct {
    Storage
    team *Team
    sm fsm.StateMachine
    AttrTitle string
    AttrQueueUid string
    AttrPriority int
}

func setupTaksStateMachine(task *Task) fsm.StateMachine {
    rules := []fsm.Rule{
        {From: "created", Event: "enqueue", To: "queued", Action: "setQueueUid"},
        {From: "queued", Event: "dequeue", To: "created", Action: "resetQueueUid"},
        {From: "queued", Event: "offer", To: "offered"},
    }
    sm := fsm.NewStateMachine(rules, task)
    return sm
}

func NewTask(attributes Attrs) (task *Task) {
    task = newModel(&Task{}, &attributes).(*Task)
    if task.AttrPriority == PriorityNone { task.AttrPriority = PriorityMedium }
    task.sm = setupTaksStateMachine(task)
    return task
}

func CreateTask(attributes Attrs) (task *Task) {
    task = NewTask(attributes)
    task.Save()
    return task
}

func (task *Task) StateMachineCallback(action string, args []interface{}) {
    switch action {
    case "setQueueUid":
        task.AttrQueueUid = args[0].(string)
    case "resetQueueUid":
        task.AttrQueueUid = ""
    }
}

func (team Task) Title() string { return team.AttrTitle }

func (team Task) QueueUid() string { return team.AttrQueueUid }

func (team *Task) SetPriority(priority int) { team.AttrPriority = priority }

func (team Task) Priority() int { return team.AttrPriority }

func (task *Task) SetTeam(team *Team) { task.team = team }

func (task Task) Team() (team *Team) { return task.team }

func (task Task) Status() Status { return statusFromString[task.sm.CurrentState] }

func (task *Task) Enqueue(queue *Queue) bool {
    if error := task.sm.Process("enqueue", queue.Uid()); error != nil {
        return false
    }
    if !queue.InsertTask(task) {
        task.sm.Process("dequeue", queue.Uid())
        return false
    }
    return true
}

func (task *Task) Offer() bool {
    if error := task.sm.Process("offer"); error != nil {
        return false
    }
    return true
}


/*
 * Tasks
 */

type Tasks struct {
    Collection
}

func toTaskSlice(source []interface{}) []*Task {
    tasks := make([]*Task, 0, len(source))
    for _, task := range source {
        tasks = append(tasks, task.(*Task))
    }
    return tasks
}

func NewTasks(owner identification.Identity) (tasks *Tasks) {
    tasks = new(Tasks)
    tasks.Collection = NewCollection(func(attributes Attrs, parent interface{}) interface{} {
        task := CreateTask(attributes)
        task.SetTeam(parent.(*Team))
        return task
    }, owner)
    return tasks
}

func (tasks *Tasks) Create(attributes Attrs) (task *Task) {
    return tasks.Collection.Create(attributes).(*Task)
}

func (tasks Tasks) Slice() []*Task {
    return toTaskSlice(tasks.Collection.Slice())
}

func (tasks Tasks) Find(uid string) *Task {
    if found := tasks.Collection.Find(uid); found != nil {
        return found.(*Task)
    }
    return nil
}

func (tasks Tasks) FindAll(uids []string) []*Task {
    return toTaskSlice(tasks.Collection.FindAll(uids))
}

func (tasks Tasks) Select(tester func(interface{}) bool) []*Task {
    return toTaskSlice(tasks.Collection.Select(tester))
}
