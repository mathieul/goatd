package dispatch

import (
    "goatd/app/model"
)


/*
 * FindTaskForTeammate
 */
func FindTaskForTeammate(store *model.Store, teammate *model.Teammate) *model.Task {
    teamUid := teammate.TeamUid()
    queuesWithTasks := store.Queues.Select(func (item interface{}) bool {
        queue := item.(*model.Queue)
        // TODO: select only queues with a task waiting
        return queue.TeamUid() == teamUid && queue.NextTaskUid() != ""
    })
    if len(queuesWithTasks) == 0 { return nil }
    taskUid := queuesWithTasks[0].NextTaskUid()
    task := store.Tasks.Find(taskUid)
    return task
}