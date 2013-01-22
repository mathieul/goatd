package acceptance_test

import (
    . "launchpad.net/gocheck"
    "testing"
    "fmt"
    "time"
    "goatd/app/event"
    "goatd/app/model"
)

func Test(t *testing.T) { TestingT(t) }

type AcceptanceSuite struct{
    busManager *event.BusManager
    store *model.Store
}

var _ = Suite(&AcceptanceSuite{})

func (s *AcceptanceSuite) SetUpTest(c *C) {
    s.busManager = event.NewBusManager()
    s.busManager.Start()
    s.store = model.NewStore(s.busManager)
    s.store.Start()
}

func (s *AcceptanceSuite) TearDownTest(c *C) {
    s.store.Stop()
    s.busManager.Stop()
}

func (s *AcceptanceSuite) TestAssignsATaskToATeamMate(c *C) {
    aLittleBit := 100 * time.Millisecond

    // provisioning
    team := s.store.Teams.Create(model.A{"Name": "Jones Household"})
    distributor := dispatch.NewDistributor(team, s.busManager, s.store)
    mate := s.store.Teammates.Create(model.A{"Name": "Jack"}, team)
    queue := s.store.Queues.Create(model.A{"Name": "Duties"}, team)
    skill := s.store.Skills.Create(model.A{}, team, mate, queue)

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
    c.Assert(mate.Status(), Equals, model.StatusSignedOut)
    mate.SignIn()
    c.Assert(mate.Status(), Equals, model.StatusOnBreak)

    task := team.Tasks.Create(model.A{"Title": "take out the trash"})
    c.Assert(task.Status(), Equals, model.StatusCreated)
    task.Enqueue(queue)
    c.Assert(task.Status(), Equals, model.StatusQueued)

    mate.MakeAvailable()
    time.Sleep(aLittleBit)
    c.Assert(mate.Status(), Equals, model.StatusOffered)
    c.Assert(mate.CurrentTask(), DeepEquals, task)
    c.Assert(task.Status(), Equals, model.StatusOffered)

    c.Assert(eventOne.Kind, Equals, event.EventOfferTask)
    c.Assert(eventOne.Data[0], Equals, mate.Uid())
    c.Assert(eventOne.Data[1], Equals, task.Uid())

    mate.AcceptTask(task)
    time.Sleep(aLittleBit)
    c.Assert(mate.Status(), Equals, model.StatusBusy)
    c.Assert(mate.CurrentTask(), DeepEquals, task)
    c.Assert(task.Status(), Equals, model.StatusAssigned)
    c.Assert(queue.QueuedTasks(), DeepEquals, []*model.Task{task})

    c.Assert(eventTwo.Kind, Equals, event.EventAcceptTask)
    c.Assert(eventTwo.Data[0], Equals, mate.Uid())
    c.Assert(eventTwo.Data[1], Equals, task.Uid())

    mate.FinishTask(task)
    time.Sleep(aLittleBit)
    c.Assert(model.StatusWrappingUp, Equals, mate.Status())
    c.Assert(mate.CurrentTask(), IsNil)
    c.Assert(task.Status(), Equals, model.StatusCompleted)
    c.Assert(queue.QueuedTasks(), DeepEquals, []*model.Task{})

    c.Assert(eventThree.Kind, Equals, event.EventCompleteTask)
    c.Assert(eventThree.Data[0], Equals, mate.Uid())
    c.Assert(eventThree.Data[1], Equals, task.Uid())

    mate.StartOtherWork()
    c.Assert(mate.Status(), Equals, model.StatusOtherWork)

    mate.GoOnBreak()
    c.Assert(mate.Status(), Equals, model.StatusOnBreak)

    mate.SignOut()
    c.Assert(mate.Status(), Equals, model.StatusSignedOut)
}
