package event_test

import (
    . "launchpad.net/gocheck"
    "testing"
    "goatd/app/event"
    "time"
)

func Test(t *testing.T) { TestingT(t) }

type EventSuite struct {
    busManager *event.BusManager
    identity *event.Identity
}

var _ = Suite(&EventSuite{})

func (s *EventSuite) SetUpTest(c *C) {
    s.busManager = event.NewBusManager()
    s.busManager.Start()
    s.identity = event.NewIdentity("Team")
}

func (s *EventSuite) TearDownTest(c *C) {
    s.busManager.Stop()
}

func (s *EventSuite) TestPublishEvent(c *C) {
    var received event.Event
    go func() {
        incoming := s.busManager.SubscribeToAll()
        received = <- incoming
    }()
    s.busManager.PublishEvent(event.OfferTask, s.identity, []interface{}{"string"})
    time.Sleep(200 * time.Millisecond)
    c.Assert(received.Kind, Equals, event.OfferTask)
    c.Assert(received.Identity, Equals, s.identity)
    c.Assert(received.Data[0].(string), Equals, "string")
}

func (s *EventSuite) TestSubscribingToSomeEvents(c *C) {
    var e11, e12, e21, e22 event.Event
    go func() {
        incoming := s.busManager.SubscribeTo([]event.Kind{event.OfferTask, event.CompleteTask})
        e11 = <- incoming
        e12 = <- incoming
    }()
    go func() {
        incoming := s.busManager.SubscribeToEvent(event.CompleteTask)
        e21 = <- incoming
        e22 = <- incoming
    }()
    s.busManager.PublishEvent(event.CompleteTask, s.identity, []interface{}{"complete 1"})
    s.busManager.PublishEvent(event.OfferTask, s.identity, []interface{}{"offer"})

    time.Sleep(200 * time.Millisecond)
    c.Assert(e11.Data[0].(string), Equals, "complete 1")
    c.Assert(e12.Data[0].(string), Equals, "offer")
    c.Assert(e21.Data[0].(string), Equals, "complete 1")
    c.Assert(e22.Data, IsNil)
}
