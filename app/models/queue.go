package models

import (
    "goatd/app/identification"
)

/*
 * Queue
 */

type Queue struct {
    Storage
    team *Team
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

func (queue *Queue) SetTeam(team *Team) {
    queue.team = team
}

func (queue Queue) Team() (team *Team) {
    return queue.team
}


/*
 * Queues
 */

type Queues struct {
    Collection
}

func toQueueSlice(source []interface{}) []*Queue {
    queues := make([]*Queue, 0, len(source))
    for _, queue := range source {
        queues = append(queues, queue.(*Queue))
    }
    return queues
}

func NewQueues(owner identification.Identity) (queues *Queues) {
    queues = new(Queues)
    queues.Collection = NewCollection(func(attributes Attrs, parent interface{}) interface{} {
        queue := CreateQueue(attributes)
        queue.SetTeam(parent.(*Team))
        return queue
    }, owner)
    return queues
}

func (queues *Queues) Create(attributes Attrs) (queue *Queue) {
    return queues.Collection.Create(attributes).(*Queue)
}

func (queues Queues) Slice() []*Queue {
    return toQueueSlice(queues.Collection.Slice())
}

func (queues Queues) Find(uid string) *Queue {
    if found := queues.Collection.Find(uid); found != nil {
        return found.(*Queue)
    }
    return nil
}

func (queues Queues) FindAll(uids []string) []*Queue {
    return toQueueSlice(queues.Collection.FindAll(uids))
}

func (queues Queues) Select(tester func(interface{}) bool) (result []*Queue) {
    return toQueueSlice(queues.Collection.Select(tester))
}
