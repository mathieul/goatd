package distribution

import (
    "goatd/app/models"
    "goatd/app/event"
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
        incoming := event.Manager().SubscribeTo([]event.Kind{event.KindTeammateAvailable})
        for theEvent := range incoming {
            teammate := theEvent.Identity.Value().(*models.Teammate)
            distributor.FindAndAssignTaskForTeammate(teammate)
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
    task := TaskSelectorByOldestNextTask(queues)
    teammate.OfferTask(task)
}

func TaskSelectorByOldestNextTask(queues []*models.Queue) *models.Task {
    // TODO, for now just return the first one :D
    return queues[0].QueuedTasks()[0]
}