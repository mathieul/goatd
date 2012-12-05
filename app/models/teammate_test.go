package models_test

import (
    . "launchpad.net/gocheck"
    "testing"
    "strings"
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

func (s *TeammateSuite) TestReturnSlice(c *C) {
    t1 := s.teammates.Create(models.Attrs{"Name": "Jon"})
    t2 := s.teammates.Create(models.Attrs{"Name": "Egret"})
    t3 := s.teammates.Create(models.Attrs{"Name": "Aria"})
    c.Assert(s.teammates.Slice(), DeepEquals, []*models.Teammate{t1, t2,t3})
}

func (s *TeammateSuite) TestFindTeammate(c *C) {
    s.teammates.Create(models.Attrs{"Name": "Jon"})
    egret := s.teammates.Create(models.Attrs{"Name": "Egret"})
    c.Assert(s.teammates.Find(egret.Uid()), DeepEquals, egret)
    c.Assert(s.teammates.Find("unknown"), IsNil)
}

func (s *TeammateSuite) TestFindTeammateSlice(c *C) {
    t1 := s.teammates.Create(models.Attrs{"Name": "One"})
    s.teammates.Create(models.Attrs{"Name": "Two"})
    t3 := s.teammates.Create(models.Attrs{"Name": "Three"})
    c.Assert(s.teammates.FindAll([]string{t1.Uid(), t3.Uid()}),
             DeepEquals,
             []*models.Teammate{t1, t3})
}

func (s *TeammateSuite) TestSelectTeammates(c *C) {
    tyrion := s.teammates.Create(models.Attrs{"Name": "Tyrion Lannister"})
    s.teammates.Create(models.Attrs{"Name": "Jon Snow"})
    jamie := s.teammates.Create(models.Attrs{"Name": "Jamie Lannister"})
    c.Assert(s.teammates.Select(func (item interface{}) bool {
            teammate := item.(*models.Teammate)
            return strings.Contains(teammate.Name(), "Lannister")
        }),
        DeepEquals,
        []*models.Teammate{tyrion, jamie})
}

func (s *TeammateSuite) TestSignInSignOutTeammate(c *C) {
    agent := s.teammates.Create(models.Attrs{"Name": "Agent"})
    c.Assert(agent.Status(), Equals, models.StatusSignedOut)
    agent.SignIn()
    c.Assert(agent.Status(), Equals, models.StatusOnBreak)
    agent.SignOut()
    c.Assert(agent.Status(), Equals, models.StatusSignedOut)
}
