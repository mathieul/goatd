package state_machine_test

import (
    . "launchpad.net/gocheck"
    "goatd/app/state_machine"
)

type StateMachineSuite struct{
    sm *state_machine.StateMachine
}

var _ = Suite(&StateMachineSuite{})

func (s *StateMachineSuite) TestCreateStateMachine(c *C) {
    // list states
    // for each event, declare transitions (from, to, action)
}
