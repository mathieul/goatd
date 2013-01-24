package sm_test

import (
    . "launchpad.net/gocheck"
    "testing"
    "goatd/app/sm"
)

func Test(t *testing.T) { TestingT(t) }

/*
 * Setup
 */

const (
    eventOpen sm.Event = iota + 1
    eventClose
    eventLock
)

const (
    statusOpened sm.Status = iota + 1
    statusClosed
    statusLocked
    statusLockedOpen
)

type StateMachineSuite struct {
    sm *sm.StateMachine
}

var _ = Suite(&StateMachineSuite{})

func (s *StateMachineSuite) TestCreateStateMachineSetsInitialStatus(c *C) {
    sm := sm.NewStateMachine(statusClosed, func (b sm.Builder) {})
    c.Assert(sm.Status(), Equals, statusClosed)
}

func (s *StateMachineSuite) TestUniTransitionEventsNoAction(c *C) {
    sm := sm.NewStateMachine(statusClosed, func (b sm.Builder) {
        b.Event(eventOpen, statusClosed, statusOpened, sm.NoAction)
        b.Event(eventClose, statusOpened, statusClosed, sm.NoAction)
    })
    c.Assert(sm.Trigger(eventOpen), Equals, true)
    c.Assert(sm.Status(), Equals, statusOpened)
}

func (s *StateMachineSuite) TestUniTransitionEventsWithAction(c *C) {
    state := "not set"
    sm := sm.NewStateMachine(statusClosed, func (b sm.Builder) {
        b.Event(eventOpen, statusClosed, statusOpened, func (args []interface{}) bool {
            state = "opened"
            return true
        })
        b.Event(eventClose, statusOpened, statusClosed, func (args []interface{}) bool {
            state = "closed"
            return true
        })
    })
    c.Assert(sm.Trigger(eventOpen), Equals, true)
    c.Assert(state, Equals, "opened")
    c.Assert(sm.Trigger(eventClose), Equals, true)
    c.Assert(state, Equals, "closed")
}

func (s *StateMachineSuite) TestMultiTransitionEvents(c *C) {
    sm := sm.NewStateMachine(statusDraft, func (b sm.Builder) {
        b.Event(eventOpen, statusClosed, statusOpened, sm.NoAction)
        b.Event(eventClose, statusOpened, statusClosed, sm.NoAction)
        b.Event(eventLock, func (t sm.Transitioner) {
            t.Transition(statusClosed, statusLocked, sm.NoAction)
            t.Transition(statusOpened, statusLockedOpen, sm.NoAction)
        })
    })
    c.Assert(sm.Trigger(eventLock), Equals, true)
    c.Assert(sm.Status(), Equals, statusLockedOpened)
}
