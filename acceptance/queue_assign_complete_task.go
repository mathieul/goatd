package acceptance_test

import (
    . "launchpad.net/gocheck"
    "testing"
    "fmt"
    "time"
    "goatd/app/event"
    "goatd/app/model"
    "goatd/app/store"
)

func Test(t *testing.T) { TestingT(t) }

type AcceptanceSuite struct{
    busManager *event.BusManager
    store *store.MemoryStore
}

var _ = Suite(&AcceptanceSuite{})

func (s *AcceptanceSuite) SetUpTest(c *C) {
    s.busManager = event.NewBusManager()
    s.busManager.Start()
    s.store.Start()
}

func (s *AcceptanceSuite) TearDownTest(c *C) {
    s.store.Stop()
    s.busManager.Stop()
}


func (s *AcceptanceSuite) TestAssignsATaskToATeamMate(c *C) {
    aLittleBit := 100 * time.Millisecond

    # provisioning
    team := s.store.Teams.Create(model.A{"Name": "Jones Household"})
    distributor := dispatch.NewDistributor(team)
    mate := team.Teammates.Create(model.A{"Name": "Jack"})
    queue := team.Queues.Create(model.A{"Name": "Duties"})

    # keep track of events received
    var eventOne, eventTwo, eventThree event.Event
    go func() {
        incoming := event.Manager().SubscribeTo([]event.Kind{
            event.KindOfferTask, event.KindAcceptTask, event.KindCompleteTask,
        })
        eventOne = <-incoming; fmt.Println("eventOne:", eventOne)
        eventTwo = <-incoming; fmt.Println("eventTwo:", eventTwo)
        eventThree = <-incoming; fmt.Println("eventThree:", eventThree)
    }()

    distributor.AddTeammateToQueue(queue, mate, model.LevelHigh)
    skills := team.Skills.Select(func (item interface{}) bool {
        skill := item.(*model.Skill)
        return skill.TeammateUid() == mate.Uid() && skill.QueueUid() == queue.Uid()
    })
    c.Assert(len(skills), Equals, 1)
    c.Assert(skills[0].Level(), Equals, model.LevelHigh)

    c.Assert(mate.Status(), Equals, model.StatusSignedOut)
    mate.SignIn()
    c.Assert(mate.Status(), Equals, model.StatusOnBreak)

    task := team.Tasks.Create(model.Attrs{"Title": "thank Jones family"})
    c.Assert(task.Status(), Equals, model.StatusCreated)
    task.Enqueue(queue)
    c.Assert(task.Status(), Equals, model.StatusQueued)

    mate.MakeAvailable()
    time.Sleep(aLittleBit)
    c.Assert(mate.Status(), Equals, model.StatusOffered)
    c.Assert(mate.CurrentTask(), DeepEquals, task)
    c.Assert(task.Status(), Equals, model.StatusOffered)

    c.Assert(eventOne.Kind, Equals, event.KindOfferTask)
    c.Assert(eventOne.Data[0], Equals, mate.Uid())
    c.Assert(eventOne.Data[1], Equals, task.Uid())

    mate.AcceptTask(task)
    time.Sleep(aLittleBit)
    c.Assert(mate.Status(), Equals, model.StatusBusy)
    c.Assert(mate.CurrentTask(), DeepEquals, task)
    c.Assert(task.Status(), Equals, model.StatusAssigned)
    c.Assert(queue.QueuedTasks(), DeepEquals, []*model.Task{task})

    c.Assert(eventTwo.Kind, Equals, event.KindAcceptTask)
    c.Assert(eventTwo.Data[0], Equals, mate.Uid())
    c.Assert(eventTwo.Data[1], Equals, task.Uid())

    mate.FinishTask(task)
    time.Sleep(aLittleBit)
    c.Assert(model.StatusWrappingUp, Equals, mate.Status())
    c.Assert(mate.CurrentTask(), IsNil)
    c.Assert(task.Status(), Equals, model.StatusCompleted)
    c.Assert(queue.QueuedTasks(), DeepEquals, []*model.Task{})

    c.Assert(eventThree.Kind, Equals, event.KindCompleteTask)
    c.Assert(eventThree.Data[0], Equals, mate.Uid())
    c.Assert(eventThree.Data[1], Equals, task.Uid())

    mate.StartOtherWork()
    c.Assert(mate.Status(), Equals, model.StatusOtherWork)

    mate.GoOnBreak()
    c.Assert(mate.Status(), Equals, model.StatusOnBreak)

    mate.SignOut()
    c.Assert(mate.Status(), Equals, model.StatusSignedOut)
}
