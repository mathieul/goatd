package model

import (
    "reflect"
    "log"
    "fmt"
)

const (
    attributePrefix = "Attr"
)


/*
 * Basic types
 */
type A map[string]interface{}


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

