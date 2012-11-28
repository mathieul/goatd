package models_test

import (
    . "launchpad.net/gocheck"
    "testing"
    "goatd/app/models"
)

func Test(t *testing.T) { TestingT(t) }

type TeamSuite struct{
	team *models.Team
}

var _ = Suite(&TeamSuite{})

func (s *TeamSuite) SetUpTest(c *C) {
	s.team = models.CreateTeam(models.Attrs{"Name": "MyTeam"})
}

func (s *TeamSuite) TearDownTest(c *C) {
	s.team = nil
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

func (s *TeamSuite) TestCreateTeammate(c *C) {
	teammate := s.team.Teammates.Create(models.Attrs{"Name": "Kirk Hammett"})
    c.Assert(teammate.Name(), Equals, "Kirk Hammett")
    c.Assert(teammate.TeamUid(), Equals, s.team.Uid())
    c.Assert(teammate.Persisted(), Equals, true)
}
