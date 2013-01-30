package dispatch

import (
    "goatd/app/model"
)


/*
 * Manager
 */

type Manager struct {
    store *model.Store
}

func (manager Manager) QueueTask(queue *model.Queue, task *model.Task) bool {
    if !queue.AddTask(task.Uid()) { return false }
    if task.Enqueue(queue.Uid())  { return true }
    queue.DelTask(task.Uid())
    return false
}

func (manager Manager) MakeTeammateAvailable(teammate *model.Teammate) bool {
    if !teammate.MakeAvailable() { return false }
    if task := FindTaskForTeammate(manager.store, teammate); task != nil {
        // for now we just try offering this task: if it doesn't work we give up
        if !task.Offer(teammate.Uid()) { return true }
        if !teammate.OfferTask(task.Uid()) {
            task.Requeue()
        }
    }
    return true
}

func (manager Manager) AcceptTask(teammate *model.Teammate, task *model.Task) bool {
    if !task.Assign(teammate.Uid())    { return false }
    if teammate.AcceptTask(task.Uid()) { return true }
    task.Requeue()
    return false
}

func (manager Manager) FinishTask(teammate *model.Teammate, task *model.Task) bool {
    if task.TeammateUid() != teammate.Uid() ||
        teammate.CurrentTask() == nil ||
        teammate.CurrentTask().Uid() != task.Uid() { return false }
    if queue := manager.store.Queues.Find(task.QueueUid()); queue != nil {
        queue.DelTask(task.Uid())
    }
    if !task.Complete() { return false }
    teammate.FinishTask(task.Uid())
    return true
}

func NewManager(store *model.Store) *Manager {
    return &Manager{store}
}
