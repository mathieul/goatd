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
    Kind
    Identity
    Data []string
}
type EventBus chan *Event


/*
 * Identity
 */

type Loner interface {
    Uid() string
}

type Identity struct {
    kind string
    uid string
    value Loner
}

func NewIdentity(kind, uid string, value Loner) *Identity {
    return &Identity{kind, uid, value}
}

func (identity *Identity) Set(kind, uid string, value Loner) {
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

func (identity Identity) Value() Loner {
    return identity.value
}

func (identity Identity) AddToAttributes(attributes map[string]interface{}) map[string]interface{} {
    attributes[identity.kind + "Uid"] = identity.uid
    return attributes
}


/*
 * busManager
 */
type busManager struct {
    incoming chan Event
    outgoings []chan Event
    done chan bool
}

func (busManager *busManager) Start() {
    busManager.incoming = make(chan Event, 0)
    busManager.done = make(chan bool, 0)
    busManager.outgoings = []chan Event{}
    go func() {
        for {
            select {
            case event := <- busManager.incoming:
                for _, outgoing := range busManager.outgoings {
                    outgoing <- event
                }
            case <- busManager.done:
                break
            }
        }
    }()
}

func (busManager *busManager) Stop() {
    busManager.done <- true
    busManager.incoming = nil
    busManager.outgoings = nil
    busManager.done = nil
}

func (busManager *busManager) PublishEvent(kind Kind, identity Identity, data []string) {
    event := Event{kind, identity, data}
    busManager.incoming <- event
}

func (busManager *busManager) SubscribeToAllEvents() (<-chan Event) {
    outgoing := make(chan Event, 0)
    busManager.outgoings = append(busManager.outgoings, outgoing)
    return outgoing
}

var manager busManager
func Manager() *busManager {
    return &manager
}
