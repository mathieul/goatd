package models_test

import (
    . "launchpad.net/gocheck"
    "testing"
    "goatd/app/models"
)

func Test(t *testing.T) { TestingT(t) }

type QueueSuite struct {
    queues *models.Queues
}

var _ = Suite(&QueueSuite{})

func (s *QueueSuite) SetUpTest(c *C) {
    s.queues = models.NewQueues("Team", "inc")
}

func (s *QueueSuite) TestCreateQueue(c *C) {
    queue := s.queues.Create(models.Attrs{"Name": "Sales"})
    c.Assert(queue.Name(), Equals, "Sales")
    c.Assert(queue.TeamUid(), Equals, "inc")
    c.Assert(queue.Persisted(), Equals, true)
}
