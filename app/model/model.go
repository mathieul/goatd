package model

type Fields []string
type Attributes map[string]interface{}

type AttributeOwner struct {
	hasField map[string]bool
	attributes Attributes
}

func normalizeAttributes(hasField map[string]bool, attributes *Attributes) *Attributes {
	normalized := make(Attributes, len(*attributes))
	for key, value := range *attributes {
		if hasField[key] {
			normalized[key] = value
		}
	}
	return &normalized
}

func NewAttributeOwner(fields *Fields, attributes *Attributes) *AttributeOwner {
	hasField := make(map[string]bool, len(*fields))
	for _, name := range *fields {
		hasField[name] = true
	}
	attributes = normalizeAttributes(hasField, attributes)
	return &AttributeOwner{hasField, *attributes}
}

func (owner *AttributeOwner) Get(key string) (value interface{}, found bool) {
	value, found = owner.attributes[key]
	return value, found
}