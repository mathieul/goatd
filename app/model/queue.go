package model

import (
    "goatd/app/event"
    "goatd/app/sm"
)

/*
 * Global
 */

const (
    initialQueueLength = 10
)

/*
 * Queue
 */

type Queue struct {
    *event.Identity
    busManager *event.BusManager
    store *Store
    queuedTasks []*Task
    isReady bool
    AttrName string
    AttrTeamUid string
}

func NewQueue(attributes A) (queue *Queue) {
    queue = newModel(&Queue{}, &attributes).(*Queue)
    queue.Identity = event.NewIdentity("Queue")
    queue.queuedTasks = make([]*Task, 0, initialQueueLength)
    queue.isReady = false
    return queue
}

func (queue *Queue) Copy() Model {
    identity := queue.Identity.Copy()
    var queuedTasks []*Task
    copy(queuedTasks, queue.queuedTasks)
    return &Queue{identity, nil, nil, queuedTasks, queue.isReady,
        queue.AttrName, queue.AttrTeamUid}
}

func (queue *Queue) SetupComs(busManager *event.BusManager, store *Store) {
    queue.busManager = busManager
    queue.store = store
}

func (queue *Queue) Update(name string, value interface{}) bool {
    setAttributeValue(queue, name, value)
    return queue.store.Update(KindQueue, queue.Uid(), name, value)
}

func (queue Queue) Status(_ ...sm.Status) sm.Status {
    return StatusNone
}

func (queue Queue) Name() string {
    return queue.AttrName
}

func (queue Queue) TeamUid() string {
    return queue.AttrTeamUid
}

func (queue Queue) IsReady() bool {
    return queue.isReady
}

func (queue Queue) QueuedTasks() []*Task {
    return queue.queuedTasks
}

// func (queue *Queue) InsertTask(task *Task) bool {
//     queue.queuedTasks = append(queue.queuedTasks, task)
//     queue.isReady = true
//     return true
// }

// func (queue *Queue) RemoveTask(taskToRemove *Task) bool {
//     index, uid := -1, taskToRemove.Uid()
//     for i, task := range queue.queuedTasks {
//         if task.Uid() == uid {
//             index = i
//             break
//         }
//     }
//     queue.queuedTasks = append(queue.queuedTasks[:index], queue.queuedTasks[index + 1:]...)
//     if len(queue.queuedTasks) == 0 {
//         queue.isReady = false
//     }
//     return true
// }


/*
 * QueueStoreProxy
 */

type QueueStoreProxy struct {
    store *Store
}

func toQueueSlice(source []Model) []*Queue {
    queues := make([]*Queue, 0, len(source))
    for _, queue := range source {
        queues = append(queues, queue.(*Queue))
    }
    return queues
}

func (proxy *QueueStoreProxy) Create(attributes A, owners ...event.Identified) *Queue {
    for _, owner := range owners { attributes = owner.AddToAttributes(attributes) }
    return proxy.store.Create(KindQueue, attributes).(*Queue)
}

func (proxy *QueueStoreProxy) Find(uid string) *Queue {
    if value := proxy.store.Find(KindQueue, uid); value != nil { return value.(*Queue) }
    return nil
}

func (proxy *QueueStoreProxy) FindAll(uids []string) []*Queue {
    values := proxy.store.FindAll(KindQueue, uids)
    return toQueueSlice(values)
}

func (proxy *QueueStoreProxy) Select(tester func(interface{}) bool) []*Queue {
    values := proxy.store.Select(KindQueue, tester)
    return toQueueSlice(values)
}
