package model

import (
    "fmt"
    "goatd/app/event"
    "goatd/app/sm"
)

const (
    KindNone Kind = iota
    KindTeam
    KindTeammate
    KindTask
    KindQueue
    KindSkill
)

func (kind Kind) String() string {
    var value string
    switch kind {
    case KindNone:      value = "None"
    case KindTeam:      value = "Team"
    case KindTeammate:  value = "Teammate"
    case KindTask:      value = "Task"
    case KindQueue:     value = "Queue"
    case KindSkill:     value = "Skill"
    default:            value = fmt.Sprintf("Unknown(%d)", kind)
    }
    return fmt.Sprintf("<Kind{%s}>", value)
}

const (
    OpNone Operation = iota
    OpCreate
    OpUpdate
    OpDestroy
    OpDestroyAll
    OpSetStatus
    OpFind
    OpFindAll
    OpSelect
    OpEach
    OpCount
    OpAddTask
    OpDelTask
)

func (operation Operation) String() string {
    var value string
    switch operation {
    case OpNone:        value = "None"
    case OpCreate:      value = "Create"
    case OpUpdate:      value = "Update"
    case OpDestroy:     value = "Destroy"
    case OpDestroyAll:  value = "Destroy all"
    case OpSetStatus:   value = "Set status"
    case OpFind:        value = "Find"
    case OpFindAll:     value = "Find all"
    case OpSelect:      value = "Select"
    case OpEach:        value = "Each"
    case OpCount:       value = "Count"
    case OpAddTask:     value = "Add task"
    case OpDelTask:     value = "Del task"
    default:            value = fmt.Sprintf("Unknown(%d)", operation)
    }
    return fmt.Sprintf("<Operation{%s}>", value)
}

/*
 * Global and init
 */

var persisted *persistentStorage
func init() {
    persisted = newPersistentStorage()
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
 * Store API
 */

type Store struct {
    busManager  *event.BusManager
    Teams       *TeamStoreProxy
    Teammates   *TeammateStoreProxy
    Tasks       *TaskStoreProxy
    Queues      *QueueStoreProxy
    Skills      *SkillStoreProxy
}

func NewStore(busManager *event.BusManager) (store *Store) {
    store = new(Store)
    store.busManager = busManager
    store.Teams = &TeamStoreProxy{store}
    store.Teammates = &TeammateStoreProxy{store}
    store.Tasks = &TaskStoreProxy{store}
    store.Queues = &QueueStoreProxy{store}
    store.Skills = &SkillStoreProxy{store}
    return store
}

func (store *Store) Create(kind Kind, attributes A) Model {
    args := []interface{}{attributes}
    persisted.Request <- Request{kind, OpCreate, args}
    if value := <- persisted.Response; value != nil {
        model := value.(Model)
        model.SetupComs(store.busManager, store)
        return model
    }
    return nil
}

func (store *Store) Find(kind Kind, uid string) Model {
    args := []interface{}{uid}
    persisted.Request <- Request{kind, OpFind, args}
    if value := <- persisted.Response; value != nil {
        model := value.(Model)
        model.SetupComs(store.busManager, store)
        return model
    }
    return nil
}

func (store *Store) Update(kind Kind, uid, name string, value interface{}) bool {
    args := []interface{}{uid, name, value}
    persisted.Request <- Request{kind, OpUpdate, args}
    if value := <- persisted.Response; value != nil {
        return value.(bool)
    }
    return false
}

func (store *Store) Destroy(kind Kind, uid string) Model {
    args := []interface{}{uid}
    persisted.Request <- Request{kind, OpDestroy, args}
    if value := <- persisted.Response; value != nil {
        model := value.(Model)
        model.SetupComs(store.busManager, store)
        return model
    }
    return nil
}

func (store *Store) DestroyAll(kind Kind) {
    persisted.Request <- Request{kind, OpDestroyAll, []interface{}{}}
    <- persisted.Response
}

func (store *Store) SetStatus(kind Kind, uid string, oldStatus, newStatus sm.Status) bool {
    args := []interface{}{uid, oldStatus, newStatus}
    persisted.Request <- Request{kind, OpSetStatus, args}
    if value := <- persisted.Response; value != nil {
        return value.(bool)
    }
    return false
}

func (store *Store) FindAll(kind Kind, uids []string) []Model {
    args := []interface{}{uids}
    persisted.Request <- Request{kind, OpFindAll, args}
    values := <- persisted.Response
    models := make([]Model, 0, len(values.([]Model)))
    for _, value := range values.([]Model) {
        model := value.(Model)
        model.SetupComs(store.busManager, store)
        models = append(models, model)
    }
    return models
}

func (store *Store) Select(kind Kind, tester func(interface{}) bool) []Model {
    args := []interface{}{tester}
    persisted.Request <- Request{kind, OpSelect, args}
    values := <- persisted.Response
    models := make([]Model, 0, len(values.([]Model)))
    for _, value := range values.([]Model) {
        model := value.(Model)
        model.SetupComs(store.busManager, store)
        models = append(models, model)
    }
    return models
}

func (store *Store) Each(kind Kind, iterator func(interface{})) {
    args := []interface{}{iterator}
    persisted.Request <- Request{kind, OpEach, args}
    <- persisted.Response
}

func (store *Store) Count(kind Kind) int {
    persisted.Request <- Request{kind, OpCount, []interface{}{}}
    value := <- persisted.Response
    return value.(int)
}

func (store *Store) AddTask(uid, taskUid string) *Queue {
    args := []interface{}{uid, taskUid}
    persisted.Request <- Request{KindQueue, OpAddTask, args}
    if value := <- persisted.Response; value != nil {
        return value.(*Queue)
    }
    return nil
}

func (store *Store) DelTask(uid, taskUid string) *Queue {
    args := []interface{}{uid, taskUid}
    persisted.Request <- Request{KindQueue, OpDelTask, args}
    if value := <- persisted.Response; value != nil {
        return value.(*Queue)
    }
    return nil
}
