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
var eventOpen, eventClose sm.Event
var statusOpened, statusClosed sm.Status

func init() {
    eventOpen, eventClose = 1, 2
    statusOpened, statusClosed = 1, 2
}

type StateMachineSuite struct {
    sm *sm.StateMachine
}

var _ = Suite(&StateMachineSuite{})

func (s *StateMachineSuite) TestCreateStateMachine(c *C) {
    sm := sm.NewStateMachine(statusClosed, func (builder sm.Builder) {
        builder.Event(eventOpen, statusClosed, statusOpened, sm.NoAction)
        builder.Event(eventClose, statusOpened, statusClosed, func (args []interface{}) bool {
            return true
        })
    })
    c.Assert(sm.Status(), Equals, statusClosed)
    c.Assert(sm.Trigger(eventOpen), Equals, true)
    c.Assert(sm.Status(), Equals, statusOpened)
}
