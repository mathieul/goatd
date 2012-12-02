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
