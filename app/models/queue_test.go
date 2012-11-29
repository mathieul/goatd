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
    teammate := models.NewQueue(models.Attrs{"Name": "Sales", "TeamUid": "inc"})
    c.Assert(teammate.Name(), Equals, "Sales")
    c.Assert(teammate.TeamUid(), Equals, "inc")
    c.Assert(teammate.Persisted(), Equals, false)
}

func (s *QueueSuite) TestCreateQueueWithNameAndTeamUid(c *C) {
    teammate := models.CreateQueue(models.Attrs{"Name": "Support", "TeamUid": "inc"})
    c.Assert(teammate.Name(), Equals, "Support")
    c.Assert(teammate.TeamUid(), Equals, "inc")
    c.Assert(teammate.Persisted(), Equals, true)
}
