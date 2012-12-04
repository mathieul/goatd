package event

import (
    "goatd/app/identification"
)

const (
    OfferTask Kind = iota
    AssignTask
    CompleteTask
)

var allKinds []Kind

/*
 * initialization
 */
var manager *busManager
func init() {
    manager = new(busManager)
    allKinds = []Kind{OfferTask, AssignTask, CompleteTask}
}

func Manager() *busManager {
    return manager
}


/*
 * Basic types
 */
type Kind int
type Event struct {
    Kind
    identification.Identity
    Data []string
}
type EventBus chan *Event


/*
 * busManager
 */
type busManager struct {
    incoming chan Event
    outgoings map[Kind][]chan Event
    done chan bool
}

func (busManager *busManager) Init() {
    busManager.incoming = make(chan Event, 0)
    busManager.done = make(chan bool, 0)
    busManager.outgoings = make(map[Kind][]chan Event, len(allKinds))
    for _, kind := range allKinds {
        busManager.outgoings[kind] = make([]chan Event, 0, 5)
    }
}

func (busManager *busManager) Start() {
    busManager.Init()
    go func() {
        for {
            select {
            case event := <- busManager.incoming:
                for _, outgoing := range busManager.outgoings[event.Kind] {
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

func (busManager *busManager) PublishEvent(kind Kind, identity identification.Identity, data []string) {
    event := Event{kind, identity, data}
    busManager.incoming <- event
}

func (busManager *busManager) SubscribeTo(kinds []Kind) (<-chan Event) {
    outgoing := make(chan Event, 0)
    for _, kind := range kinds {
        busManager.outgoings[kind] = append(busManager.outgoings[kind], outgoing)
    }
    return outgoing
}

func (busManager *busManager) SubscribeToEvent(kind Kind) (<-chan Event) {
    return busManager.SubscribeTo([]Kind{kind})
}

func (busManager *busManager) SubscribeToAll() (<-chan Event) {
    return busManager.SubscribeTo(allKinds)
}
