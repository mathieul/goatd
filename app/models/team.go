package models

type Team struct {
	Storage
    AttrName string
}

func NewTeam(attributes Attrs) (team *Team) {
    team = new(Team)
    team.Init()
    for name, value := range attributes {
	    setAttributeValue(team, name, value)
    }
    return team
}

func CreateTeam(attributes Attrs) (team *Team) {
	team = NewTeam(attributes)
	team.Save(true)
	return team
}

func (team *Team) Name() string {
    return team.AttrName
}
