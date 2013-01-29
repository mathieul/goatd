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
    queuedTaskUids []string
    AttrName string
    AttrTeamUid string
}

func NewQueue(attributes A) (queue *Queue) {
    queue = newModel(&Queue{}, &attributes).(*Queue)
    queue.Identity = event.NewIdentity("Queue")
    queue.queuedTaskUids = make([]string, 0, initialQueueLength)
    return queue
}

func (queue *Queue) Copy() Model {
    identity := queue.Identity.Copy()
    var queuedTaskUids []string
    copy(queuedTaskUids, queue.queuedTaskUids)
    return &Queue{identity, nil, nil, queuedTaskUids,
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

func (queue Queue) QueuedTaskUids() []string {
    return queue.queuedTaskUids
}

func (queue *Queue) AddTask(taskUid string) bool {
    queue.queuedTaskUids = append(queue.queuedTaskUids, taskUid)
    // TODO: push to persitent storage
    return true
}

func (queue *Queue) RemoveTask(taskUid string) bool {
    index := -1
    for i, uid := range queue.queuedTaskUids {
        if uid == taskUid {
            index = i
            break
        }
    }
    queue.queuedTaskUids = append(queue.queuedTaskUids[:index], queue.queuedTaskUids[index + 1:]...)
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
