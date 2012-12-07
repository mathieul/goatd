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
}

func NewDistributor(team *models.Team) (distributor *Distributor) {
    distributor = &Distributor{team, make(map[Event][]CallbackFunc)}
    distributor.setupListeners()
    return distributor
}

func (distributor *Distributor) setupListeners() {
    if !event.Manager().Running() {
        panic("Event manager is not running.")
    }
    go func() {
        incoming := event.Manager().SubscribeTo([]event.Kind{event.KindTeammateAvailable})
        for event.Manager().Running() {
            select {
            case event := <- incoming:
                // TODO
                panic(event)
            }
        }
    }()
}

func (distributor Distributor) Team() *models.Team {
    return distributor.team
}

func (distributor Distributor) On(event Event, callback CallbackFunc) {
    callbackSlice := distributor.callbacks[event]
    distributor.callbacks[event] = append(callbackSlice, callback)
}

func (distributor Distributor) Trigger(event Event, parameters []interface{}) {
    if callbackSlice, found := distributor.callbacks[event]; found {
        for _, callback := range callbackSlice {
            callback(event, parameters)
        }
    }
}

func (distributor Distributor) AddTeammateToQueue(queue *models.Queue,
        teammate *models.Teammate, level int) bool {
    skills := distributor.team.Skills
    attributes := models.Attrs{"TeammateUid": teammate.Uid(),
        "QueueUid": queue.Uid(), "Level": level, "Enabled": true}
    if skill := skills.Create(attributes); skill == nil { return false }
    return true
}
