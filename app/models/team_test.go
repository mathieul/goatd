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

func (s *TeamSuite) TestFindTeam(c *C) {
    uid := s.teams.Create(models.Attrs{"Name": "Metallica"}).Uid()
    s.teams.Create(models.Attrs{"Name": "Masada"})
    team := s.teams.Find(uid)
    c.Assert(team.Name(), Equals, "Metallica")
    team = s.teams.Find("unknown")
    c.Assert(team, IsNil)
}

func (s *TeamSuite) TestFindAllTeams(c *C) {
    uid1 := s.teams.Create(models.Attrs{"Name": "One"}).Uid()
    s.teams.Create(models.Attrs{"Name": "Two"})
    s.teams.Create(models.Attrs{"Name": "Three"})
    uid2 := s.teams.Create(models.Attrs{"Name": "Four"}).Uid()
    s.teams.Create(models.Attrs{"Name": "Five"})

    found := s.teams.FindAll([]string{uid1, uid2})
    c.Assert(found[0].Name(), DeepEquals, "One")
    c.Assert(found[1].Name(), DeepEquals, "Four")
}

func (s *TeamSuite) TestCreateTeammate(c *C) {
    teammate := s.team.Teammates.Create(models.Attrs{"Name": "Kirk Hammett"})
    c.Assert(teammate.Name(), Equals, "Kirk Hammett")
    c.Assert(teammate.TeamUid(), Equals, s.team.Uid())
    c.Assert(teammate.Team(), DeepEquals, s.team)
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
