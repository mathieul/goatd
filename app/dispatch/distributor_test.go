package dispatch_test

import (
    . "launchpad.net/gocheck"
    "goatd/app/model"
    "goatd/app/dispatch"
)

type DistributorSuite struct {
    store *model.Store
    distributor *dispatch.Distributor
    team *model.Team
}

var _ = Suite(&DistributorSuite{})

func (s *DistributorSuite) SetUpTest(c *C) {
    s.store = model.NewStore(nil)
    s.distributor = dispatch.NewDistributor(s.store)
    s.team = s.store.Teams.Create(model.A{"Name": "James Bond"})
}

func (s *DistributorSuite) TestCreateDistributor(c *C) {
    distributor := dispatch.NewDistributor(s.store)
    c.Assert(distributor.Store, Equals, s.store)
}

func (s *DistributorSuite) TestQueueTask(c *C) {
    queue := s.store.Queues.Create(model.A{"Name": "Daniel Craig"}, s.team)
    task := s.store.Tasks.Create(model.A{"Title": "Skyfall"}, s.team)
    c.Assert(s.distributor.QueueTask(queue, task), Equals, true)
    c.Assert(task.Status(), Equals, model.StatusQueued)
    c.Assert(task.QueueUid(), Equals, queue.Uid())
    c.Assert(queue.QueuedTaskUids(), DeepEquals, []string{task.Uid()})
    c.Assert(s.distributor.QueueTask(queue, task), Equals, false)
}
