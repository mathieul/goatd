package event_test

import (
    . "launchpad.net/gocheck"
    "testing"
    "goatd/app/event"
    "time"
)

func Test(t *testing.T) { TestingT(t) }

type EventSuite struct{
    identity event.Identity
}

var _ = Suite(&EventSuite{})

func (s *EventSuite) SetUpTest(c *C) {
    event.Manager().Start()
    s.identity.Set("Team", "qazwsx098")
}

func (s *EventSuite) TearDownTest(c *C) {
    event.Manager().Stop()
}

func (s *EventSuite) TestPublishEvent(c *C) {
    var received event.Event
    go func() {
        outgoing := event.Manager().SubscribeToAllEvents()
        received = <- outgoing
    }()
    event.Manager().PublishEvent(event.OfferTask, s.identity, []string{"string"})
    time.Sleep(200 * time.Millisecond)
    c.Assert(received.Kind, Equals, event.OfferTask)
    c.Assert(received.Identity, Equals, s.identity)
    c.Assert(received.Data[0], Equals, "string")
}
