package model

import (
    "reflect"
    "log"
    "fmt"
    "goatd/app/event"
)

const (
    attributePrefix = "Attr"
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


/*
 * Init
 */
var statusFromString map[string]Status
var statusToString map[Status]string

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
        StatusNone: "none",
        StatusSignedOut: "signed-out",
        StatusOnBreak: "on-break",
        StatusWaiting: "waiting",
        StatusOffered: "offered",
        StatusBusy: "busy",
        StatusWrappingUp: "wrapping-up",
        StatusOtherWork: "other-work",
        StatusCreated: "created",
        StatusQueued: "queued",
        StatusAssigned: "assigned",
        StatusCompleted: "completed",
    }
}


/*
 * Basic types and interfaces
 */
type A map[string]interface{}
type Status int

func (status Status) String() string {
    return fmt.Sprintf("Status{%s}", statusToString[status])
}

type Model interface {
    Uid() string
    MakeActive(*event.BusManager)
    Copy() Model
}


/*
 * Helpers
 */
func setAttributeValue(destination Model, name string, value interface{}) {
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

func newModel(model Model, attributes *A) interface{} {
    value := reflect.ValueOf(model)
    kind := value.Elem().Type().Kind()
    if kind != reflect.Struct {
        log.Fatal(fmt.Errorf("newModel(): model must be a Struct, and it is a %q", kind))
    }
    for name, value := range *attributes {
        setAttributeValue(model, name, value)
    }
    return model
}

func simpleMethodCall(model Model, methodName string) interface{} {
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

