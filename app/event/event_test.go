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
    event.Manager().PublishEvent(event.KindOfferTask, s.identity, []string{"string"})
    time.Sleep(200 * time.Millisecond)
    c.Assert(received.Kind, Equals, event.KindOfferTask)
    c.Assert(received.Identity, Equals, s.identity)
    c.Assert(received.Data[0], Equals, "string")
}

func (s *EventSuite) TestSubscribingToSomeEvents(c *C) {
    var e11, e12, e21, e22 event.Event
    go func() {
        incoming := event.Manager().SubscribeTo([]event.Kind{event.KindOfferTask, event.KindCompleteTask})
        e11 = <- incoming
        e12 = <- incoming
    }()
    go func() {
        incoming := event.Manager().SubscribeToEvent(event.KindCompleteTask)
        e21 = <- incoming
        e22 = <- incoming
    }()
    event.Manager().PublishEvent(event.KindCompleteTask, s.identity, []string{"complete 1"})
    event.Manager().PublishEvent(event.KindOfferTask, s.identity, []string{"offer"})

    time.Sleep(200 * time.Millisecond)
    c.Assert(e11.Data[0], Equals, "complete 1")
    c.Assert(e12.Data[0], Equals, "offer")
    c.Assert(e21.Data[0], Equals, "complete 1")
    c.Assert(e22.Data, IsNil)
}
