package model

import (
    "fmt"
    "log"
)

const (
    KindNone Kind = iota
    KindTeam
    KindTeammate
    KindTask
)

func (kind Kind) String() string {
    var value string
    switch kind {
    case KindNone:      value = "None"
    case KindTeam:      value = "Team"
    case KindTeammate:  value = "Teammate"
    case KindTask:      value = "Task"
    default:            value = fmt.Sprintf("Unknown(%d)", kind)
    }
    return fmt.Sprintf("<Kind{%s}>", value)
}

const (
    OpNone Operation = iota
    OpCreate
    OpFind
    OpFindAll
    OpSelect
)

func (operation Operation) String() string {
    var value string
    switch operation {
    case OpNone:    value = "None"
    case OpCreate:  value = "Create"
    case OpFind:    value = "Find"
    case OpFindAll: value = "FindAll"
    case OpSelect:  value = "Select"
    default:        value = fmt.Sprintf("Unknown(%d)", operation)
    }
    return fmt.Sprintf("<Operation{%s}>", value)
}

/*
 * Global and init
 */

var persisted *persistentStore
func init() {
    persisted = newPersistentStore()
    persisted.start()
}


/*
 * Request & Response
 */
type Kind int
type Operation int

type Request struct {
    Kind
    Operation
    args []interface{}
}


/*
 * Persistent store
 */

type persistentStore struct {
    Request chan Request
    Response chan interface{}
    collections map[Kind]*Collection
}

func newPersistentStore() (store *persistentStore) {
    store = new(persistentStore)
    store.Request = make(chan Request, 0)
    store.Response = make(chan interface{}, 0)
    store.collections = make(map[Kind]*Collection)
    store.collections[KindTeam] = NewCollection(func(attributes A) Model {
        return NewTeam(attributes)
    }, nil)
    store.collections[KindTeammate] = NewCollection(func(attributes A) Model {
        return NewTeammate(attributes)
    }, nil)
    store.collections[KindTask] = NewCollection(func(attributes A) Model {
        return NewTask(attributes)
    }, nil)
    return store
}

func copyModels(models []Model) []Model {
    copied := make([]Model, 0, len(models))
    for _, model := range models {
        copied = append(copied, model.Copy())
    }
    return copied
}

func (store *persistentStore) processRequest(request Request, collection *Collection) (response interface{}) {
    switch request.Operation {
    case OpCreate:
        response = collection.New(request.args[0].(A))
    case OpFind:
        if model := collection.Find(request.args[0].(string)); model != nil {
            response = model.Copy()
        }
    case OpFindAll:
        models := collection.FindAll(request.args[0].([]string))
        response = copyModels(models)
    case OpSelect:
        tester := request.args[0].(func(interface{}) bool)
        models := collection.Select(tester)
        response = copyModels(models)
    default:
        log.Printf("Unknown operation %v\n", request.Operation)
    }
    return response
}

func (store *persistentStore) respondToRequests() {
    for {
        request := <- store.Request
        if collection := store.collections[request.Kind]; collection != nil {
            store.Response <- store.processRequest(request, collection)
        } else {
            log.Printf("No collection found for kind %v\n", request.Kind)
            store.Response <- nil
        }
    }
}

func (store *persistentStore) start() {
    go store.respondToRequests()
}

/*
 * Store API
 */

type Store struct {
    Teams     *TeamStoreProxy
    Teammates *TeammateStoreProxy
    Tasks     *TaskStoreProxy
}

func NewStore() (store *Store) {
    store = new(Store)
    store.Teams = &TeamStoreProxy{store}
    store.Teammates = &TeammateStoreProxy{store}
    store.Tasks = &TaskStoreProxy{store}
    return store
}

func (store *Store) Create(kind Kind, attributes A) Model {
    args := []interface{}{attributes}
    persisted.Request <- Request{kind, OpCreate, args}
    if value := <- persisted.Response; value != nil {
        return value.(Model)
    }
    return nil
}

func (store *Store) Find(kind Kind, uid string) Model {
    args := []interface{}{uid}
    persisted.Request <- Request{kind, OpFind, args}
    if value := <- persisted.Response; value != nil {
        return value.(Model)
    }
    return nil
}

func (store *Store) FindAll(kind Kind, uids []string) []Model {
    args := []interface{}{uids}
    persisted.Request <- Request{kind, OpFindAll, args}
    values := <- persisted.Response
    models := make([]Model, 0, len(values.([]Model)))
    for _, value := range values.([]Model) {
        models = append(models, value.(Model))
    }
    return models
}

func (store *Store) Select(kind Kind, tester func(interface{}) bool) []Model {
    args := []interface{}{tester}
    persisted.Request <- Request{kind, OpSelect, args}
    values := <- persisted.Response
    models := make([]Model, 0, len(values.([]Model)))
    for _, value := range values.([]Model) {
        models = append(models, value.(Model))
    }
    return models
}
