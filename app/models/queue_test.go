package models_test

import (
    . "launchpad.net/gocheck"
    "testing"
    "goatd/app/models"
)

func Test(t *testing.T) { TestingT(t) }

type QueueSuite struct{
}

var _ = Suite(&QueueSuite{})

func (s *QueueSuite) SetUpTest(c *C) {
}

func (s *QueueSuite) TearDownTest(c *C) {
}

func (s *QueueSuite) TestNewQueueWithNameAndTeamUid(c *C) {
    queue := models.NewQueue(models.Attrs{"Name": "Sales", "TeamUid": "inc"})
    c.Assert(queue.Name(), Equals, "Sales")
    c.Assert(queue.TeamUid(), Equals, "inc")
    c.Assert(queue.Persisted(), Equals, false)
}

func (s *QueueSuite) TestCreateQueueWithNameAndTeamUid(c *C) {
    queue := models.CreateQueue(models.Attrs{"Name": "Support", "TeamUid": "inc"})
    c.Assert(queue.Name(), Equals, "Support")
    c.Assert(queue.TeamUid(), Equals, "inc")
    c.Assert(queue.Persisted(), Equals, true)
}
