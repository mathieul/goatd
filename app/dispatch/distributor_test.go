package dispatch_test

import (
    . "launchpad.net/gocheck"
    "goatd/app/model"
    "goatd/app/dispatch"
)

type DistributorSuite struct {
    store *model.Store
}

var _ = Suite(&DistributorSuite{})

func (s *DistributorSuite) SetUpTest(c *C) {
    s.store = model.NewStore(nil)
}

func (s *DistributorSuite) TestCreateDistributor(c *C) {
    distributor := dispatch.NewDistributor(s.store)
    c.Assert(distributor.Store, Equals, s.store)
}
