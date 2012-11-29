package models

type Task struct {
    Storage
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
