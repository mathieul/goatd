package distribution_test

import (
    . "launchpad.net/gocheck"
    "testing"
    "goatd/app/distribution"
    "goatd/app/models"
    "time"
)

func Test(t *testing.T) { TestingT(t) }

type DistributorSuite struct {
    teams *models.Teams
    team *models.Team
}

var _ = Suite(&DistributorSuite{})

func (s *DistributorSuite) SetUpTest(c *C) {
    s.teams = models.NewTeams()
    s.team = s.teams.Create(models.Attrs{"Name": "Get To Work"})
}

func (s *DistributorSuite) TestNewDistributor(c *C) {
    distributor := distribution.NewDistributor(s.team)
    c.Assert(distributor.Team(), DeepEquals, s.team)
}

func (s *DistributorSuite) TestBindingEvents(c *C) {
    var mockQueue, mockMate, mockTask string
    aLittleBit := 100 * time.Millisecond

    distributor := distribution.NewDistributor(s.team)
    distributor.On(distribution.EventOfferTask,
        func (event distribution.Event, parameters []interface{}) {
            mockQueue = parameters[0].(string)
            mockMate = parameters[1].(string)
            mockTask = parameters[2].(string)
    })
    distributor.Trigger(distribution.EventOfferTask, []interface{}{"the queue", "the teammate", "the task"})
    time.Sleep(aLittleBit)
    c.Assert(mockQueue, Equals, "the queue")
    c.Assert(mockMate, Equals, "the teammate")
    c.Assert(mockTask, Equals, "the task")
}

func (s *DistributorSuite) TestAddTeammateToQueue(c *C) {
    teammate := s.team.Teammates.Create(models.Attrs{"Name": "The Mate"})
    queue := s.team.Queues.Create(models.Attrs{"Name": "The Queue"})
    distributor := distribution.NewDistributor(s.team)

    distributor.AddTeammateToQueue(queue, teammate, models.LevelLow)
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
