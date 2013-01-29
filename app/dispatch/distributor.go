package dispatch

import (
    "goatd/app/model"
)


/*
 * Distributor
 */

type Distributor struct {
    *model.Store
}

func (distributor Distributor) QueueTask(queue *model.Queue, task *model.Task) bool {
    if queue.AddTask(task.Uid()) {
        if !task.Enqueue(queue.Uid()) {
            queue.RemoveTask(task.Uid())
            return false
        }
        return true
    }
    return false
}

func NewDistributor(store *model.Store) *Distributor {
    return &Distributor{store}
}
