package model_test

import (
    . "launchpad.net/gocheck"
    // "strings"
    "goatd/app/event"
    "goatd/app/model"
)

type TaskOwner struct {
    *event.Identity
}

type TaskSuite struct{
    store *model.Store
    owner *TaskOwner
}

var _ = Suite(&TaskSuite{})

func (s *TaskSuite) SetUpTest(c *C) {
    s.store = model.NewStore()
    s.owner = &TaskOwner{event.NewIdentity("Team")}
}

func (s *TaskSuite) TestCreateTask(c *C) {
    task := s.store.Tasks.Create(model.A{"Title": "Call back the milkman"}, s.owner)
    c.Assert(task.Title(), Equals, "Call back the milkman")
    c.Assert(task.TeamUid(), Equals, s.owner.Uid())
}

// func (s *TaskSuite) TestTaskHasPriority(c *C) {
//     task := s.tasks.Create(models.Attrs{"Title": "Blah"})
//     c.Assert(task.Priority(), Equals, models.PriorityMedium)
//     task.SetPriority(models.PriorityHigh)
//     c.Assert(task.Priority(), Equals, models.PriorityHigh)
// }

// func (s *TaskSuite) TestReturnsSlice(c *C) {
//     q1 := s.tasks.Create(models.Attrs{"Title": "One"})
//     q2 := s.tasks.Create(models.Attrs{"Title": "Two"})
//     q3 := s.tasks.Create(models.Attrs{"Title": "Three"})
//     c.Assert(s.tasks.Slice(), DeepEquals, []*models.Task{q1, q2, q3})
// }

// func (s *TaskSuite) TestFindTask(c *C) {
//     s.tasks.Create(models.Attrs{"Title": "One"})
//     q2 := s.tasks.Create(models.Attrs{"Title": "Two"})
//     c.Assert(s.tasks.Find(q2.Uid()), DeepEquals, q2)
//     c.Assert(s.tasks.Find("unknown"), IsNil)
// }

// func (s *TaskSuite) TestFindTaskSlice(c *C) {
//     q1 := s.tasks.Create(models.Attrs{"Title": "One"})
//     s.tasks.Create(models.Attrs{"Title": "Two"})
//     q3 := s.tasks.Create(models.Attrs{"Title": "Three"})
//     c.Assert(s.tasks.FindAll([]string{q1.Uid(), q3.Uid()}),
//              DeepEquals,
//              []*models.Task{q1, q3})
// }

// func (s *TaskSuite) TestSelectTasks(c *C) {
//     tyrion := s.tasks.Create(models.Attrs{"Title": "Tyrion Lannister"})
//     s.tasks.Create(models.Attrs{"Title": "Jon Snow"})
//     jamie := s.tasks.Create(models.Attrs{"Title": "Jamie Lannister"})
//     c.Assert(s.tasks.Select(func (item interface{}) bool {
//             task := item.(*models.Task)
//             return strings.Contains(task.Title(), "Lannister")
//         }),
//         DeepEquals,
//         []*models.Task{tyrion, jamie})
// }

// func (s *TaskSuite) TestEnqueueTask(c *C) {
//     queue := s.team.Queues.Create(models.Attrs{"Name": "My TODOs"})
//     task := s.tasks.Create(models.Attrs{"Title": "Clean-up my room"})
//     c.Assert(task.Status(), Equals, models.StatusCreated)
//     c.Assert(task.Enqueue(queue), Equals, true)
//     c.Assert(task.Status(), Equals, models.StatusQueued)
//     c.Assert(task.QueueUid(), Equals, queue.Uid())
//     c.Assert(queue.QueuedTasks(), DeepEquals, []*models.Task{task})

//     c.Assert(task.Enqueue(queue), Equals, false)
// }
