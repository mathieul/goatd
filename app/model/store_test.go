package model_test

import (
    . "launchpad.net/gocheck"
    "goatd/app/model"
)

type StoreSuite struct{
    store *model.Store
}

var _ = Suite(&StoreSuite{})

func (s *StoreSuite) SetUpTest(c *C) {
    s.store = model.NewStore(nil)
}

func (s *StoreSuite) TestCreateTeam(c *C) {
    team := s.store.Teams.Create(model.A{"Name": "Aria & Tyrion"})
    c.Assert(team, Not(IsNil))
    c.Assert(team.Name(), Equals, "Aria & Tyrion")

    teamFound := s.store.Teams.Find(team.Uid())
    c.Assert(teamFound.Uid(), Equals, team.Uid())
    c.Assert(teamFound.Name(), Equals, "Aria & Tyrion")
    c.Assert(teamFound, Not(Equals), team)
}
