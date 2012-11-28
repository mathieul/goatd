package models

type Team struct {
	Storage
    AttrName string
}

func NewTeam(attributes Attrs) *Team {
	return newModel(&Team{}, &attributes).(*Team)
}

func CreateTeam(attributes Attrs) (team *Team) {
	return createModel(&Team{}, &attributes).(*Team)
}

func (team *Team) Name() string {
    return team.AttrName
}
