package distribution_test

import (
    . "launchpad.net/gocheck"
    "testing"
    "time"
    "goatd/app/identification"
    "goatd/app/event"
    "goatd/app/models"
    "goatd/app/distribution"
)

func Test(t *testing.T) { TestingT(t) }

type ResourceTrackerSuite struct {
    teams *models.Teams
    team *models.Team
    tracker *distribution.ResourceTracker
    aLittleBit time.Duration
}

var _ = Suite(&ResourceTrackerSuite{})

func (s *ResourceTrackerSuite) SetUpTest(c *C) {
    event.Manager().Start()
    s.teams = models.NewTeams(identification.NoIdentity())
    s.team = s.teams.Create(models.Attrs{"Name": "Get To Work"})
    s.tracker = distribution.NewResourceTracker(s.team)
    s.aLittleBit = 100 * time.Millisecond
}

func (s *ResourceTrackerSuite) TearDownTest(c *C) {
    event.Manager().Stop()
}

func (s *ResourceTrackerSuite) TestTrackTeammateQueues(c *C) {
    agent := s.team.Teammates.Create(models.Attrs{"Name": "John Doe"})
    sales := s.team.Queues.Create(models.Attrs{"Name": "Sales"})
    support := s.team.Queues.Create(models.Attrs{"Name": "Support"})
    s.team.Queues.Create(models.Attrs{"Name": "Legal"})
    s.team.Skills.Create(models.Attrs{"TeammateUid": agent.Uid(), "QueueUid": sales.Uid()})
    s.team.Skills.Create(models.Attrs{"TeammateUid": agent.Uid(), "QueueUid": support.Uid()})
    time.Sleep(s.aLittleBit)

    task := s.team.Tasks.Create(models.Attrs{"Title": "Do It"})
    task.Enqueue(sales)
    queues := s.tracker.TeammateQueuesReady(agent)
    c.Assert(queues, DeepEquals, []*models.Queue{sales})
}
