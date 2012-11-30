package event_test

import (
    . "launchpad.net/gocheck"
    "testing"
    "goatd/app/event"
    "time"
)

func Test(t *testing.T) { TestingT(t) }

type EventSuite struct{
    manager event.BusManager
    identity *event.Identity
}

var _ = Suite(&EventSuite{})

func (s *EventSuite) SetUpTest(c *C) {
    // s.manager.Start()
    // s.identity = event.NewIdentity("Team", "qazwsx098")
}

func (s *EventSuite) TearDownTest(c *C) {
    // s.manager.Stop()
}

func (s *EventSuite) TestPublishEvent(c *C) {
    manager := new(event.BusManager)
    manager.Start()

    var received event.Event
    go func() {
        outgoing := manager.SubscribeToAllEvents()
        received = <- outgoing
    }()
    manager.PublishEvent(event.OfferTask)
    time.Sleep(200 * time.Millisecond)
    c.Assert(received.Kind, Equals, event.OfferTask)

    manager.Stop()
}
