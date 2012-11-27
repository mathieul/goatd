package models

import (
	"reflect"
	"log"
	"fmt"
	"os"
)

type Attrs map[string]interface{}

type Team struct {
    AttrUid string
    AttrName string
}

const (
	attributePrefix = "Attr"
	randomDevice = "/dev/urandom"
)

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
			field.SetInt(value.(int64))
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

func CreateTeam(attributes Attrs) (team *Team) {
    team = new(Team)
    for name, value := range attributes {
	    setAttributeValue(team, name, value)
    }
    if team.AttrUid == "" {
    	team.AttrUid = generateUid()
    }
    return team
}

func (team *Team) Uid() string {
    return team.AttrUid
}

func (team *Team) Name() string {
    return team.AttrName
}
