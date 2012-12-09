package event

/*
 * Identity
 */

type Identity struct {
    kind string
    uid string
    value interface{}
}

func NewIdentity(kind, uid string, value interface{}) *Identity {
    return &Identity{kind, uid, value}
}

func NoIdentity() Identity {
    return Identity{"", "", nil}
}

func (identity *Identity) Set(kind, uid string, value interface{}) {
    identity.kind = kind
    identity.uid = uid
    identity.value = value
}

func (identity Identity) Kind() string {
    return identity.kind
}

func (identity Identity) Uid() string {
    return identity.uid
}

func (identity Identity) Value() interface{} {
    return identity.value
}

func (identity Identity) AddToAttributes(attributes map[string]interface{}) map[string]interface{} {
    attributes[identity.kind + "Uid"] = identity.uid
    return attributes
}
