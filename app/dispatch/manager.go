package dispatch

import (
    "goatd/app/model"
)


/*
 * Manager
 */

type Manager struct {
    *model.Store
}

func (manager Manager) QueueTask(queue *model.Queue, task *model.Task) bool {
    if !queue.AddTask(task.Uid()) { return false }
    if task.Enqueue(queue.Uid())  { return true }
    queue.RemoveTask(task.Uid())
    return false
}

func (manager Manager) MakeTeammateAvailable(teammate *model.Teammate) bool {
    if !teammate.MakeAvailable() { return false }
    FindAndAssignTaskForTeammate(teammate)
    return true
}

func NewManager(store *model.Store) *Manager {
    return &Manager{store}
}
