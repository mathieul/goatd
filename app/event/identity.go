package event

import (
    "os"
    "fmt"
    "log"
)


/*
 * Global
 */

const (
    randomDevice = "/dev/urandom"
)

/*
 * Helpers
 */

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


/*
 * Identity
 */

type Identity struct {
    kind string
    uid string
}

func NewIdentity(values ...string) (identity *Identity) {
    identity = new(Identity)
    switch len(values) {
    case 0:
        // nothing to do, no identity
    case 1:
        identity.kind = values[0]
        identity.uid = generateUid()
    default:
        identity.kind, identity.uid = values[0], values[1]
    }
    return identity
}

func NoIdentity(identity *Identity) bool {
    if identity.uid == "" {
        return true
    }
    return false
}

func (identity *Identity) Set(kind, uid string) {
    identity.kind = kind
    identity.uid = uid
}

func (identity Identity) Kind() string {
    return identity.kind
}

func (identity *Identity) SetKind(value string) {
    identity.kind = value
}

func (identity Identity) Uid() string {
    return identity.uid
}

func (identity *Identity) SetUid(value string) {
    identity.uid = value
}

func (identity *Identity) Copy(original *Identity) *Identity {
    return &Identity{original.kind, original.uid}
}

func (identity Identity) AddToAttributes(attributes map[string]interface{}) map[string]interface{} {
    attributes[identity.kind + "Uid"] = identity.uid
    return attributes
}
