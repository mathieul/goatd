package distribution_test

import (
    . "launchpad.net/gocheck"
    "testing"
    "goatd/app/distribution"
    "goatd/app/models"
    "time"
)

func Test(t *testing.T) { TestingT(t) }

type DistributorSuite struct{
}

var _ = Suite(&DistributorSuite{})

func (s *DistributorSuite) SetUpTest(c *C) {
}

func (s *DistributorSuite) TearDownTest(c *C) {
}

func (s *DistributorSuite) TestNewDistributor(c *C) {
    team := models.NewTeam(models.Attrs{"Name": "Get To Work"})
    distributor := distribution.NewDistributor(team)
    c.Assert(distributor.Team(), DeepEquals, team)
}

func (s *DistributorSuite) TestBindingEvents(c *C) {
    var mockQueue, mockMate, mockTask string
    aLittleBit := 100 * time.Millisecond

    distributor := distribution.NewDistributor(nil)
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
    team := models.NewTeam(models.Attrs{"Name": "The Team"})
    teammate := s.team.Teammates.Create(models.Attrs{"Name": "The Mate"})
    queue := s.team.Queues.Create(models.Attrs{"Name": "The Queue"})

    distributor := distribution.NewDistributor(team)
    distributor.AddTeammateToQueue(queue, teammate, models.LevelLow)
    c.Assert(queue.Teammates().Slice(), DeepEquals, []models.Teammate{teammate})
    c.Assert(teammate.Queues().Slice(), DeepEquals, []models.Queue{queue})
    query := models.Attrs{"TeammateUid": teammate.Uid(), "QueueUid": queue.Uid()}
    skills := team.Skills().Select(query)
    c.Assert(len(skills), Equals, 1)
    c.Assert(skills[0].QueueUid(), Equals, queue.Uid())
    c.Assert(skills[0].TeammateUid(), Equals, teammate.Uid())
    c.Assert(skills[0].Level(), Equals, models.LevelLow)
    c.Assert(skills[0].Enabled(), Equals, true)
}
