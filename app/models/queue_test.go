package models_test

import (
    . "launchpad.net/gocheck"
    "testing"
    "goatd/app/models"
    "strings"
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

func (s *QueueSuite) TestFindQueue(c *C) {
    s.queues.Create(models.Attrs{"Name": "One"})
    q2 := s.queues.Create(models.Attrs{"Name": "Two"})
    c.Assert(s.queues.Find(q2.Uid()), DeepEquals, q2)
}

func (s *QueueSuite) TestFindQueueSlice(c *C) {
    q1 := s.queues.Create(models.Attrs{"Name": "One"})
    s.queues.Create(models.Attrs{"Name": "Two"})
    q3 := s.queues.Create(models.Attrs{"Name": "Three"})
    c.Assert(s.queues.FindAll([]string{q1.Uid(), q3.Uid()}),
             DeepEquals,
             []*models.Queue{q1, q3})
}

func (s *QueueSuite) TestSelectQueues(c *C) {
    tyrion := s.queues.Create(models.Attrs{"Name": "Tyrion Lanister"})
    s.queues.Create(models.Attrs{"Name": "Jon Snow"})
    jamie := s.queues.Create(models.Attrs{"Name": "Jamie Lanister"})
    c.Assert(s.queues.Select(func (queue *models.Queue) bool {
            if strings.Contains(queue.Name(), "Lanister") {
                return true
            }
            return false
        }),
        DeepEquals,
        []*models.Queue{tyrion, jamie})
}
