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