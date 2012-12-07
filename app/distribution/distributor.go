package distribution

import (
    "goatd/app/models"
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

func NewDistributor(team *models.Team) *Distributor {
    return &Distributor{team, make(map[Event][]CallbackFunc)}
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
