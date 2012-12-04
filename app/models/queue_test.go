package models_test

import (
    . "launchpad.net/gocheck"
    "testing"
    "goatd/app/models"
)

func Test(t *testing.T) { TestingT(t) }

type QueueSuite struct {
    team *models.Team
    queues *models.Queues
}

var _ = Suite(&QueueSuite{})

func (s *QueueSuite) SetUpTest(c *C) {
    s.team = models.NewTeam(models.Attrs{"Name": "Hello, inc."})
    s.queues = s.team.Queues
}

func (s *QueueSuite) TestCreateQueue(c *C) {
    queue := s.queues.Create(models.Attrs{"Name": "Sales"})
    c.Assert(queue.Name(), Equals, "Sales")
    c.Assert(queue.TeamUid(), Equals, s.team.Uid())
    c.Assert(queue.Persisted(), Equals, true)
}

func (s *QueueSuite) TestReturnsSlice(c *C) {
    q1 := s.queues.Create(models.Attrs{"Name": "One"})
    q2 := s.queues.Create(models.Attrs{"Name": "Two"})
    q3 := s.queues.Create(models.Attrs{"Name": "Three"})
    c.Assert(s.queues.Slice(), DeepEquals, []*models.Queue{q1, q2, q3})
}

// func (s *QueueSuite) TestReturnAssignedTeammates(c *C) {
//     queue := s.queues.Create(models.Attrs{"Name": "James Bond"})
//     craig := t.team.Teammates.Create(models.Attrs{"Name": "Daniel Craig"})
//     caine := t.team.Teammates.Create(models.Attrs{"Name": "Michael Caine"})
//     connery := t.team.Teammates.Create(models.Attrs{"Name": "Sean Connery"})
//     c.Assert(queue.Teammates(), DeepEquals, []*models.Teammate{craig, connery})
// }
