package models_test

import (
    . "launchpad.net/gocheck"
    "testing"
    "goatd/app/models"
)

func Test(t *testing.T) { TestingT(t) }

type TeammateSuite struct{
    team *models.Team
    teammates *models.Teammates
}

var _ = Suite(&TeammateSuite{})

func (s *TeammateSuite) SetUpTest(c *C) {
    s.team = models.NewTeam(models.Attrs{"Name": "Jon Snow & Egret"})
    s.teammates = s.team.Teammates
}

func (s *TeammateSuite) TestCreateTeammate(c *C) {
    teammate := s.teammates.Create(models.Attrs{"Name": "Jon"})
    c.Assert(teammate.Name(), Equals, "Jon")
    c.Assert(teammate.TeamUid(), Equals, s.team.Uid())
    c.Assert(teammate.Persisted(), Equals, true)
}

func (s *TeammateSuite) TestCreateTeammateWithTeam(c *C) {
    teammate := s.teammates.Create(models.Attrs{"Name": "Egret"})
    c.Assert(teammate.Team(), DeepEquals, s.team)
}
