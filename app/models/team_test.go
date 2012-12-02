package models_test

import (
    . "launchpad.net/gocheck"
    "testing"
    "goatd/app/models"
)

func Test(t *testing.T) { TestingT(t) }

type TeamSuite struct{
    teams *models.Teams
    team *models.Team
}

var _ = Suite(&TeamSuite{})

func (s *TeamSuite) SetUpTest(c *C) {
    s.teams = models.NewTeams()
    s.team = s.teams.Create(models.Attrs{"Name": "Jon Snow & Egret"})
}

func (s *TeamSuite) TestCreateTeam(c *C) {
    c.Assert(s.team.Name(), Equals, "Jon Snow & Egret")
    c.Assert(len(s.team.Uid()), Equals, 8 + 1 + 8)
    c.Assert(s.team.Persisted(), Equals, true)
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
