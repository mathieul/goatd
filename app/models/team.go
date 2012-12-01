package models

type Team struct {
	Storage
	Teammates *Teammates
	Queues *Queues
	Skills *Skills
    AttrName string
}

func NewTeam(attributes Attrs) (team *Team) {
	team = newModel(&Team{}, &attributes).(*Team)
	team.Teammates = NewTeammates("Team", team.Uid())
	team.Queues = NewQueues("Team", team.Uid())
	team.Skills = NewSkills("Team", team.Uid())
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
