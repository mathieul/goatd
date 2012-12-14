package store

import (
    "fmt"
    "goatd/app/model"
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
    collections map[Kind]*model.Collection
}

func newPersistentStore() (store *persistentStore) {
    store = new(persistentStore)
    store.Request = make(chan request, 0)
    store.Response = make(chan interface{}, 0)
    store.collections = make(map[Kind]*model.Collection)
    store.collections[KindTeam] = model.NewCollection(func(attributes model.A) model.Model {
        team := model.NewTeam(attributes)
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
                response = collection.New(request.args[0].(model.A))
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

func (store *Store) Create(kind Kind, attributes model.A) model.Model {
    args := []interface{}{attributes}
    persisted.Request <- request{kind, OpCreate, args}
    value := <- persisted.Response
    return value.(model.Model)
}

func (store *Store) Find(kind Kind, uid string) model.Model {
    args := []interface{}{uid}
    persisted.Request <- request{kind, OpFind, args}
    if value := <- persisted.Response; value != nil {
        return value.(model.Model)
    }
    return nil
}

func (store *Store) CreateTeam(attributes model.A) *model.Team {
    return store.Create(KindTeam, attributes).(*model.Team)
}

func (store *Store) FindTeam(uid string) *model.Team {
    if value := store.Find(KindTeam, uid); value != nil {
        return value.(*model.Team)
    }
    return nil
}
