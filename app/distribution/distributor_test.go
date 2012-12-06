package distribution_test

import (
    . "launchpad.net/gocheck"
    "testing"
    "goatd/app/identification"
    "goatd/app/distribution"
    "goatd/app/models"
    "strings"
)

func Test(t *testing.T) { TestingT(t) }

type DistributorSuite struct {
    teams *models.Teams
    team *models.Team
}

var _ = Suite(&DistributorSuite{})

func (s *DistributorSuite) SetUpTest(c *C) {
    s.teams = models.NewTeams(identification.NoIdentity())
    s.team = s.teams.Create(models.Attrs{"Name": "Get To Work"})
}

func (s *DistributorSuite) TestNewDistributor(c *C) {
    distributor := distribution.NewDistributor(s.team)
    c.Assert(distributor.Team(), DeepEquals, s.team)
}

func (s *DistributorSuite) TestAddTeammateToQueue(c *C) {
    teammate := s.team.Teammates.Create(models.Attrs{"Name": "The Mate"})
    queue := s.team.Queues.Create(models.Attrs{"Name": "The Queue"})
    distributor := distribution.NewDistributor(s.team)

    result := distributor.AddTeammateToQueue(queue, teammate, models.LevelLow)
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

func (s *DistributorSuite) TestEnqueueTask(c *C) {
    distributor := distribution.NewDistributor(s.team)
    queue := s.team.Queues.Create(models.Attrs{"Name": "Support"})
    task := s.team.Tasks.Create(models.Attrs{"Title": "My printer is not working"})

    result := distributor.EnqueueTask(queue, task, models.PriorityMedium)
    c.Assert(result, Equals, true)
    c.Assert(task.Status(), Equals, models.StatusQueued)
    c.Assert(task.Priority(), Equals, models.PriorityMedium)
    queuedTasks := s.team.Tasks.Select(func (item interface{}) bool {
        task := item.(*models.Task)
        return strings.Contains(task.QueueUid(), queue.Uid())
    })
    c.Assert(queuedTasks, DeepEquals, []*models.Task{task})
    result = distributor.EnqueueTask(queue, task, models.PriorityHigh)
    c.Assert(result, Equals, false)
}
