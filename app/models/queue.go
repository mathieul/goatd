package models

import (
    "goatd/app/event"
)

/*
 * Queue
 */

type Queue struct {
    Storage
    AttrName string
    AttrTeamUid string
}

func NewQueue(attributes Attrs) *Queue {
    return newModel(&Queue{}, &attributes).(*Queue)
}

func CreateQueue(attributes Attrs) (queue *Queue) {
    queue = NewQueue(attributes)
    queue.Save()
    return queue
}

func (queue Queue) Name() string {
    return queue.AttrName
}

func (queue Queue) TeamUid() string {
    return queue.AttrTeamUid
}


/*
 * Queues
 */

type Queues struct {
    owner event.Identity
    items []*Queue
}

func NewQueues(owner event.Identity) (queues *Queues) {
    queues = new(Queues)
    queues.owner = owner
    return queues
}

func (queues *Queues) Create(attributes Attrs) (queue *Queue) {
    queue = CreateQueue(queues.owner.AddToAttributes(attributes))
    queues.items = append(queues.items, queue)
    return queue
}
