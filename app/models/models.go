package models

import (
    "reflect"
    "log"
    "fmt"
    "os"
)

const (
    attributePrefix = "Attr"
    randomDevice = "/dev/urandom"
)

const (
    StatusSignedOut Status = iota
    StatusOnBreak
    StatusOffered
    StatusBusy
    StatusAssigned
    StatusWrappingUp
    StatusCompleted
    StatusOtherWork
)

const (
    LevelLow int = iota
    LevelMedium
    LevelHigh
)

/*
 * Basic types
 */
type Attrs map[string]interface{}
type Status int

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

func simpleMethodCall(model interface{}, methodName string) (result interface{}) {
    value := reflect.ValueOf(model)
    kind := value.Elem().Type().Kind()
    if kind != reflect.Struct {
        log.Fatal(fmt.Errorf("simpleMethodCall(): model must be a Struct, and it is a %q", kind))
    }
    method := value.MethodByName(methodName)
    if !method.IsValid() {
        log.Fatal(fmt.Errorf("simpleMethodCall(): model must have a %q method", methodName))
    }
    return method.Call([]reflect.Value{})
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
