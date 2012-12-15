package model

import (
    "fmt"
)

const (
    KindNone Kind = iota
    KindTeam
)

const (
    OpNone Operation = iota
    OpCreate
    OpFind
)

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
                response = collection.Find(request.args[0].(string)).Copy()
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

type Store struct {}

func NewStore() *Store {
    return &Store{}
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

func (store *Store) CreateTeam(attributes A) *Team {
    return store.Create(KindTeam, attributes).(*Team)
}

func (store *Store) FindTeam(uid string) *Team {
    if value := store.Find(KindTeam, uid); value != nil {
        return value.(*Team)
    }
    return nil
}
