package acceptance_test

import (
    . "launchpad.net/gocheck"
    "testing"
    "fmt"
    "goatd/app/event"
    "goatd/app/model"
    "goatd/app/dispatch"
)

func Test(t *testing.T) { TestingT(t) }

type AcceptanceSuite struct {
    busManager *event.BusManager
    store *model.Store
}

var _ = Suite(&AcceptanceSuite{})

func (s *AcceptanceSuite) SetUpTest(c *C) {
    s.busManager = event.NewBusManager()
    s.busManager.Start()
    s.store = model.NewStore(s.busManager)
}

func (s *AcceptanceSuite) TearDownTest(c *C) {
    s.busManager.Stop()
}

func (s *AcceptanceSuite) TestAssignsATaskToATeamMate(c *C) {
    // provisioning
    team := s.store.Teams.Create(model.A{"Name": "Jones Household"})
    manager := dispatch.NewManager(s.store)
    teammate := s.store.Teammates.Create(model.A{"Name": "Jack"}, team)
    queue := s.store.Queues.Create(model.A{"Name": "Duties"}, team)
    s.store.Skills.Create(model.A{"TeammateUid": teammate.Uid(),
        "QueueUid": queue.Uid()}, team)

    // keep track of events received
    var eventOne, eventTwo, eventThree event.Event
    go func() {
        incoming := s.busManager.SubscribeTo([]event.Kind{
            event.OfferTask, event.AcceptTask, event.CompleteTask,
        })
        eventOne = <-incoming; fmt.Println("eventOne:", eventOne)
        eventTwo = <-incoming; fmt.Println("eventTwo:", eventTwo)
        eventThree = <-incoming; fmt.Println("eventThree:", eventThree)
    }()

    // create task and get have the teammate process it
    c.Assert(teammate.Status(), Equals, model.StatusSignedOut)
    teammate.SignIn()
    c.Assert(teammate.Status(), Equals, model.StatusOnBreak)

    task := s.store.Tasks.Create(model.A{"Title": "take out the trash"}, team)
    c.Assert(task.Status(), Equals, model.StatusCreated)
    manager.QueueTask(queue, task)
    c.Assert(task.Status(), Equals, model.StatusQueued)

    manager.MakeTeammateAvailable(teammate)
    task = task.Reload()
    c.Assert(teammate.Status(), Equals, model.StatusOffered)
    c.Assert(teammate.CurrentTask().Uid(), DeepEquals, task.Uid())
    c.Assert(task.Status(), Equals, model.StatusOffered)

    c.Assert(eventOne.Kind, Equals, event.OfferTask)
    c.Assert(eventOne.Data[0], Equals, teammate.Uid())
    c.Assert(eventOne.Data[1], Equals, task.Uid())

    manager.AcceptTask(teammate, task)
    c.Assert(teammate.Status(), Equals, model.StatusBusy)
    c.Assert(teammate.CurrentTask(), DeepEquals, task)
    c.Assert(task.Status(), Equals, model.StatusAssigned)
    c.Assert(queue.NextTaskUid(), DeepEquals, task.Uid())

    c.Assert(eventTwo.Kind, Equals, event.AcceptTask)
    c.Assert(eventTwo.Data[0], Equals, teammate.Uid())
    c.Assert(eventTwo.Data[1], Equals, task.Uid())

    manager.FinishTask(teammate, task)
    c.Assert(teammate.Status(), Equals, model.StatusWrappingUp)
    c.Assert(teammate.CurrentTask(), IsNil)
    c.Assert(task.Status(), Equals, model.StatusCompleted)
    c.Assert(queue.NextTaskUid(), DeepEquals, "")

    c.Assert(eventThree.Kind, Equals, event.CompleteTask)
    c.Assert(eventThree.Data[0], Equals, teammate.Uid())
    c.Assert(eventThree.Data[1], Equals, task.Uid())

    teammate.StartOtherWork()
    c.Assert(teammate.Status(), Equals, model.StatusOtherWork)

    teammate.GoOnBreak()
    c.Assert(teammate.Status(), Equals, model.StatusOnBreak)

    teammate.SignOut()
    c.Assert(teammate.Status(), Equals, model.StatusSignedOut)
}
