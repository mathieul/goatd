package dispatch_test

import (
    . "launchpad.net/gocheck"
    "goatd/app/event"
    "goatd/app/model"
    "goatd/app/dispatch"
)

type DistributorSuite struct {
}

var _ = Suite(&DistributorSuite{})

func (s *DistributorSuite) SetUpTest(c *C) {
    s.busManager = event.NewBusManager()
    s.busManager.Start()
    s.store = model.NewStore(s.busManager)
    s.team = s.store.Teams.Create(model.A{"Name": "James Bond"})
}
