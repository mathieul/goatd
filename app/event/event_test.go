package event_test

import (
    . "launchpad.net/gocheck"
    "testing"
    "goatd/app/event"
)

func Test(t *testing.T) { TestingT(t) }

type EventSuite struct{
    manager event.BusManager
    identity *event.Identity
}

var _ = Suite(&EventSuite{})

func (s *EventSuite) SetUpTest(c *C) {
    s.manager.Start()
    s.identity = event.NewIdentity("Team", "qazwsx098")
}

func (s *EventSuite) TearDownTest(c *C) {
    s.manager.Stop()
}

func (s *EventSuite) TestPublishEvent(c *C) {
    done := make(chan bool)
    var one, two string
    go func () {
        // TODO: register and listen
        done <- true
    }()
    s.manager.PublishEvent(event.OfferTask, s.identity, []string{"Forty Two", "twelve"})
    <- done
    s.Assert()
}
