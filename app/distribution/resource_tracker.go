package distribution

import (
    "fmt"
    "goatd/app/event"
    "goatd/app/models"
)


/*
 * Distributor
 */
type ResourceTracker struct {
    team *models.Team
    queuesByTeammate map[string][]string
}

func (tracker *ResourceTracker) TeammateQueuesReady(teammate *models.Teammate) []*models.Queue {
    uids := tracker.queuesByTeammate[teammate.Uid()]
    fmt.Println("TeammateQueuesReady: uids =", uids)
    allQueues := tracker.team.Queues.FindAll(uids)
    fmt.Println("TeammateQueuesReady: allQueues =", allQueues)
    if allQueues == nil { return []*models.Queue{} }
    queues := make([]*models.Queue, 0, len(allQueues))
    for _, queue := range allQueues {
        if queue.IsReady() {
            queues = append(queues, queue)
        }
    }
    return queues
}

func (tracker *ResourceTracker) startTracking() {
    go func() {
        incoming := event.Manager().SubscribeTo([]event.Kind{event.KindSkillCreated})
        for theEvent := range incoming {
            queueUid, teammateUid := theEvent.Data[0], theEvent.Data[1]
            uids := tracker.queuesByTeammate[teammateUid]
            if uids == nil {
                uids = []string{queueUid}
            } else {
                uids = append(uids, queueUid)
            }
            tracker.queuesByTeammate[teammateUid] = uids
            fmt.Println("startTracking: uids =", uids)
        }
    }()
}

func NewResourceTracker(team *models.Team) (tracker *ResourceTracker) {
    tracker = new(ResourceTracker)
    tracker.team = team
    numberTeams := len(team.Teammates.Slice())
    tracker.queuesByTeammate = make(map[string][]string, numberTeams)
    tracker.startTracking()
    return tracker
}
