package model

import (
	"log"
	"fmt"
)

var supportedTypes map[string]bool

func init() {
	supportedTypes = map[string]bool{
		"string": true,
	}
}

type Fields map[string]string
type Attributes map[string]interface{}

type AttributeOwner struct {
	typeFor Fields
	attributes Attributes
}

func normalizeAttributes(fields Fields, attributes Attributes) *Attributes {
	normalized := make(Attributes, len(attributes))
	for key, value := range attributes {
		if _, found := fields[key]; found {
			normalized[key] = value
		}
	}
	return &normalized
}

func checkFields(fields Fields) {
	for _, fieldType := range fields {
		if !supportedTypes[fieldType] {
			log.Fatal(fmt.Errorf("checkFields(): field type %q is not supported", fieldType))
		}
	}
}

func NewAttributeOwner(fields Fields, attributes Attributes) *AttributeOwner {
	normalized := normalizeAttributes(fields, attributes)
	checkFields(fields)
	return &AttributeOwner{fields, *normalized}
}

func (owner *AttributeOwner) Get(key string) (value interface{}, found bool) {
	if _, typeFound := owner.typeFor[key]; typeFound {
		value, found = owner.attributes[key]

	} else {
		found = false
	}
	return value, found
}
