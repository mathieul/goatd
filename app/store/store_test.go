package store_test

import (
    . "launchpad.net/gocheck"
    "testing"
    "goatd/app/model"
    "goatd/app/store"
)

func Test(t *testing.T) { TestingT(t) }


type StoreSuite struct{
    store *store.Store
}

var _ = Suite(&StoreSuite{})

func (s *StoreSuite) SetUpTest(c *C) {
    s.store = store.NewStore()
}

func (s *StoreSuite) TestCreateTeam(c *C) {
    team := s.store.CreateTeam(model.A{"Name": "Aria & Tyrion"})
    c.Assert(team, Not(IsNil))
    c.Assert(team.Name(), Equals, "Aria & Tyrion")

    teamFound := s.store.FindTeam(team.Uid())
    c.Assert(teamFound.Uid(), Equals, team.Uid())
    c.Assert(teamFound.Name(), Equals, "Aria & Tyrion")
    c.Assert(teamFound, Not(Equals), team)
}
