package dispatch_test

import (
    . "launchpad.net/gocheck"
    "goatd/app/event"
    "goatd/app/model"
    "goatd/app/dispatch"
)

type ManagerSuite struct {
    store      *model.Store
    busManager *event.BusManager
    manager    *dispatch.Manager
    team       *model.Team
}

var _ = Suite(&ManagerSuite{})

func (s *ManagerSuite) SetUpTest(c *C) {
    s.busManager = event.NewBusManager()
    s.busManager.Start()
    s.store = model.NewStore(s.busManager)
    s.manager = dispatch.NewManager(s.busManager, s.store)
    s.team = s.store.Teams.Create(model.A{"Name": "James Bond"})
}

func (s *ManagerSuite) TearDownTest(c *C) {
    s.busManager.Stop()
}

func (s *ManagerSuite) TestQueueTask(c *C) {
    queue := s.store.Queues.Create(model.A{"Name": "Daniel Craig"}, s.team)
    task := s.store.Tasks.Create(model.A{"Title": "Skyfall"}, s.team)
    c.Assert(s.manager.QueueTask(queue, task), Equals, true)
    c.Assert(task.Status(), Equals, model.StatusQueued)
    c.Assert(task.QueueUid(), Equals, queue.Uid())
    c.Assert(queue.NextTaskUid(), DeepEquals, task.Uid())
    c.Assert(s.manager.QueueTask(queue, task), Equals, false)
}

func (s *ManagerSuite) TestMakeTeammateAvailable(c *C) {
    queue := s.store.Queues.Create(model.A{"Name": "Daniel Craig"}, s.team)
    task := s.store.Tasks.Create(model.A{"Title": "Skyfall"}, s.team)
    s.manager.QueueTask(queue, task)

    teammate := s.store.Teammates.Create(model.A{"Name": "007"}, s.team)
    s.store.Skills.Create(model.A{"Level": model.LevelHigh}, teammate, queue)
    teammate.SignIn()

    c.Assert(s.manager.MakeTeammateAvailable(teammate), Equals, true)
    c.Assert(teammate.Status(), Equals, model.StatusOffered)
    c.Assert(teammate.TaskUid(), DeepEquals, task.Uid())
    c.Assert(task.Reload().Status(), Equals, model.StatusOffered)
}

func (s *ManagerSuite) TestAcceptTask(c *C) {
    queue := s.store.Queues.Create(model.A{"Name": "Daniel Craig"}, s.team)
    task := s.store.Tasks.Create(model.A{"Title": "Skyfall"}, s.team)
    teammate := s.store.Teammates.Create(model.A{"Name": "007"}, s.team)
    s.store.Skills.Create(model.A{"Level": model.LevelHigh}, teammate, queue)
    teammate.SignIn()
    s.manager.QueueTask(queue, task)
    s.manager.MakeTeammateAvailable(teammate)
    task = task.Reload()

    c.Assert(s.manager.AcceptTask(teammate, task), Equals, true)
    c.Assert(teammate.Status(), Equals, model.StatusBusy)
    c.Assert(teammate.TaskUid(), DeepEquals, task.Uid())
    c.Assert(task.Status(), Equals, model.StatusAssigned)
}

func (s *ManagerSuite) TestFinishTask(c *C) {
    queue := s.store.Queues.Create(model.A{"Name": "Daniel Craig"}, s.team)
    task := s.store.Tasks.Create(model.A{"Title": "Skyfall"}, s.team)
    teammate := s.store.Teammates.Create(model.A{"Name": "007"}, s.team)
    s.store.Skills.Create(model.A{"Level": model.LevelHigh}, teammate, queue)
    teammate.SignIn()
    s.manager.QueueTask(queue, task)
    s.manager.MakeTeammateAvailable(teammate)
    task = task.Reload()
    s.manager.AcceptTask(teammate, task)

    c.Assert(s.manager.FinishTask(teammate, task), Equals, true)
    c.Assert(teammate.Status(), Equals, model.StatusWrappingUp)
    c.Assert(teammate.TaskUid(), Equals, "")
    c.Assert(task.Status(), Equals, model.StatusCompleted)
    c.Assert(queue.Reload().NextTaskUid(), DeepEquals, "")
}
