package model

import (
    "goatd/app/event"
)

/*
 * Team
 */

type Team struct {
    *event.Identity
    // *Teammates
    // *Queues
    // *Skills
    // *Tasks
    AttrName string
}

func (team Team) Name() string {
    return team.AttrName
}

func NewTeam(attributes A) (team *Team) {
    team = newModel(&Team{}, &attributes).(*Team)
    team.Identity = event.NewIdentity("Team")
    // team.Teammates = NewTeammates(*team.identity)
    // team.Queues = NewQueues(*team.identity)
    // team.Skills = NewSkills(*team.identity)
    // team.Tasks = NewTasks(*team.identity)
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

func NewTeams() (teams *Teams) {
    teams = new(Teams)
    teams.Collection = NewCollection(func(attributes A) interface{} {
        team := NewTeam(attributes)
        return team
    }, nil)
    return teams
}

func (teams *Teams) New(attributes A) (team *Team) {
    return teams.Collection.New(attributes).(*Team)
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
