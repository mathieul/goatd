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

func (team *Queue) Name() string {
    return team.AttrName
}

func (team *Queue) TeamUid() string {
    return team.AttrTeamUid
}

/*
 * Queues
 */

type Queues struct {
	owner *event.Identity
	items []*Queue
}

func NewQueues(kind, uid string) (queues *Queues) {
	queues = new(Queues)
	queues.owner = event.NewIdentity(kind, uid)
	return queues
}

func (queues *Queues) Create(attributes Attrs) (queue *Queue) {
	queue = CreateQueue(queues.owner.AddToAttributes(attributes))
	queues.items = append(queues.items, queue)
	return queue
}
