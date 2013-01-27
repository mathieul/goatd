package model

import (
    "reflect"
    "log"
    "fmt"
    "goatd/app/sm"
    "goatd/app/event"
)

const (
    attributePrefix = "Attr"
)

const (
    EventNone sm.Event = iota
    EventSignIn
    EventGoOnBreak
    EventSignOut
    EventMakeAvailable
    EventOfferTask
    EventAcceptTask
    EventRejectTask
    EventFinishTask
    EventStartOtherWork
    EventEnqueue
    EventDequeue
    EventOffer
    EventRequeue
    EventAssign
    EventComplete
)

const (
    StatusNone sm.Status = iota
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
func init() {
    sm.RegisterEvent(EventNone, "None")
    sm.RegisterEvent(EventSignIn, "Sign-in")
    sm.RegisterEvent(EventGoOnBreak, "Go on break")
    sm.RegisterEvent(EventSignOut, "Sign-out")
    sm.RegisterEvent(EventMakeAvailable, "Make available")
    sm.RegisterEvent(EventOfferTask, "Offer task")
    sm.RegisterEvent(EventAcceptTask, "Accept task")
    sm.RegisterEvent(EventRejectTask, "Reject task")
    sm.RegisterEvent(EventFinishTask, "Finish task")
    sm.RegisterEvent(EventStartOtherWork, "Start other work")
    sm.RegisterEvent(EventEnqueue, "Enqueue")
    sm.RegisterEvent(EventDequeue, "Dequeue")
    sm.RegisterEvent(EventOffer, "Offer")
    sm.RegisterEvent(EventRequeue, "Requeue")
    sm.RegisterEvent(EventAssign, "Assign")
    sm.RegisterEvent(EventComplete, "Complete")

    sm.RegisterStatus(StatusNone, "None")
    sm.RegisterStatus(StatusSignedOut, "Signed out")
    sm.RegisterStatus(StatusOnBreak, "On break")
    sm.RegisterStatus(StatusWaiting, "Waiting")
    sm.RegisterStatus(StatusOffered, "Offered")
    sm.RegisterStatus(StatusBusy, "Busy")
    sm.RegisterStatus(StatusWrappingUp, "Wrapping up")
    sm.RegisterStatus(StatusCompleted, "Completed")
    sm.RegisterStatus(StatusOtherWork, "Other work")
    sm.RegisterStatus(StatusCreated, "Created")
    sm.RegisterStatus(StatusQueued, "Queued")
    sm.RegisterStatus(StatusAssigned, "Assigned")
}


/*
 * Basic types and interfaces
 */
type A map[string]interface{}

type Model interface {
    Uid() string
    SetupComs(*event.BusManager, *Store)
    Copy() Model
    IsCopy() bool
    Status(...sm.Status) sm.Status
}


/*
 * Helpers
 */
func setFieldValue(model Model, name string, value interface{}) {
    destValue := reflect.ValueOf(model).Elem()
    if destValue.Type().Kind() != reflect.Struct {
        log.Fatal(fmt.Errorf("setAttributeValue(): model must be a pointer to a Struct, not %v", destValue.Type().Kind()))
    }
    field := destValue.FieldByName(name)
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

func setAttributeValue(model Model, name string, value interface{}) {
    setFieldValue(model, attributePrefix + name, value)
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

