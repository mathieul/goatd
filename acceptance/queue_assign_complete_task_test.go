package acceptance_test

import (
    . "launchpad.net/gocheck"
    "testing"
    "fmt"
    "time"
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
    aLittleBit := 100 * time.Millisecond

    // provisioning
    team := s.store.Teams.Create(model.A{"Name": "Jones Household"})
    distributor := dispatch.NewDistributor(s.store)
    teammate := s.store.Teammates.Create(model.A{"Name": "Jack"}, team)
    queue := s.store.Queues.Create(model.A{"Name": "Duties"}, team)
    skill := s.store.Skills.Create(model.A{"TeammateUid": teammate.Uid(),
        "QueueUid": queue.Uid()}, team)

    // keep track of events received
    var eventOne, eventTwo, eventThree event.Event
    go func() {
        incoming := event.Manager().SubscribeTo([]event.Kind{
            event.EventOfferTask, event.EventAcceptTask, event.EventCompleteTask,
        })
        eventOne = <-incoming; fmt.Println("eventOne:", eventOne)
        eventTwo = <-incoming; fmt.Println("eventTwo:", eventTwo)
        eventThree = <-incoming; fmt.Println("eventThree:", eventThree)
    }()

    // create task and get have the teammate process it
    c.Assert(teammate.Status(), Equals, model.StatusSignedOut)
    teammate.SignIn()
    c.Assert(teammate.Status(), Equals, model.StatusOnBreak)

    task := team.Tasks.Create(model.A{"Title": "take out the trash"})
    c.Assert(task.Status(), Equals, model.StatusCreated)
    task.Enqueue(queue)
    c.Assert(task.Status(), Equals, model.StatusQueued)

    teammate.MakeAvailable()
    time.Sleep(aLittleBit)
    c.Assert(teammate.Status(), Equals, model.StatusOffered)
    c.Assert(teammate.CurrentTask(), DeepEquals, task)
    c.Assert(task.Status(), Equals, model.StatusOffered)

    c.Assert(eventOne.Kind, Equals, event.EventOfferTask)
    c.Assert(eventOne.Data[0], Equals, teammate.Uid())
    c.Assert(eventOne.Data[1], Equals, task.Uid())

    teammate.AcceptTask(task)
    time.Sleep(aLittleBit)
    c.Assert(teammate.Status(), Equals, model.StatusBusy)
    c.Assert(teammate.CurrentTask(), DeepEquals, task)
    c.Assert(task.Status(), Equals, model.StatusAssigned)
    c.Assert(queue.QueuedTasks(), DeepEquals, []*model.Task{task})

    c.Assert(eventTwo.Kind, Equals, event.EventAcceptTask)
    c.Assert(eventTwo.Data[0], Equals, teammate.Uid())
    c.Assert(eventTwo.Data[1], Equals, task.Uid())

    teammate.FinishTask(task)
    time.Sleep(aLittleBit)
    c.Assert(model.StatusWrappingUp, Equals, teammate.Status())
    c.Assert(teammate.CurrentTask(), IsNil)
    c.Assert(task.Status(), Equals, model.StatusCompleted)
    c.Assert(queue.QueuedTasks(), DeepEquals, []*model.Task{})

    c.Assert(eventThree.Kind, Equals, event.EventCompleteTask)
    c.Assert(eventThree.Data[0], Equals, teammate.Uid())
    c.Assert(eventThree.Data[1], Equals, task.Uid())

    teammate.StartOtherWork()
    c.Assert(teammate.Status(), Equals, model.StatusOtherWork)

    teammate.GoOnBreak()
    c.Assert(teammate.Status(), Equals, model.StatusOnBreak)

    teammate.SignOut()
    c.Assert(teammate.Status(), Equals, model.StatusSignedOut)
}
