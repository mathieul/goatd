package distribution_test

import (
    . "launchpad.net/gocheck"
    "testing"
    "time"
    "goatd/app/identification"
    "goatd/app/models"
    "goatd/app/event"
    "goatd/app/distribution"
)

func Test(t *testing.T) { TestingT(t) }

type DistributorSuite struct {
    teams *models.Teams
    team *models.Team
    distributor *distribution.Distributor
    aLittleBit time.Duration
}

var _ = Suite(&DistributorSuite{})

func (s *DistributorSuite) SetUpTest(c *C) {
    event.Manager().Start()
    s.aLittleBit = 100 * time.Millisecond
    s.teams = models.NewTeams(identification.NoIdentity())
    s.team = s.teams.Create(models.Attrs{"Name": "Get To Work"})
    s.distributor = distribution.NewDistributor(s.team)
}

func (s *DistributorSuite) TearDownTest(c *C) {
    event.Manager().Stop()
}

func (s *DistributorSuite) TestNewDistributor(c *C) {
    c.Assert(s.distributor.Team(), DeepEquals, s.team)
}

func (s *DistributorSuite) TestAddTeammateToQueue(c *C) {
    teammate := s.team.Teammates.Create(models.Attrs{"Name": "The Mate"})
    queue := s.team.Queues.Create(models.Attrs{"Name": "The Queue"})

    result := s.distributor.AddTeammateToQueue(queue, teammate, models.LevelLow)
    c.Assert(result, Equals, true)
    teammateUid, queueUid := teammate.Uid(), queue.Uid()
    skills := s.team.Skills.Select(func (item interface{}) bool {
        skill := item.(*models.Skill)
        return skill.TeammateUid() == teammateUid && skill.QueueUid() == queueUid
    })
    c.Assert(len(skills), Equals, 1)
    c.Assert(skills[0].QueueUid(), Equals, queue.Uid())
    c.Assert(skills[0].TeammateUid(), Equals, teammate.Uid())
    c.Assert(skills[0].Level(), Equals, models.LevelLow)
    c.Assert(skills[0].Enabled(), Equals, true)
}

func (s *DistributorSuite) TestAssignTaskWhenTeammateAvailable(c *C) {
    teammate := s.team.Teammates.Create(models.Attrs{"Name": "The Mate"})
    queue := s.team.Queues.Create(models.Attrs{"Name": "The Queue"})
    task := s.team.Tasks.Create(models.Attrs{"Title": "Do It"})
    task.Enqueue(queue)
    teammate.SignIn()

    teammate.MakeAvailable()
    time.Sleep(s.aLittleBit)
    c.Assert(teammate.Status(), Equals, models.StatusOffered)
    c.Assert(task.Status(), Equals, models.StatusOffered)
}
