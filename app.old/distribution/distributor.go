package distribution

import (
    "goatd/app.old/models"
    "goatd/app.old/event"
)


/*
 * Basic types
 */
type Event int
type CallbackFunc func (Event, []interface{})


/*
 * Distributor
 */
type Distributor struct {
    team *models.Team
    callbacks map[Event][]CallbackFunc
    tracker *ResourceTracker
}

func NewDistributor(team *models.Team) (distributor *Distributor) {
    if !event.Manager().Running() {
        panic("Event manager must be running to instantiate a distributor.")
    }
    distributor = new(Distributor)
    distributor.team = team
    distributor.callbacks = make(map[Event][]CallbackFunc)
    distributor.tracker = NewResourceTracker(team)
    distributor.monitorDistributionTriggers()
    return distributor
}

func (distributor *Distributor) monitorDistributionTriggers() {
    go func() {
        incoming := event.Manager().SubscribeTo([]event.Kind{
            event.KindTeammateAvailable,
            event.KindAcceptTask,
            event.KindCompleteTask,
        })
        for theEvent := range incoming {
            teammate := theEvent.Identity.Value().(*models.Teammate)
            switch theEvent.Kind {
            case event.KindTeammateAvailable:
                distributor.FindAndAssignTaskForTeammate(teammate)
            case event.KindAcceptTask:
                distributor.AssignTask(teammate, theEvent.Data[1])
            case event.KindCompleteTask:
                distributor.CompleteTask(teammate, theEvent.Data[1])
            }
        }
    }()
}

func (distributor Distributor) Team() *models.Team {
    return distributor.team
}

func (distributor Distributor) AddTeammateToQueue(queue *models.Queue,
        teammate *models.Teammate, level int) bool {
    skills := distributor.team.Skills
    attributes := models.Attrs{"TeammateUid": teammate.Uid(),
        "QueueUid": queue.Uid(), "Level": level, "Enabled": true}
    if skill := skills.Create(attributes); skill == nil { return false }
    return true
}

func (distributor *Distributor) FindAndAssignTaskForTeammate(teammate *models.Teammate) {
    queues := distributor.tracker.TeammateQueuesReady(teammate)
    if task := TaskSelectorByOldestNextTask(queues); task != nil {
        if teammate.OfferTask(task) {
            if task.Offer() {
                event.Manager().PublishEvent(event.KindOfferTask, distributor.team.Identity(),
                    []string{teammate.Uid(), task.Uid()})
            } else {
                teammate.RejectTask(task)
            }
        }
    }
}

func (distributor *Distributor) AssignTask(teammate *models.Teammate, taskUid string) {
    task := distributor.team.Tasks.Find(taskUid)
    task.Assign()
}

func (distributor *Distributor) CompleteTask(teammate *models.Teammate, taskUid string) {
    task := distributor.team.Tasks.Find(taskUid)
    task.Complete()
    queue := distributor.team.Queues.Find(task.QueueUid())
    queue.RemoveTask(task)
}

func TaskSelectorByOldestNextTask(queues []*models.Queue) *models.Task {
    // TODO, for now just return the first one :D
    if len(queues) == 0 { return nil }
    tasks := queues[0].QueuedTasks()
    if len(tasks) == 0 { return nil }
    return tasks[0]
}