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
    // identity *Identity
    // data []string
}
type EventBus chan *Event


/*
 * BusManager
 */
type BusManager struct {
    incoming chan Event
    outgoings []chan Event
    done chan bool
}

func (busManager *BusManager) Start() {
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

func (busManager *BusManager) Stop() {
    busManager.done <- true
    busManager.incoming = nil
    busManager.outgoings = nil
    busManager.done = nil
}

func (busManager *BusManager) PublishEvent(kind Kind) {
    event := Event{kind}
    busManager.incoming <- event
}

func (busManager *BusManager) SubscribeToAllEvents() (<-chan Event) {
    outgoing := make(chan Event, 0)
    busManager.outgoings = append(busManager.outgoings, outgoing)
    return outgoing
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
