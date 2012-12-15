package model

import (
    "fmt"
)

const (
    KindNone Kind = iota
    KindTeam
)

func (kind Kind) String() string {
    var value string
    switch kind {
    case KindNone:  value = "Opera"
    case KindTeam:  value = "Create"
    default:        value = fmt.Sprintf("Unknown(%d)", kind)
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
    case OpNone:    value = "Opera"
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

type request struct {
    Kind
    Operation
    args []interface{}
}


/*
 * Persistent store
 */

type persistentStore struct {
    Request chan request
    Response chan interface{}
    collections map[Kind]*Collection
}

func newPersistentStore() (store *persistentStore) {
    store = new(persistentStore)
    store.Request = make(chan request, 0)
    store.Response = make(chan interface{}, 0)
    store.collections = make(map[Kind]*Collection)
    store.collections[KindTeam] = NewCollection(func(attributes A) Model {
        team := NewTeam(attributes)
        return team
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

func (store *persistentStore) start() {
    go func() {
        for {
            var response interface{}
            request := <- store.Request
            collection := store.collections[request.Kind]

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
                panic(fmt.Errorf("Unknown operation %v\n", request.Operation))
            }

            store.Response <- response
        }
    }()
}

/*
 * Store API
 */

type Store struct {
    Teams *TeamStoreProxy
}

func NewStore() (store *Store) {
    store = new(Store)
    store.Teams = &TeamStoreProxy{store}
    return store
}

func (store *Store) Create(kind Kind, attributes A) Model {
    args := []interface{}{attributes}
    persisted.Request <- request{kind, OpCreate, args}
    value := <- persisted.Response
    return value.(Model)
}

func (store *Store) Find(kind Kind, uid string) Model {
    args := []interface{}{uid}
    persisted.Request <- request{kind, OpFind, args}
    if value := <- persisted.Response; value != nil {
        return value.(Model)
    }
    return nil
}

func (store *Store) FindAll(kind Kind, uids []string) []Model {
    args := []interface{}{uids}
    persisted.Request <- request{kind, OpFindAll, args}
    values := <- persisted.Response
    models := make([]Model, 0, len(values.([]Model)))
    for _, value := range values.([]Model) {
        models = append(models, value.(Model))
    }
    return models
}

func (store *Store) Select(kind Kind, tester func(interface{}) bool) []Model {
    args := []interface{}{tester}
    persisted.Request <- request{kind, OpSelect, args}
    values := <- persisted.Response
    models := make([]Model, 0, len(values.([]Model)))
    for _, value := range values.([]Model) {
        models = append(models, value.(Model))
    }
    return models
}
