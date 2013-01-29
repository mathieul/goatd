package dispatch_test

import (
    . "launchpad.net/gocheck"
    "goatd/app/event"
    "goatd/app/model"
    "goatd/app/dispatch"
)

type DistributorSuite struct {
    busManager *event.BusManager
    store      *model.Store
    team       *model.Team
}

var _ = Suite(&DistributorSuite{})

func (s *DistributorSuite) SetUpTest(c *C) {
    s.busManager = event.NewBusManager()
    s.busManager.Start()
    s.store = model.NewStore(s.busManager)
    s.team = s.store.Teams.Create(model.A{"Name": "James Bond"})
}

func (s *DistributorSuite) TestFindTaskForTeammateWhenOne(c *C) {
    task := s.store.Tasks.Create(model.A{"Title": "Skyfall"}, s.team)
    queue := s.store.Queues.Create(model.A{"Name": "Daniel Craig"}, s.team)
    teammate := s.store.Teammates.Create(model.A{"Name": "007"}, s.team)
    s.store.Skills.Create(model.A{"Level": model.LevelHigh}, teammate, queue)
    manager := dispatch.NewManager(s.store)
    manager.QueueTask(queue, task)
    task2 := s.store.Tasks.Find(task.Uid())
    queue2 := s.store.Queues.Find(queue.Uid())
    c.Assert(task2.Status(), Equals, model.StatusQueued)
    c.Assert(task2.QueueUid(), Equals, queue.Uid())
    c.Assert(queue2.QueuedTaskUids(), DeepEquals, []string{task.Uid()})

    found := dispatch.FindTaskForTeammate(s.store, teammate)
    c.Assert(found, Not(IsNil))
    c.Assert(found.Uid(), Equals, task.Uid())
}
