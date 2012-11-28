package models_test

import (
    . "launchpad.net/gocheck"
    "testing"
    "goatd/app/models"
)

func Test(t *testing.T) { TestingT(t) }

type TeamSuite struct{
}

var _ = Suite(&TeamSuite{})

func (s *TeamSuite) SetUpTest(c *C) {
}

func (s *TeamSuite) TearDownTest(c *C) {
}

func (s *TeamSuite) TestNewTeamWithName(c *C) {
    team := models.NewTeam(models.Attrs{"Name": "Jon Snow & Egret"})
    c.Assert(team.Name(), Equals, "Jon Snow & Egret")
    c.Assert(len(team.Uid()), Equals, 8 + 1 + 8)
    c.Assert(team.Persisted(), Equals, false)
}

func (s *TeamSuite) TestCreateTeamWithName(c *C) {
    team := models.CreateTeam(models.Attrs{"Name": "Metallica"})
    c.Assert(team.Name(), Equals, "Metallica")
    c.Assert(team.Persisted(), Equals, true)
}
