package model

import (
    "github.com/petar/GoLLRB/llrb"
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
    Tasks *llrb.Tree
    numberTasks int
    nextTaskUid string
    AttrName string
    AttrTeamUid string
}

func NewQueue(attributes A) (queue *Queue) {
    queue = newModel(&Queue{}, &attributes).(*Queue)
    queue.Identity = event.NewIdentity("Queue")
    queue.Tasks = llrb.New(TaskLess)
    return queue
}

func (queue *Queue) Copy() Model {
    identity := queue.Identity.Copy()
    return &Queue{identity, nil, nil, nil, queue.NumberTasks(),
        queue.NextTaskUid(), queue.AttrName, queue.AttrTeamUid}
}

func (queue *Queue) SetupComs(busManager *event.BusManager, store *Store) {
    queue.busManager = busManager
    queue.store = store
}

func (queue *Queue) Update(name string, value interface{}) bool {
    setAttributeValue(queue, name, value)
    return queue.store.Update(KindQueue, queue.Uid(), name, value)
}

func (queue Queue) Reload() *Queue {
    if found := queue.store.Queues.Find(queue.Uid()); found != nil {
        return found
    }
    return nil
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

func (queue Queue) NextTaskUid() string {
    if queue.Tasks == nil {
        return queue.nextTaskUid
    }
    if item := queue.Tasks.Min(); item != nil {
        return item.(*Task).Uid()
    }
    return ""
}

func (queue Queue) NumberTasks() int {
    if queue.Tasks == nil {
        return queue.numberTasks
    }
    return queue.Tasks.Len()
}

func (queue Queue) CalculateTaskWeight(task *Task) int64 {
    return task.Created()
}

func (queue *Queue) PersistAddTask(task *Task) bool {
    if queue.Tasks == nil { return false }
    // update task inline
    task.AttrWeight = queue.CalculateTaskWeight(task)
    queue.Tasks.InsertNoReplace(task)
    return true
}

func (queue *Queue) AddTask(taskUid string) bool {
    updated := queue.store.AddTask(queue.Uid(), taskUid)
    queue.nextTaskUid = updated.NextTaskUid()
    queue.numberTasks = updated.NumberTasks()
    return true
}

func (queue *Queue) PersistDelTask(task *Task) bool {
    if queue.Tasks == nil { return false }
    queue.Tasks.Delete(task)
    return true
}

func (queue *Queue) DelTask(taskUid string) bool {
    updated := queue.store.DelTask(queue.Uid(), taskUid)
    queue.nextTaskUid = updated.NextTaskUid()
    queue.numberTasks = updated.NumberTasks()
    return true
}


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
