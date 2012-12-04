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

func (queue Queue) Uid() string {
    return queue.Storage.Uid()
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

func NewQueues(owner identification.Identity) (queues *Queues) {
    queues = new(Queues)
    queues.Collection = NewCollection(func(attributes Attrs, lonerTeam interface{}) interface{} {
        queue := CreateQueue(attributes)
        queue.SetTeam(lonerTeam.(*Team))
        return queue
    }, owner)
    return queues
}

func (queues *Queues) Create(attributes Attrs) (queue *Queue) {
    return queues.Collection.Create(attributes).(*Queue)
}

func (queues Queues) Slice() []*Queue {
    result := make([]*Queue, 0, len(queues.Items))
    for _, pointer := range queues.Items {
        result = append(result, pointer.(*Queue))
    }
    return result
}

func (queues Queues) Find(uid string) *Queue {
    return queues.Collection.Find(uid).(*Queue)
}

func (queues Queues) FindAll(uids []string) (result []*Queue) {
    for _, queue := range queues.Collection.FindAll(uids) {
        result = append(result, queue.(*Queue))
    }
    return result
}
