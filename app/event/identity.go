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

func NewIdentity(kind, uid string) (identity *Identity) {
    identity = &Identity{kind, uid}
    if uid == "" {
        identity.uid = generateUid()
    }
    return identity
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
