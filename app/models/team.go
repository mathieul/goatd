package models

import (
    "goatd/app/identification"
)

/*
 * Team
 */

type Team struct {
    Storage
    identity *identification.Identity
    Teammates *Teammates
    Queues *Queues
    Skills *Skills
    AttrName string
}

func (team Team) Uid() string {
    return team.Storage.Uid()
}

func (team Team) Name() string {
    return team.AttrName
}

func NewTeam(attributes Attrs) (team *Team) {
    team = newModel(&Team{}, &attributes).(*Team)
    team.identity = identification.NewIdentity("Team", team.Uid(), team)
    team.Teammates = NewTeammates(*team.identity)
    team.Queues = NewQueues(*team.identity)
    team.Skills = NewSkills(*team.identity)
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
    items []*Team
}

func NewTeams() *Teams {
    return new(Teams)
}

func (teams *Teams) Create(attributes Attrs) (team *Team) {
    team = CreateTeam(attributes)
    teams.items = append(teams.items, team)
    return team
}

func (teams Teams) FindAll(uids []string) (found []*Team) {
    for _, candidate := range teams.items {
        candidateUid := candidate.Uid()
        for _, uid := range uids {
            if candidateUid == uid {
                found = append(found, candidate)
            }
        }
    }
    return found
}

func (teams Teams) Find(uid string) *Team {
    found := teams.FindAll([]string{uid})
    if len(found) == 0 {
        return nil
    }
    return found[0]
}
