package models

import (
    "goatd/app.old/identification"
)

/*
 * Team
 */

type Team struct {
    Storage
    identity *identification.Identity
    *Teammates
    *Queues
    *Skills
    *Tasks
    AttrName string
}

func (team Team) Uid() string {
    return team.Storage.Uid()
}

func (team Team) Name() string {
    return team.AttrName
}

func (team Team) TasksQueued(queue *Queue) []*Task {
    queueUid := queue.Uid()
    return team.Tasks.Select(func (item interface{}) bool {
        task := item.(*Task)
        return task.Status() == StatusQueued && task.QueueUid() == queueUid
    })
}

func (team Team) Identity() identification.Identity {
    return *team.identity
}

func NewTeam(attributes Attrs) (team *Team) {
    team = newModel(&Team{}, &attributes).(*Team)
    team.identity = identification.NewIdentity("Team", team.Uid(), team)
    team.Teammates = NewTeammates(*team.identity)
    team.Queues = NewQueues(*team.identity)
    team.Skills = NewSkills(*team.identity)
    team.Tasks = NewTasks(*team.identity)
    return team
}

func CreateTeam(attributes Attrs) (team *Team) {
    team = NewTeam(attributes)
    team.Save()
    return team
}


/*
 * Teams
 */

type Teams struct {
    Collection
}

func toTeamSlice(source []interface{}) []*Team {
    teams := make([]*Team, 0, len(source))
    for _, team := range source {
        teams = append(teams, team.(*Team))
    }
    return teams
}

func NewTeams(owner identification.Identity) (teams *Teams) {
    teams = new(Teams)
    teams.Collection = NewCollection(func(attributes Attrs, parent interface{}) interface{} {
        team := CreateTeam(attributes)
        return team
    }, owner)
    return teams
}

func (teams *Teams) Create(attributes Attrs) (team *Team) {
    return teams.Collection.Create(attributes).(*Team)
}

func (teams Teams) Slice() []*Team {
    return toTeamSlice(teams.Collection.Slice())
}

func (teams Teams) Find(uid string) *Team {
    if found := teams.Collection.Find(uid); found != nil {
        return found.(*Team)
    }
    return nil
}

func (teams Teams) FindAll(uids []string) []*Team {
    return toTeamSlice(teams.Collection.FindAll(uids))
}

func (teams Teams) Select(tester func(interface{}) bool) (result []*Team) {
    return toTeamSlice(teams.Collection.Select(tester))
}
