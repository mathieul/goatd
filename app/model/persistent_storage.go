package model

import (
    "log"
    "goatd/app/sm"
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
    storage.collections[KindQueue] = NewCollection(func(attributes A) Model {
        return NewQueue(attributes)
    }, nil)
    storage.collections[KindSkill] = NewCollection(func(attributes A) Model {
        return NewSkill(attributes)
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
        attributes := request.args[0].(A)
        model := collection.New(attributes)
        response = model.Copy()
    case OpUpdate:
        uid := request.args[0].(string)
        name, value := request.args[1].(string), request.args[2]
        if model := collection.Find(uid); model != nil {
            setAttributeValue(model, name, value)
            response = true
        } else {
            response = false
        }
    case OpSetStatus:
        uid := request.args[0].(string)
        oldStatus, newStatus := request.args[1].(sm.Status), request.args[2].(sm.Status)
        response = false
        if model := collection.Find(uid); model != nil {
            currentStatus := model.Status()
            if currentStatus == oldStatus {
                model.Status(newStatus)
                response = true
            }
        }
    case OpFind:
        uid := request.args[0].(string)
        if model := collection.Find(uid); model != nil {
            response = model.Copy()
        }
    case OpFindAll:
        uids := request.args[0].([]string)
        models := collection.FindAll(uids)
        response = copyModels(models)
    case OpSelect:
        selector := request.args[0].(func (interface{}) bool)
        models := collection.Select(selector)
        response = copyModels(models)
    case OpAddTask, OpDelTask:
        uid, taskUid := request.args[0].(string), request.args[1].(string)
        model := collection.Find(uid)
        taskModel := storage.collections[KindTask].Find(taskUid)
        if model != nil && taskModel != nil {
            queue := model.(*Queue)
            if request.Operation == OpAddTask {
                queue.PersistAddTask(taskModel.(*Task))
            } else {
                queue.PersistDelTask(taskModel.(*Task))
            }
            response = queue.Copy()
        }
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
