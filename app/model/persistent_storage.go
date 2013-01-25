package model

import (
    "log"
)

/*
 * Persistent storage
 */

type persistentStorage struct {
    Request chan Request
    Response chan interface{}
    collections map[Kind]*Collection
}

func newPersistentStorage() (storage *persistentStorage) {
    storage = new(persistentStorage)
    storage.Request = make(chan Request, 0)
    storage.Response = make(chan interface{}, 0)
    storage.collections = make(map[Kind]*Collection)
    storage.collections[KindTeam] = NewCollection(func(attributes A) Model {
        return NewTeam(attributes)
    }, nil)
    storage.collections[KindTeammate] = NewCollection(func(attributes A) Model {
        return NewTeammate(attributes)
    }, nil)
    storage.collections[KindTask] = NewCollection(func(attributes A) Model {
        return NewTask(attributes)
    }, nil)
    return storage
}

func copyModels(models []Model) []Model {
    copied := make([]Model, 0, len(models))
    for _, model := range models {
        copied = append(copied, model.Copy())
    }
    return copied
}

func (storage *persistentStorage) processRequest(request Request, collection *Collection) (response interface{}) {
    switch request.Operation {
    case OpCreate:
        model := collection.New(request.args[0].(A))
        response = model.Copy()
    case OpFind:
        if model := collection.Find(request.args[0].(string)); model != nil {
            response = model.Copy()
        }
    case OpFindAll:
        models := collection.FindAll(request.args[0].([]string))
        response = copyModels(models)
    case OpSelect:
        tester := request.args[0].(func (interface{}) bool)
        models := collection.Select(tester)
        response = copyModels(models)
    default:
        log.Printf("Unknown operation %v\n", request.Operation)
    }
    return response
}

func (storage *persistentStorage) respondToRequests() {
    for {
        request := <- storage.Request
        if collection := storage.collections[request.Kind]; collection != nil {
            storage.Response <- storage.processRequest(request, collection)
        } else {
            log.Printf("No collection found for kind %v\n", request.Kind)
            storage.Response <- nil
        }
    }
}

func (storage *persistentStorage) start() {
    go storage.respondToRequests()
}
