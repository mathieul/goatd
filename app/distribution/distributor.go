package distribution

import (
    "goatd/app/models"
)

const (
    EventOfferTask Event = 1
    EventAssignTask
    EventCompleteTask
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
