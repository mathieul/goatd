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

func (s *QueueSuite) TestDestroyQueue(c *C) {
    queue := s.store.Queues.Create(model.A{"Name": "Destroy me"}, s.owner)
    c.Assert(queue.Destroy().Uid(), Equals, queue.Uid())
    c.Assert(s.store.Queues.Find(queue.Uid()), IsNil)
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

func (s *QueueSuite) TestCountAndEachQueues(c *C) {
    s.store.Queues.DestroyAll()
    c.Assert(s.store.Queues.Count(), Equals, 0)

    s.store.Queues.Create(model.A{"Name": "One"})
    s.store.Queues.Create(model.A{"Name": "Two"})

    c.Assert(s.store.Queues.Count(), Equals, 2)
    names := make([]string, 0)
    s.store.Queues.Each(func (item interface{}) {
        queue := item.(*model.Queue)
        names = append(names, queue.Name())
    })
    c.Assert(names, DeepEquals, []string{"One", "Two"})
}

func (s *QueueSuite) TestPersistAddTask(c *C) {
    queue := model.NewQueue(model.A{})
    task1 := model.NewTask(model.A{"Title": "One", "Created": int64(41)})
    task2 := model.NewTask(model.A{"Title": "Two", "Created": int64(42)})
    c.Assert(queue.NumberTasks(), Equals, 0)
    queue.PersistAddTask(task1)
    c.Assert(task1.Weight(), Equals, int64(41))
    queue.PersistAddTask(task2)
    c.Assert(task2.Weight(), Equals, int64(42))
    c.Assert(queue.NumberTasks(), Equals, 2)
    c.Assert(queue.NextTaskUid(), Equals, task1.Uid())
}

func (s *QueueSuite) TestAddTask(c *C) {
    caller := s.store.Queues.Create(model.A{"Name": "Caller"})
    task1 := s.store.Tasks.Create(model.A{"Title": "One"})
    task2 := s.store.Tasks.Create(model.A{"Title": "Two"})
    task3 := s.store.Tasks.Create(model.A{"Title": "Three"})
    c.Assert(caller.AddTask(task1.Uid()), Equals, true)
    c.Assert(caller.NextTaskUid(), Equals, task1.Uid())
    c.Assert(caller.NumberTasks(), Equals, 1)
    caller.AddTask(task2.Uid())
    caller.AddTask(task3.Uid())
    c.Assert(caller.NextTaskUid(), Equals, task1.Uid())
    c.Assert(caller.NumberTasks(), Equals, 3)
    persisted := s.store.Queues.Find(caller.Uid())
    c.Assert(persisted.NextTaskUid(), Equals, task1.Uid())
    c.Assert(persisted.NumberTasks(), Equals, 3)
}

func (s *QueueSuite) TestPersistDelTask(c *C) {
    queue := model.NewQueue(model.A{})
    task1 := model.NewTask(model.A{"Title": "One", "Created": int64(41)})
    task2 := model.NewTask(model.A{"Title": "Two", "Created": int64(42)})
    queue.PersistAddTask(task1)
    queue.PersistAddTask(task2)
    queue.PersistDelTask(task1)
    c.Assert(queue.NumberTasks(), Equals, 1)
    c.Assert(queue.NextTaskUid(), Equals, task2.Uid())
    queue.PersistDelTask(task2)
    c.Assert(queue.NumberTasks(), Equals, 0)
    c.Assert(queue.NextTaskUid(), Equals, "")
}

func (s *QueueSuite) TestDelTask(c *C) {
    caller := s.store.Queues.Create(model.A{"Name": "Caller"})
    task1 := s.store.Tasks.Create(model.A{"Title": "One", "Created": int64(1)})
    task2 := s.store.Tasks.Create(model.A{"Title": "Two", "Created": int64(2)})
    task3 := s.store.Tasks.Create(model.A{"Title": "Three", "Created": int64(3)})
    caller.AddTask(task1.Uid())
    caller.AddTask(task2.Uid())
    caller.AddTask(task3.Uid())
    caller.DelTask(task1.Uid())
    c.Assert(caller.NextTaskUid(), Equals, task2.Uid())
    c.Assert(caller.NumberTasks(), Equals, 2)
    persisted := s.store.Queues.Find(caller.Uid())
    c.Assert(persisted.NextTaskUid(), Equals, task2.Uid())
    c.Assert(persisted.NumberTasks(), Equals, 2)
}
