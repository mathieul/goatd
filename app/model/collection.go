package model

import (
    "goatd/app/event"
)

/*
 * Collection
 */

type CollectionCreator func (A) Model
type Collection struct {
    creator CollectionCreator
    Items []Model
    owner *event.Identity
}

func NewCollection(creator CollectionCreator, owner *event.Identity) (collection *Collection) {
    collection = new(Collection)
    collection.creator = creator
    collection.owner = owner
    return collection
}

func (collection *Collection) New(attributes A) Model {
    if collection.owner != nil {
        attributes = collection.owner.AddToAttributes(attributes)
    }
    model := collection.creator(attributes)
    collection.Items = append(collection.Items, model)
    return model
}

func (collection *Collection) Slice() []Model {
    return collection.Items
}

func (collection Collection) Find(uid string) Model {
    found := collection.FindAll([]string{uid})
    if len(found) == 0 {
        return nil
    }
    return found[0]
}

func (collection Collection) FindAll(uids []string) (found []Model) {
    for _, candidate := range collection.Items {
        candidateUid := simpleMethodCall(candidate, "Uid").(string)
        for _, uid := range uids {
            if candidateUid == uid {
                found = append(found, candidate)
            }
        }
    }
    return found
}

func (collection Collection) Select(tester func(interface{}) bool) (result []Model) {
    result = make([]Model, 0)
    for _, item := range collection.Items {
        if tester(item) {
            result = append(result, item)
        }
    }
    return result
}
