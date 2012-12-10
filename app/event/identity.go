package event

/*
 * Identity
 */

type Identity struct {
    kind string
    uid string
}

func NewIdentity(kind, uid string) *Identity {
    return &Identity{kind, uid}
}

func NoIdentity() Identity {
    return Identity{"", ""}
}

func (identity *Identity) Set(kind, uid string) {
    identity.kind = kind
    identity.uid = uid
}

func (identity Identity) Kind() string {
    return identity.kind
}

func (identity Identity) Uid() string {
    return identity.uid
}

func (identity Identity) AddToAttributes(attributes map[string]interface{}) map[string]interface{} {
    attributes[identity.kind + "Uid"] = identity.uid
    return attributes
}
