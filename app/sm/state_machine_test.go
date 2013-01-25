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

type StateMachineSuite struct {}

var _ = Suite(&StateMachineSuite{})

func (s *StateMachineSuite) TestCreateStateMachineSetsInitialStatus(c *C) {
    stateMachine := sm.NewStateMachine(statusClosed, func (b sm.Builder) {})
    c.Assert(stateMachine.Status(), Equals, statusClosed)
}

func (s *StateMachineSuite) TestRegisterEvent(c *C) {
    c.Assert(eventOpen.String(), Equals, "Event{?Unknown?}")
    sm.RegisterEvent(eventOpen, "Open")
    c.Assert(eventOpen.String(), Equals, "Event{Open}")
}

func (s *StateMachineSuite) TestRegisterStatus(c *C) {
    c.Assert(statusLocked.String(), Equals, "Status{?Unknown?}")
    sm.RegisterStatus(statusLocked, "Locked")
    c.Assert(statusLocked.String(), Equals, "Status{Locked}")
}

func (s *StateMachineSuite) TestSingleTransitionEventsNoAction(c *C) {
    stateMachine := sm.NewStateMachine(statusClosed, func (b sm.Builder) {
        b.Event(eventOpen, statusClosed, statusOpened, sm.NoAction)
        b.Event(eventClose, statusOpened, statusClosed, sm.NoAction)
    })
    c.Assert(stateMachine.Trigger(eventOpen), Equals, true)
    c.Assert(stateMachine.Status(), Equals, statusOpened)
    c.Assert(stateMachine.Trigger(eventClose), Equals, true)
    c.Assert(stateMachine.Status(), Equals, statusClosed)
}

func (s *StateMachineSuite) TestUniTransitionEventsWithAction(c *C) {
    state := "not set"
    stateMachine := sm.NewStateMachine(statusClosed, func (b sm.Builder) {
        b.Event(eventOpen, statusClosed, statusOpened, func (args []interface{}) bool {
            state = "opened"
            return true
        })
        b.Event(eventClose, statusOpened, statusClosed, func (args []interface{}) bool {
            state = "closed"
            return true
        })
    })
    c.Assert(stateMachine.Trigger(eventOpen), Equals, true)
    c.Assert(state, Equals, "opened")
    c.Assert(stateMachine.Trigger(eventClose), Equals, true)
    c.Assert(state, Equals, "closed")
}

func (s *StateMachineSuite) TestActionPreventingStatusChange(c *C) {
    stateMachine := sm.NewStateMachine(statusClosed, func (b sm.Builder) {
        b.Event(eventOpen, statusClosed, statusOpened, func (args []interface{}) bool {
            return false
        })
    })
    c.Assert(stateMachine.Trigger(eventOpen), Equals, false)
}

func (s *StateMachineSuite) TestTriggeringWithArguments(c *C) {
    var hello string
    var fortyTwo int
    stateMachine := sm.NewStateMachine(statusClosed, func (b sm.Builder) {
        b.Event(eventOpen, statusClosed, statusOpened, func (args []interface{}) bool {
            hello = args[0].(string)
            fortyTwo = args[1].(int)
            return true
        })
    })
    c.Assert(stateMachine.Trigger(eventOpen, "hello", 42), Equals, true)
    c.Assert(hello, Equals, "hello")
    c.Assert(fortyTwo, Equals, 42)
}

func (s *StateMachineSuite) TestMultiTransitionEvents(c *C) {
    stateMachine := sm.NewStateMachine(statusClosed, func (b sm.Builder) {
        b.Event(eventOpen, statusClosed, statusOpened, sm.NoAction)
        b.Event(eventClose, statusOpened, statusClosed, sm.NoAction)
        b.Event(eventLock, func (b sm.Builder) {
            b.Transition(statusClosed, statusLocked, sm.NoAction)
            b.Transition(statusOpened, statusLockedOpen, sm.NoAction)
        })
    })
    c.Assert(stateMachine.Trigger(eventLock), Equals, true)
    c.Assert(stateMachine.Status(), Equals, statusLocked)
}

func (s *StateMachineSuite) TestSetTriggerValidator(c *C) {
    stateMachine := sm.NewStateMachine(statusClosed, func (b sm.Builder) {
        b.Event(eventOpen, statusClosed, statusOpened, sm.NoAction)
        b.Event(eventClose, statusOpened, statusClosed, sm.NoAction)
    })
    stateMachine.SetTriggerValidator(func (args ...interface{}) bool {
        return args[0].(bool)
    })
    c.Assert(stateMachine.Trigger(eventOpen, false), Equals, false)
    c.Assert(stateMachine.Status(), Equals, statusClosed)
    c.Assert(stateMachine.Trigger(eventClose, true), Equals, false)
}
