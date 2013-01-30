package dispatch

import (
    "goatd/app/event"
    "goatd/app/model"
)


/*
 * Manager
 */

type Manager struct {
    busManager *event.BusManager
    store      *model.Store
}

func (manager Manager) QueueTask(queue *model.Queue, task *model.Task) bool {
    if !queue.AddTask(task.Uid()) { return false }
    if task.Enqueue(queue.Uid())  { return true }
    queue.DelTask(task.Uid())
    return false
}

func (manager Manager) MakeTeammateAvailable(teammate *model.Teammate) bool {
    if !teammate.MakeAvailable() { return false }
    manager.busManager.PublishEvent(event.TeammateAvailable,
        teammate.Identity, []interface{}{})
    if task := FindTaskForTeammate(manager.store, teammate); task != nil {
        // for now we just try offering this task: if it doesn't work we give up
        if !task.Offer(teammate.Uid()) { return true }
        if teammate.OfferTask(task.Uid()) {
            manager.busManager.PublishEvent(event.OfferTask, teammate.Identity,
                []interface{}{teammate.Uid(), task.Uid()})
        } else {
            task.Requeue()
        }
    }
    return true
}

func (manager Manager) AcceptTask(teammate *model.Teammate, task *model.Task) bool {
    if !task.Assign(teammate.Uid())    { return false }
    if teammate.AcceptTask(task.Uid()) {
        manager.busManager.PublishEvent(event.AcceptTask, teammate.Identity,
            []interface{}{teammate.Uid(), task.Uid()})
        return true
    }
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
    manager.busManager.PublishEvent(event.CompleteTask, teammate.Identity,
        []interface{}{teammate.Uid(), task.Uid()})
    return true
}

func NewManager(busManager *event.BusManager, store *model.Store) *Manager {
    return &Manager{busManager, store}
}
