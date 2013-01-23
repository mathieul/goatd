package model_test

import (
    . "launchpad.net/gocheck"
    "strings"
    "goatd/app/event"
    "goatd/app/model"
)

type TaskOwner struct {
    *event.Identity
}

type TaskSuite struct {
    store *model.Store
    owner *TaskOwner
}

var _ = Suite(&TaskSuite{})

func (s *TaskSuite) SetUpTest(c *C) {
    s.store = model.NewStore(nil)
    s.owner = &TaskOwner{event.NewIdentity("Team")}
}

func (s *TaskSuite) TestCreateTask(c *C) {
    task := s.store.Tasks.Create(model.A{"Title": "Call back the milkman"}, s.owner)
    c.Assert(task.Title(), Equals, "Call back the milkman")
    c.Assert(task.TeamUid(), Equals, s.owner.Uid())
}

func (s *TaskSuite) TestFindTask(c *C) {
    s.store.Tasks.Create(model.A{"Title": "One"}, s.owner)
    q2 := s.store.Tasks.Create(model.A{"Title": "Two"}, s.owner)
    c.Assert(s.store.Tasks.Find(q2.Uid()), DeepEquals, q2)
    c.Assert(s.store.Tasks.Find("unknown"), IsNil)
}

func (s *TaskSuite) TestFindTaskSlice(c *C) {
    q1 := s.store.Tasks.Create(model.A{"Title": "One"}, s.owner)
    s.store.Tasks.Create(model.A{"Title": "Two"}, s.owner)
    q3 := s.store.Tasks.Create(model.A{"Title": "Three"}, s.owner)
    c.Assert(s.store.Tasks.FindAll([]string{q1.Uid(), q3.Uid()}),
             DeepEquals,
             []*model.Task{q1, q3})
}

func (s *TaskSuite) TestSelectTasks(c *C) {
    tyrion := s.store.Tasks.Create(model.A{"Title": "Tyrion Lannister"}, s.owner)
    s.store.Tasks.Create(model.A{"Title": "Jon Snow"}, s.owner)
    jamie := s.store.Tasks.Create(model.A{"Title": "Jamie Lannister"}, s.owner)
    c.Assert(s.store.Tasks.Select(func (item interface{}) bool {
            task := item.(*model.Task)
            return strings.Contains(task.Title(), "Lannister")
        }),
        DeepEquals,
        []*model.Task{tyrion, jamie})
}

// func (s *TaskSuite) TestEnqueueTask(c *C) {
//     queue := s.team.Queues.Create(model.A{"Name": "My TODOs"})
//     task := s.store.Tasks.Create(model.A{"Title": "Clean-up my room"}, s.owner)
//     c.Assert(task.Status(), Equals, model.StatusCreated)
//     c.Assert(task.Enqueue(queue), Equals, true)
//     c.Assert(task.Status(), Equals, model.StatusQueued)
//     c.Assert(task.QueueUid(), Equals, queue.Uid())
//     c.Assert(queue.QueuedTasks(), DeepEquals, []*model.Task{task})

//     c.Assert(task.Enqueue(queue), Equals, false)
// }
