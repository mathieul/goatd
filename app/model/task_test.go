package model_test

import (
    . "launchpad.net/gocheck"
    "strings"
    "goatd/app/event"
    "goatd/app/model"
)

func taskTitles(tasks []*model.Task) (titles []string) {
    for _, task := range tasks {
        titles = append(titles, task.Title())
    }
    return titles
}

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
    found := s.store.Tasks.Find(q2.Uid())
    c.Assert(found.Title(), DeepEquals, q2.Title())
    c.Assert(s.store.Tasks.Find("unknown"), IsNil)
}

func (s *TaskSuite) TestFindAllTasks(c *C) {
    q1 := s.store.Tasks.Create(model.A{"Title": "One"}, s.owner)
    s.store.Tasks.Create(model.A{"Title": "Two"}, s.owner)
    q3 := s.store.Tasks.Create(model.A{"Title": "Three"}, s.owner)
    foundTasks := s.store.Tasks.FindAll([]string{q1.Uid(), q3.Uid()})
    c.Assert(taskTitles(foundTasks), DeepEquals, []string{"One", "Three"})
}

func (s *TaskSuite) TestUpdateTask(c *C) {
    task := s.store.Tasks.Create(model.A{"Title": "Jamie Lannister"}, s.owner)
    task.Update("Title", "Tyrion Lannister")
    c.Assert(task.Title(), Equals, "Tyrion Lannister")
    found := s.store.Tasks.Find(task.Uid())
    c.Assert(found.Title(), Equals, "Tyrion Lannister")
}

func (s *TaskSuite) TestSelectTasks(c *C) {
    s.store.Tasks.Create(model.A{"Title": "Tyrion Lannister"}, s.owner)
    s.store.Tasks.Create(model.A{"Title": "Jon Snow"}, s.owner)
    s.store.Tasks.Create(model.A{"Title": "Jamie Lannister"}, s.owner)
    selectedTasks := s.store.Tasks.Select(func (item interface{}) bool {
        task := item.(*model.Task)
        return strings.Contains(task.Title(), "Lannister")
    })
    c.Assert(taskTitles(selectedTasks), DeepEquals, []string{"Tyrion Lannister", "Jamie Lannister"})
}

func (s *TaskSuite) TestEnqueueTask(c *C) {
    queueUid := "abcd1234"
    task := s.store.Tasks.Create(model.A{"Title": "Clean-up my room"}, s.owner)
    c.Assert(task.Status(), Equals, model.StatusCreated)
    c.Assert(task.Enqueue(queueUid), Equals, true)
    c.Assert(task.Status(), Equals, model.StatusQueued)
    c.Assert(task.QueueUid(), Equals, queueUid)
    c.Assert(task.Enqueue(queueUid), Equals, false)
}

func (s *TaskSuite) TestDequeueTask(c *C) {
    queueUid := "abcd1234"
    task := s.store.Tasks.Create(model.A{"Title": "Clean-up my room"}, s.owner)
    task.Enqueue(queueUid)
    c.Assert(task.Dequeue("blah"), Equals, false)
    c.Assert(task.Dequeue(queueUid), Equals, true)
    c.Assert(task.Status(), Equals, model.StatusCreated)
    c.Assert(task.QueueUid(), Equals, "")
    c.Assert(task.Enqueue(queueUid), Equals, true)
}
