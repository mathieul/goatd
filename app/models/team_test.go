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

func (s *TeamSuite) TestCreateQueue(c *C) {
    queue := s.team.Queues.Create(models.Attrs{"Name": "Support"})
    c.Assert(queue.Name(), Equals, "Support")
    c.Assert(queue.TeamUid(), Equals, s.team.Uid())
    c.Assert(queue.Persisted(), Equals, true)
}

func (s *TeamSuite) TestCreateSkill(c *C) {
    skill := s.team.Skills.Create(models.Attrs{"QueueUid": "0abc", "TeammateUid": "1def"})
    c.Assert(skill.QueueUid(), Equals, "0abc")
    c.Assert(skill.TeammateUid(), Equals, "1def")
    c.Assert(skill.Persisted(), Equals, true)
}
