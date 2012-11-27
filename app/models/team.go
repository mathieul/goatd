package models

type Team struct {
	Storage
    AttrUid string
    AttrName string
}

func NewTeam(attributes Attrs) (team *Team) {
    team = new(Team)
    for name, value := range attributes {
	    setAttributeValue(team, name, value)
    }
    if team.AttrUid == "" {
    	team.AttrUid = generateUid()
    }
    return team
}

func CreateTeam(attributes Attrs) (team *Team) {
	team = NewTeam(attributes)
	team.SetPersisted(true)
	return team
}

func (team *Team) Uid() string {
    return team.AttrUid
}

func (team *Team) Name() string {
    return team.AttrName
}
