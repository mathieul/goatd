package event

import (
    "goatd/app/identification"
)

const (
    KindNone Kind = iota
    KindOfferTask
    KindAssignTask
    KindCompleteTask
    KindTeammateAvailable
)

var allKinds []Kind

/*
 * initialization
 */
var manager *busManager
func init() {
    manager = new(busManager)
    allKinds = []Kind{KindOfferTask, KindAssignTask, KindCompleteTask}
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
    running bool
}

func (busManager busManager) Running() bool {
    return busManager.running
}

func (busManager *busManager) setupChannels() {
    busManager.incoming = make(chan Event, 0)
    busManager.done = make(chan bool, 0)
    busManager.outgoings = make(map[Kind][]chan Event, len(allKinds))
    for _, kind := range allKinds {
        busManager.outgoings[kind] = make([]chan Event, 0, 5)
    }
}

func (busManager *busManager) cleanupChannels() {
    close(busManager.done)
    close(busManager.incoming)
    for _, channels := range busManager.outgoings {
        for _, channel := range channels {
            close(channel)
        }
    }
}

func (busManager *busManager) Start() {
    busManager.setupChannels()
    go func() {
        busManager.running = true
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
        busManager.running = false
        busManager.cleanupChannels()
    }()
}

func (busManager *busManager) Stop() {
    busManager.done <- true
}

func (busManager *busManager) PublishEvent(kind Kind,
        identity identification.Identity, data []string) bool {
    if busManager.running {
        event := Event{kind, identity, data}
        busManager.incoming <- event
        return true
    }
    return false
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
