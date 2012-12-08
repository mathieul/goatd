package models_test

import (
    . "launchpad.net/gocheck"
    "testing"
    "goatd/app/models"
    "strings"
    "fmt"
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
    c.Assert(s.queues.Find("unknown"), IsNil)
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
    tyrion := s.queues.Create(models.Attrs{"Name": "Tyrion Lannister"})
    s.queues.Create(models.Attrs{"Name": "Jon Snow"})
    jamie := s.queues.Create(models.Attrs{"Name": "Jamie Lannister"})
    c.Assert(s.queues.Select(func (item interface{}) bool {
            queue := item.(*models.Queue)
            return strings.Contains(queue.Name(), "Lannister")
        }),
        DeepEquals,
        []*models.Queue{tyrion, jamie})
}

func (s *QueueSuite) TestInsertTask(c *C) {
    caller := s.queues.Create(models.Attrs{"Name": "Caller"})
    c.Assert(caller.IsReady(), Equals, false)
    task := s.team.Tasks.Create(models.Attrs{"Title": "Do It"})
    c.Assert(caller.InsertTask(task), Equals, true)
    c.Assert(caller.IsReady(), Equals, true)
    c.Assert(caller.QueuedTasks(), DeepEquals, []*models.Task{task})
}

func (s *QueueSuite) TestInsertAndKeepTasksOrdered(c *C) {
    caller := s.queues.Create(models.Attrs{"Name": "Caller"})
    tasks := make([]*models.Task, 0, 10)
    for i := 0; i < 10; i++ {
        task := s.team.Tasks.Create(models.Attrs{"Title": fmt.Sprintf("#%02d", i)})
        caller.InsertTask(task)
        tasks = append(tasks, task)
    }
    c.Assert(caller.QueuedTasks(), DeepEquals, tasks)
}
