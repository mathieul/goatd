package models

import (
    "goatd/app/event"
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

 type Collectioner interface {
    Create(Attrs) event.Loner
    Slice() []event.Loner
    Find(string) event.Loner
    FindAll([]string) []event.Loner
    // Select(Attrs) []event.Loner
}

type CollectionCreator func (Attrs, event.Loner) event.Loner
type Collection struct {
    creator CollectionCreator
    items []event.Loner
    owner event.Identity
}

func NewCollection(creator CollectionCreator, owner event.Identity) (collection Collection) {
    collection = *new(Collection)
    collection.creator = creator
    collection.owner = owner
    return collection
}

func (collection *Collection) Create(attributes Attrs) event.Loner {
    attributes = collection.owner.AddToAttributes(attributes)
    model := collection.creator(attributes, collection.owner.Value())
    collection.items = append(collection.items, model)
    return model
}

func (collection Collection) Slice() []event.Loner {
    slice := make([]event.Loner, len(collection.items))
    copy(slice, collection.items)
    return slice
}

func (collection Collection) Find(uid string) event.Loner {
    found := collection.FindAll([]string{uid})
    if len(found) == 0 {
        return nil
    }
    return found[0]
}

func (collection Collection) FindAll(uids []string) (found []event.Loner) {
    for _, candidate := range collection.items {
        candidateUid := candidate.Uid()
        for _, uid := range uids {
            if candidateUid == uid {
                found = append(found, candidate)
            }
        }
    }
    return found
}

// func (collection Collection) Select(query Attrs) (found []event.Loner) {
//     for _, candidate := range collection.items {
//         match := true
//         for name, value := range query {
//             if candidate[name] != value {
//                 match = false
//             }
//         }
//         if match {
//             found = append(found, candidate)
//         }
//     }
//     return found
// }


type Queues struct {
    Collection
}

func NewQueues(owner event.Identity) (queues *Queues) {
    queues = new(Queues)
    queues.Collection = NewCollection(func(attributes Attrs, lonerTeam event.Loner) event.Loner {
        queue := CreateQueue(attributes)
        queue.SetTeam(lonerTeam.(*Team))
        return queue
    }, owner)
    return queues
}

func (queues *Queues) Create(attributes Attrs) (queue *Queue) {
    queue = CreateQueue(queues.owner.AddToAttributes(attributes))
    queues.items = append(queues.items, queue)
    return queue
}
