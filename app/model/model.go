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


/*
 * Basic types
 */
type A map[string]interface{}


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

func newModel(model interface{}, attributes *A) interface{} {
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
 * Collection
 */

type CollectionCreator func (A) interface{}
type Collection struct {
    creator CollectionCreator
    Items []interface{}
    owner *event.Identity
}

func NewCollection(creator CollectionCreator, owner *event.Identity) (collection Collection) {
    collection = *new(Collection)
    collection.creator = creator
    collection.owner = owner
    return collection
}

func (collection *Collection) New(attributes A) interface{} {
    if collection.owner != nil {
        attributes = collection.owner.AddToAttributes(attributes)
    }
    model := collection.creator(attributes)
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
