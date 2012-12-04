package event_test

import (
    . "launchpad.net/gocheck"
    "testing"
    "goatd/app/identification"
    "goatd/app/event"
    "time"
)

func Test(t *testing.T) { TestingT(t) }

type EventSuite struct{
    identity identification.Identity
}

var _ = Suite(&EventSuite{})

func (s *EventSuite) SetUpTest(c *C) {
    event.Manager().Start()
    s.identity.Set("Team", "qazwsx098", nil)
}

func (s *EventSuite) TearDownTest(c *C) {
    event.Manager().Stop()
}

func (s *EventSuite) TestPublishEvent(c *C) {
    var received event.Event
    go func() {
        incoming := event.Manager().SubscribeToAll()
        received = <- incoming
    }()
    event.Manager().PublishEvent(event.OfferTask, s.identity, []string{"string"})
    time.Sleep(200 * time.Millisecond)
    c.Assert(received.Kind, Equals, event.OfferTask)
    c.Assert(received.Identity, Equals, s.identity)
    c.Assert(received.Data[0], Equals, "string")
}

func (s *EventSuite) TestSubscribingToSomeEvents(c *C) {
    received1 := make([]event.Event, 0, 2)
    received2 := make([]event.Event, 0, 2)
    go func() {
        incoming := event.Manager().SubscribeTo([]event.Event{event.OfferTask, event.CompleteTask})
        received1[0] = <- incoming
        received1[1] = <- incoming
    }()
    go func() {
        incoming := event.Manager().SubscribeToEvent(event.CompleteTask)
        received2[0] = <- incoming
        received2[1] = <- incoming
    }()
    event.Manager().PublishEvent(event.CompleteTask, s.identity, []string{"complete 1"})
    event.Manager().PublishEvent(event.OfferTask, s.identity, []string{"offer"})
    event.Manager().PublishEvent(event.CompleteTask, s.identity, []string{"complete 2"})

    time.Sleep(200 * time.Millisecond)
    c.Assert(received1[0].Data[0], Equals, "complete 1")
    c.Assert(received1[1].Data[0], Equals, "offer")
    c.Assert(received1[0].Data[0], Equals, "complete 1")
    c.Assert(received1[1].Data[0], Equals, "complete 2")
}
