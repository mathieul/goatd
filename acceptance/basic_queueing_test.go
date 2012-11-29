package acceptance_test

import (
	. "launchpad.net/gocheck"
	"testing"
	"goatd/app/models"
	"goatd/app/distribution"
	"time"
)

func Test(t *testing.T) { TestingT(t) }

type AcceptanceS struct{
	team  *models.Team
	mate  *models.Teammate
	queue *models.Queue
}

type State struct {
	event models.Event
	mate models.Teammate
	queue models.Queue
	task models.Task
}

var _ = Suite(&AcceptanceS{})

func (s *AcceptanceS) SetUpTest(c *C) {
	s.team = models.CreateTeam(models.Attrs{"Name": "Wedding"})
	s.mate = s.team.Teammates.Create(models.Attrs{"Name": "Bride"})
	s.queue = s.team.Queues.Create(models.Attrs{"Name": "Thank you notes"})
}

func (s *AcceptanceS) TestAssignsATaskToATeamMate(c *C) {
	aLittleBit := 100 * time.Millisecond
	distributor := distribution.NewDistributor(s.team)

	var state State
	go func() {
		distributor.On(models.EventOfferTask, func (queue, mate, task) {
			state.event = models.EventOfferTask
			state.queue = queue; state.mate = mate; state.task = task
		})
		distributor.On(models.EventAssignTask, func (queue, mate, task) {
			state.event = models.EventAssignTask
			state.queue = queue; state.mate = mate; state.task = task
		})
		distributor.On(models.EventCompleteTask, func (queue, mate, task) {
			state.event = models.EventCompleteTask
			state.queue = queue; state.mate = mate; state.task = task
		})
	}()

	distributor.AddToQueue(s.queue, s.mate, "level", "high", "enabled", true)
	c.Assert([]models.Queue{s.queue}, DeepEquals, s.mate.Queues().Slice())

	c.Assert(models.StatusSignedOut, Equals, s.mate.Status())
	s.mate.signIn()
	c.Assert(models.StatusOnBreak, Equals, s.mate.Status())

	task := models.NewTask(models.Attrs{"Title": "thank Jones family"})
	c.Assert("created", Equals, task.Status())
	distributor.EnqueueTask(s.queue, task, distribution.PriorityMedium)
	c.Assert("queued", Equals, task.Status())
	c.Assert([]models.Task{task}, DeepEquals, s.queue.Tasks().Slice())

	s.mate.makeAvailable()
	time.Sleep(aLittleBit)
	c.Assert(models.StatusOffered, Equals, s.mate.Status())
	c.Assert(task, DeepEquals, s.mate.CurrentTask())
	c.Assert(models.StatusOffered, Equals, task.Status())

	c.Assert(models.EventOfferTask, Equals, state.event)
	c.Assert(s.queue, Equals, state.queue)
	c.Assert(s.mate, Equals, state.mate)
	c.Assert(task, Equals, state.task)

	s.mate.AcceptTask(task)
	time.Sleep(aLittleBit)
	c.Assert(models.StatusBusy, Equals, s.mate.Status())
	c.Assert(task, DeepEquals, s.mate.CurrentTask())
	c.Assert(models.StatusAssigned, Equals, task.Status())
	c.Assert([]models.Task{task}, DeepEquals, s.queue.Tasks().Slice())

	c.Assert(models.EventAssignTask, Equals, state.event)
	c.Assert(s.queue, Equals, state.queue)
	c.Assert(s.mate, Equals, state.mate)
	c.Assert(task, Equals, state.task)

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
