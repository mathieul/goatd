package event

const (
    OfferTask Kind = iota
    AssignTask
    CompleteTask
)


/*
 * Basic types
 */
type Kind int
type Event struct {
    kind Kind
    identity *Identity
    data []string
}
type Bus chan Event


/*
 * BusManager
 */
type BusManager struct {
    bus Bus
}

func (manager BusManager) Start() {
    manager.bus = make(Bus, 5)
}

func (manager BusManager) Stop() {}

func (manager BusManager) PublishEvent(kind Kind, identity *Identity, data []string) {
    manager.bus <- Event{kind, identity, data}
}

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
