package models

type Team struct {
	Storage
	Teammates
    AttrName string
}

func NewTeam(attributes Attrs) (team *Team) {
	team = newModel(&Team{}, &attributes).(*Team)
	team.Teammates.SetOwner("Team", team.Uid())
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
