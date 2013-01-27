package model_test

import (
    . "launchpad.net/gocheck"
    "strings"
    // "fmt"
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

// func (s *QueueSuite) TestInsertTask(c *C) {
//     caller := s.queues.Create(models.Attrs{"Name": "Caller"})
//     c.Assert(caller.IsReady(), Equals, false)
//     task := s.team.Tasks.Create(models.Attrs{"Title": "Do It"})
//     c.Assert(caller.InsertTask(task), Equals, true)
//     c.Assert(caller.IsReady(), Equals, true)
//     c.Assert(caller.QueuedTasks(), DeepEquals, []*models.Task{task})
// }

// func (s *QueueSuite) TestInsertAndKeepTasksOrdered(c *C) {
//     caller := s.queues.Create(models.Attrs{"Name": "Caller"})
//     tasks := make([]*models.Task, 0, 10)
//     for i := 0; i < 10; i++ {
//         task := s.team.Tasks.Create(models.Attrs{"Title": fmt.Sprintf("#%02d", i)})
//         caller.InsertTask(task)
//         tasks = append(tasks, task)
//     }
//     c.Assert(caller.QueuedTasks(), DeepEquals, tasks)
// }
