package models

/*
 * Team
 */

type Team struct {
    Storage
    Teammates *Teammates
    Queues *Queues
    Skills *Skills
    AttrName string
}

func NewTeam(attributes Attrs) (team *Team) {
    team = newModel(&Team{}, &attributes).(*Team)
    uid := team.Uid()
    team.Teammates = NewTeammates("Team", uid)
    team.Queues = NewQueues("Team", uid)
    team.Skills = NewSkills("Team", uid)
    return team
}

func CreateTeam(attributes Attrs) (team *Team) {
    team = NewTeam(attributes)
    team.Save()
    return team
}

func (team *Team) Name() string {
    return team.AttrName
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
