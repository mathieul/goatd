package models

import (
    "goatd/app/identification"
)

/*
 * Task
 */

type Task struct {
    Storage
    team *Team
    AttrTitle string
}

func NewTask(attributes Attrs) *Task {
    return newModel(&Task{}, &attributes).(*Task)
}

func CreateTask(attributes Attrs) (task *Task) {
    task = NewTask(attributes)
    task.Save()
    return task
}

func (team *Task) Title() string {
    return team.AttrTitle
}

func (task *Task) SetTeam(team *Team) {
    task.team = team
}

func (task Task) Team() (team *Team) {
    return task.team
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

func (tasks Tasks) Select(tester func(interface{}) bool) (result []*Task) {
    return toTaskSlice(tasks.Collection.Select(tester))
}
