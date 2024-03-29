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
    if len(uids) == 0 { return collection.Items }
    for _, candidate := range collection.Items {
        candidateUid := candidate.Uid()
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

func (collection Collection) Each(iterator func(interface{})) {
    for _, item := range collection.Items { iterator(item) }
}

func (collection Collection) Count() int {
    return len(collection.Items)
}

func (collection *Collection) Destroy(uid string) Model {
    for index, item := range collection.Items {
        if item.Uid() == uid {
            collection.Items = append(collection.Items[:index], collection.Items[index + 1:]...)
            return item
        }
    }
    return nil
}

func (collection *Collection) DestroyAll() {
    collection.Items = make([]Model, 0)
}
