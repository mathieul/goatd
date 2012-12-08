package models

import (
    "reflect"
    "log"
    "fmt"
    "os"
    "goatd/app/identification"
)

const (
    attributePrefix = "Attr"
    randomDevice = "/dev/urandom"
)

const (
    StatusNone Status = iota
    StatusSignedOut
    StatusOnBreak
    StatusWaiting
    StatusOffered
    StatusBusy
    StatusWrappingUp
    StatusCompleted
    StatusOtherWork

    StatusCreated
    StatusQueued
    StatusAssigned
)

const (
    LevelNone int = iota
    LevelLow
    LevelMedium
    LevelHigh
)

const (
    PriorityNone int = iota
    PriorityLow
    PriorityMedium
    PriorityHigh
)


/*
 * Basic types
 */
type Attrs map[string]interface{}
type Status int

func (status Status) String() string {
    return fmt.Sprintf("Status{%s}", statusToString[status])
}


/*
 * Init
 */
var statusFromString map[string]Status
var statusToString map[Status]string
var levelToString map[int]string
var boolToString map[bool]string

func init() {
    statusFromString = map[string]Status{
        "signed-out": StatusSignedOut,
        "on-break": StatusOnBreak,
        "waiting": StatusWaiting,
        "offered": StatusOffered,
        "busy": StatusBusy,
        "assigned": StatusAssigned,
        "wrapping-up": StatusWrappingUp,
        "completed": StatusCompleted,
        "other-work": StatusOtherWork,
        "created": StatusCreated,
        "queued": StatusQueued,
    }
    statusToString = map[Status]string{
        StatusNone: "None",
        StatusSignedOut: "SignedOut",
        StatusOnBreak: "OnBreak",
        StatusWaiting: "Waiting",
        StatusOffered: "Offered",
        StatusBusy: "Busy",
        StatusWrappingUp: "WrappingUp",
        StatusOtherWork: "OtherWork",
        StatusCreated: "Created",
        StatusQueued: "Queued",
        StatusAssigned: "Assigned",
        StatusCompleted: "Completed",
    }
    levelToString = map[int]string{
        LevelLow: "LevelLow",
        LevelMedium: "LevelMedium",
        LevelHigh: "LevelHigh",
    }
    boolToString = map[bool]string{
        true: "True",
        false: "False",
    }
}


/*
 * Helpers
 */
func setAttributeValue(destination interface{}, name string, value interface{}) {
    destValue := reflect.ValueOf(destination).Elem()
    if destValue.Type().Kind() != reflect.Struct {
        log.Fatal(fmt.Errorf("setAttributeValue(): destination must be a pointer to a Struct, not %v", destValue.Type().Kind()))
    }
    field := destValue.FieldByName(attributePrefix + name)
    if field.IsValid() {
        switch field.Type().Kind() {
        case reflect.String:
            field.SetString(value.(string))
        case reflect.Int:
            field.SetInt(int64(value.(int)))
        case reflect.Bool:
            field.SetBool(value.(bool))
        }
    }
}

func generateUid() string {
    data := make([]byte, 8)
    if randomizer, err := os.Open(randomDevice); err != nil {
        log.Fatal(fmt.Errorf("generateUid(): can't open random device %s (%q)", randomDevice, err))
    } else {
        defer randomizer.Close()
        randomizer.Read(data)
    }
    return fmt.Sprintf("%x-%x", data[0:4], data[4:])
}

func newModel(model interface{}, attributes *Attrs) interface{} {
    value := reflect.ValueOf(model)
    kind := value.Elem().Type().Kind()
    if kind != reflect.Struct {
        log.Fatal(fmt.Errorf("newModel(): model must be a Struct, and it is a %q", kind))
    }
    method := value.MethodByName("Init")
    if !method.IsValid() {
        log.Fatal(fmt.Errorf("newModel(): model must have a Storage field"))
    }
    method.Call([]reflect.Value{})
    for name, value := range *attributes {
        setAttributeValue(model, name, value)
    }
    return model
}

func simpleMethodCall(model interface{}, methodName string) interface{} {
    value := reflect.ValueOf(model)
    kind := value.Elem().Type().Kind()
    if kind != reflect.Struct {
        log.Fatal(fmt.Errorf("simpleMethodCall(): model must be a Struct, and it is a %q", kind))
    }
    method := value.MethodByName(methodName)
    if !method.IsValid() {
        log.Fatal(fmt.Errorf("simpleMethodCall(): model must have a %q method", methodName))
    }
    result := method.Call([]reflect.Value{})
    if len(result) != 1 {
        log.Fatal(fmt.Errorf("simpleMethodCall(): method %q must return one value", methodName))
    }
    return result[0].Interface()
}


/*
 * Storage
 */
type Storage struct {
    persisted bool
    uid string
}

func (storage *Storage) Init() {
    if storage.uid == "" {
        storage.uid = generateUid()
    }
}

func (storage Storage) Persisted() bool {
    return storage.persisted
}

func (storage Storage) Uid() string {
    return storage.uid
}

func (storage *Storage) Save() {
    storage.persisted = true
}


/*
 * Collection
 */

 type Collectioner interface {
    Create(Attrs) interface{}
    Find(string) interface{}
    FindAll([]string) []interface{}
    Select(func(interface{}) bool) []interface{}
}

type CollectionCreator func (Attrs, interface{}) interface{}
type Collection struct {
    creator CollectionCreator
    Items []interface{}
    owner identification.Identity
}

func NewCollection(creator CollectionCreator, owner identification.Identity) (collection Collection) {
    collection = *new(Collection)
    collection.creator = creator
    collection.owner = owner
    return collection
}

func (collection *Collection) Create(attributes Attrs) interface{} {
    attributes = collection.owner.AddToAttributes(attributes)
    model := collection.creator(attributes, collection.owner.Value())
    collection.Items = append(collection.Items, model)
    return model
}

func (collection *Collection) Slice() []interface{} {
    return collection.Items
}

func (collection Collection) Find(uid string) interface{} {
    found := collection.FindAll([]string{uid})
    if len(found) == 0 {
        return nil
    }
    return found[0]
}

func (collection Collection) FindAll(uids []string) (found []interface{}) {
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

func (collection Collection) Select(tester func(interface{}) bool) (result []interface{}) {
    result = make([]interface{}, 0)
    for _, item := range collection.Items {
        if tester(item) {
            result = append(result, item)
        }
    }
    return result
}
