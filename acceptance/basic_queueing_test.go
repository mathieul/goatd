package acceptance_test

import (
    . "launchpad.net/gocheck"
    "testing"
    "fmt"
    "time"
    "goatd/app/identification"
    "goatd/app/event"
    "goatd/app/models"
    "goatd/app/distribution"
)

func Test(t *testing.T) { TestingT(t) }

type AcceptanceS struct{
    teams  *models.Teams
    team  *models.Team
    mate  *models.Teammate
    queue *models.Queue
}

var _ = Suite(&AcceptanceS{})

func (s *AcceptanceS) SetUpTest(c *C) {
    event.Manager().Start()
    s.teams = models.NewTeams(identification.NoIdentity())
    s.team = s.teams.Create(models.Attrs{"Name": "Wedding"})
    s.mate = s.team.Teammates.Create(models.Attrs{"Name": "Bride"})
    s.queue = s.team.Queues.Create(models.Attrs{"Name": "Thank you notes"})
}

func (s *AcceptanceS) TearDownTest(c *C) {
    event.Manager().Stop()
}


func (s *AcceptanceS) TestAssignsATaskToATeamMate(c *C) {
    aLittleBit := 100 * time.Millisecond
    distributor := distribution.NewDistributor(s.team)

    var eventOne, eventTwo, eventThree event.Event
    go func() {
        incoming := event.Manager().SubscribeTo([]event.Kind{
            event.OfferTask, event.AssignTask, event.CompleteTask,
        })
        eventOne = <-incoming; fmt.Println("eventOne:", eventOne)
        eventTwo = <-incoming; fmt.Println("eventTwo:", eventTwo)
        eventThree = <-incoming; fmt.Println("eventThree:", eventThree)
    }()

    distributor.AddTeammateToQueue(s.queue, s.mate, models.LevelHigh)
    skills := s.team.Skills.Select(func (item interface{}) bool {
        skill := item.(*models.Skill)
        return skill.TeammateUid() == s.mate.Uid() && skill.QueueUid() == s.queue.Uid()
    })
    c.Assert(len(skills), Equals, 1)
    c.Assert(skills[0].Level(), Equals, models.LevelHigh)

    c.Assert(s.mate.Status(), Equals, models.StatusSignedOut)
    s.mate.SignIn()
    c.Assert(s.mate.Status(), Equals, models.StatusOnBreak)

    task := models.NewTask(models.Attrs{"Title": "thank Jones family"})
    c.Assert(task.Status(), Equals, models.StatusCreated)
    task.Enqueue(s.queue)
    c.Assert(task.Status(), Equals, models.StatusQueued)

    s.mate.MakeAvailable()
    time.Sleep(aLittleBit)
    c.Assert(s.mate.Status(), Equals, models.StatusOffered)
    c.Assert(s.mate.CurrentTask(), DeepEquals, task)
    c.Assert(task.Status(), Equals, models.StatusOffered)

    c.Assert(eventOne.Kind, Equals, event.OfferTask)
    c.Assert(eventOne.Data[0], Equals, s.queue.Uid())
    c.Assert(eventOne.Data[1], Equals, s.mate.Uid())
    c.Assert(eventOne.Data[2], Equals, task.Uid())

    s.mate.AcceptTask(task)
    time.Sleep(aLittleBit)
    c.Assert(models.StatusBusy, Equals, s.mate.Status())
    c.Assert(task, DeepEquals, s.mate.CurrentTask())
    c.Assert(models.StatusAssigned, Equals, task.Status())
    c.Assert(s.queue.QueuedTasks(), DeepEquals, []*models.Task{task})

    c.Assert(eventTwo.Kind, Equals, event.AssignTask)
    c.Assert(eventTwo.Data[0], Equals, s.queue.Uid())
    c.Assert(eventTwo.Data[1], Equals, s.mate.Uid())
    c.Assert(eventTwo.Data[2], Equals, task.Uid())

    s.mate.FinishTask(task)
    time.Sleep(aLittleBit)
    c.Assert(models.StatusWrappingUp, Equals, s.mate.Status())
    c.Assert(s.mate.CurrentTask(), IsNil)
    c.Assert(models.StatusCompleted, Equals, task.Status())
    c.Assert([]models.Task{}, DeepEquals, s.queue.Tasks().Slice())

    s.mate.StartOtherWork()
    c.Assert(models.StatusOtherWork, Equals, s.mate.Status())

    s.mate.GoOnBreak()
    c.Assert(models.StatusOnBreak, Equals, s.mate.Status())

    s.mate.SignOut()
    c.Assert(models.StatusSignedOut, Equals, s.mate.Status())
}
