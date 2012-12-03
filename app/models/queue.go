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

  type Modeler interface {
    Uid() string
 }

 type Collectioner interface {
    Create(Attrs) *Modeler
    Slice() []*Modeler
    Find(string) *Modeler
    FindAll([]string) []*Modeler
    Select(Attrs) []*Modeler
}

type CollectionCreator func (Attrs, Modeler) interface{}
type Collection struct {
    creator CollectionCreator
    items []Modeler
    owner event.Identity
}

func NewCollection(creator CollectionCreator, owner event.Identity) (collection *Collection) {
    collection = new(Collection)
    collection.creator = creator
    collection.owner = owner
    return collection
}

func (collection *Collection) Create(attributes Attrs) Modeler {
    attributes = collection.owner.AddToAttributes(attributes)
    model := collection.creator(attributes, collection.owner.Value())
    collection.items = append(collection.items, model)
    return model
}

func (collection Collection) Slice() []Modeler {
    slice := make([]Modeler, len(collection.items))
    copy(slice, collection.items)
    return slice
}

func (collection Collection) FindAll(uids []string) (found []Modeler) {
    for _, candidate := range collection.items.(Modeler) {
        candidateUid := candidate.Uid()
        for _, uid := range uids {
            if candidateUid == uid {
                found = append(found, candidate)
            }
        }
    }
    return found
}

func (collection Collection) Find(uid string) Modeler {
    found := collection.FindAll([]string{uid})
    if len(found) == 0 {
        return nil
    }
    return found[0]
}


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
