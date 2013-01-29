package model_test

import (
    . "launchpad.net/gocheck"
    "strings"
    "goatd/app/event"
    "goatd/app/model"
)

func queueNames(queues []*model.Queue) (names []string) {
    for _, queue := range queues {
        names = append(names, queue.Name())
    }
    return names
}

type QueueOwner struct {
    *event.Identity
}

type QueueSuite struct {
    store *model.Store
    owner *QueueOwner
}

var _ = Suite(&QueueSuite{})

func (s *QueueSuite) SetUpTest(c *C) {
    s.store = model.NewStore(nil)
    s.owner = &QueueOwner{event.NewIdentity("Team")}
}

func (s *QueueSuite) TestCreateQueue(c *C) {
    queue := s.store.Queues.Create(model.A{"Name": "Sales"}, s.owner)
    c.Assert(queue.Name(), Equals, "Sales")
    c.Assert(queue.TeamUid(), Equals, s.owner.Uid())
}

func (s *QueueSuite) TestFindQueue(c *C) {
    s.store.Queues.Create(model.A{"Name": "One"}, s.owner)
    q2 := s.store.Queues.Create(model.A{"Name": "Two"}, s.owner)
    found := s.store.Queues.Find(q2.Uid())
    c.Assert(found.Name(), DeepEquals, q2.Name())
    c.Assert(s.store.Queues.Find("unknown"), IsNil)
}

func (s *QueueSuite) TestFindAllQueues(c *C) {
    q1 := s.store.Queues.Create(model.A{"Name": "One"}, s.owner)
    s.store.Queues.Create(model.A{"Name": "Two"}, s.owner)
    q3 := s.store.Queues.Create(model.A{"Name": "Three"}, s.owner)
    foundQueues := s.store.Queues.FindAll([]string{q1.Uid(), q3.Uid()})
    c.Assert(queueNames(foundQueues), DeepEquals, []string{"One", "Three"})
}

func (s *QueueSuite) TestUpdateQueue(c *C) {
    queue := s.store.Queues.Create(model.A{"Name": "Jamie Lannister"}, s.owner)
    queue.Update("Name", "Tyrion Lannister")
    c.Assert(queue.Name(), Equals, "Tyrion Lannister")
    found := s.store.Queues.Find(queue.Uid())
    c.Assert(found.Name(), Equals, "Tyrion Lannister")
}

func (s *QueueSuite) TestSelectQueues(c *C) {
    s.store.Queues.Create(model.A{"Name": "Tyrion Lannister"}, s.owner)
    s.store.Queues.Create(model.A{"Name": "Jon Snow"}, s.owner)
    s.store.Queues.Create(model.A{"Name": "Jamie Lannister"}, s.owner)
    selectedQueues := s.store.Queues.Select(func (item interface{}) bool {
        queue := item.(*model.Queue)
        return strings.Contains(queue.Name(), "Lannister")
    })
    c.Assert(queueNames(selectedQueues), DeepEquals, []string{"Tyrion Lannister", "Jamie Lannister"})
}

func (s *QueueSuite) TestAddTask(c *C) {
    caller := s.store.Queues.Create(model.A{"Name": "Caller"})
    c.Assert(caller.AddTask("task1"), Equals, true)
    c.Assert(caller.QueuedTaskUids(), DeepEquals, []string{"task1"})
    c.Assert(caller.AddTask("task2"), Equals, true)
    c.Assert(caller.AddTask("task3"), Equals, true)
    c.Assert(caller.QueuedTaskUids(), DeepEquals, []string{"task1", "task2", "task3"})
    persisted := s.store.Queues.Find(caller.Uid())
    c.Assert(persisted.QueuedTaskUids(), DeepEquals, []string{"task1", "task2", "task3"})
}

func (s *QueueSuite) TestRemoveTask(c *C) {
    caller := s.store.Queues.Create(model.A{"Name": "Caller"})
    caller.AddTask("task1")
    caller.AddTask("task2")
    caller.AddTask("task3")
    caller.RemoveTask("task2")
    c.Assert(caller.QueuedTaskUids(), DeepEquals, []string{"task1", "task3"})
}
