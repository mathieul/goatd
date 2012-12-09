package models_test

import (
    . "launchpad.net/gocheck"
    "testing"
    "strings"
    "goatd/app.old/models"
)

func Test(t *testing.T) { TestingT(t) }

type TeammateSuite struct{
    team *models.Team
    teammates *models.Teammates
    agent *models.Teammate
}

var _ = Suite(&TeammateSuite{})

func (s *TeammateSuite) SetUpTest(c *C) {
    s.team = models.NewTeam(models.Attrs{"Name": "Jon Snow & Egret"})
    s.teammates = s.team.Teammates
    s.agent = s.teammates.Create(models.Attrs{"Name": "Agent"})
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
    c.Assert(s.teammates.Slice(), DeepEquals, []*models.Teammate{s.agent, t1, t2, t3})
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
    c.Assert(s.agent.Status(), Equals, models.StatusSignedOut)
    c.Assert(s.agent.SignIn(), Equals, true)
    c.Assert(s.agent.Status(), Equals, models.StatusOnBreak)
    c.Assert(s.agent.SignOut(), Equals, true)
    c.Assert(s.agent.Status(), Equals, models.StatusSignedOut)
}

func (s *TeammateSuite) TestChangingAvailability(c *C) {
    s.agent.SignIn()
    c.Assert(s.agent.MakeAvailable(), Equals, true)
    c.Assert(s.agent.Status(), Equals, models.StatusWaiting)
    task := s.team.Tasks.Create(models.Attrs{"Title": "Do It"})
    c.Assert(s.agent.OfferTask(task), Equals, true)
    c.Assert(s.agent.Status(), Equals, models.StatusOffered)
    c.Assert(s.agent.CurrentTask(), DeepEquals, task)
}

func (s *TeammateSuite) TestAcceptFinishTask(c *C) {
    s.agent.SignIn()
    s.agent.MakeAvailable()
    task := s.team.Tasks.Create(models.Attrs{"Title": "Do It"})
    s.agent.OfferTask(task)
    c.Assert(s.agent.AcceptTask(task), Equals, true)
    c.Assert(s.agent.Status(), Equals, models.StatusBusy)
    c.Assert(s.agent.CurrentTask(), DeepEquals, task)

    c.Assert(s.agent.FinishTask(task), Equals, true)
    c.Assert(s.agent.Status(), Equals, models.StatusWrappingUp)
    c.Assert(s.agent.CurrentTask(), IsNil)
}

func (s *TeammateSuite) TestOtherWorkOnBreakTask(c *C) {
    s.agent.SignIn()
    s.agent.MakeAvailable()
    task := s.team.Tasks.Create(models.Attrs{"Title": "Do It"})
    s.agent.OfferTask(task)
    s.agent.AcceptTask(task)
    c.Assert(s.agent.StartOtherWork(), Equals, false)
    s.agent.FinishTask(task)

    c.Assert(s.agent.StartOtherWork(), Equals, true)
    c.Assert(s.agent.Status(), Equals, models.StatusOtherWork)
    c.Assert(s.agent.GoOnBreak(), Equals, true)
    c.Assert(s.agent.Status(), Equals, models.StatusOnBreak)
    c.Assert(s.agent.StartOtherWork(), Equals, true)
    c.Assert(s.agent.Status(), Equals, models.StatusOtherWork)
}
