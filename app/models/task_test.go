package models_test

import (
    . "launchpad.net/gocheck"
    "testing"
    "goatd/app/models"
    "strings"
)

func Test(t *testing.T) { TestingT(t) }

type TaskSuite struct{
    team *models.Team
    tasks *models.Tasks
}

var _ = Suite(&TaskSuite{})

func (s *TaskSuite) SetUpTest(c *C) {
    s.team = models.NewTeam(models.Attrs{"Name": "Tyrion's Team"})
    s.tasks = s.team.Tasks
}

func (s *TaskSuite) TestCreateTaskWithTitle(c *C) {
    task := s.tasks.Create(models.Attrs{"Title": "Call back the milkman"})
    c.Assert(task.Title(), Equals, "Call back the milkman")
    c.Assert(task.Persisted(), Equals, true)
}

func (s *TaskSuite) TestTaskHasPriority(c *C) {
    task := s.tasks.Create(models.Attrs{"Title": "Blah"})
    c.Assert(task.Priority(), Equals, models.PriorityMedium)
    task.SetPriority(models.PriorityHigh)
    c.Assert(task.Priority(), Equals, models.PriorityHigh)
}

func (s *TaskSuite) TestReturnsSlice(c *C) {
    q1 := s.tasks.Create(models.Attrs{"Title": "One"})
    q2 := s.tasks.Create(models.Attrs{"Title": "Two"})
    q3 := s.tasks.Create(models.Attrs{"Title": "Three"})
    c.Assert(s.tasks.Slice(), DeepEquals, []*models.Task{q1, q2, q3})
}

func (s *TaskSuite) TestFindTask(c *C) {
    s.tasks.Create(models.Attrs{"Title": "One"})
    q2 := s.tasks.Create(models.Attrs{"Title": "Two"})
    c.Assert(s.tasks.Find(q2.Uid()), DeepEquals, q2)
    c.Assert(s.tasks.Find("unknown"), IsNil)
}

func (s *TaskSuite) TestFindTaskSlice(c *C) {
    q1 := s.tasks.Create(models.Attrs{"Title": "One"})
    s.tasks.Create(models.Attrs{"Title": "Two"})
    q3 := s.tasks.Create(models.Attrs{"Title": "Three"})
    c.Assert(s.tasks.FindAll([]string{q1.Uid(), q3.Uid()}),
             DeepEquals,
             []*models.Task{q1, q3})
}

func (s *TaskSuite) TestSelectTasks(c *C) {
    tyrion := s.tasks.Create(models.Attrs{"Title": "Tyrion Lannister"})
    s.tasks.Create(models.Attrs{"Title": "Jon Snow"})
    jamie := s.tasks.Create(models.Attrs{"Title": "Jamie Lannister"})
    c.Assert(s.tasks.Select(func (item interface{}) bool {
            task := item.(*models.Task)
            return strings.Contains(task.Title(), "Lannister")
        }),
        DeepEquals,
        []*models.Task{tyrion, jamie})
}

func (s *TaskSuite) TestSignInSignOutTask(c *C) {
    queue := s.team.Queues.Create(models.Attrs{"Name": "My TODOs"})
    task := s.tasks.Create(models.Attrs{"Title": "Clean-up my room"})
    c.Assert(task.Status(), Equals, models.StatusCreated)
    c.Assert(task.Queue(queue), Equals, true)
    c.Assert(task.Status(), Equals, models.StatusQueued)
    c.Assert(task.QueueUid(), Equals, queue.Uid())
    c.Assert(task.Queue(queue), Equals, false)
}
