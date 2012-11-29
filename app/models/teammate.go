package models

type Teammate struct {
	Storage
    AttrName string
    AttrTeamUid string
}

func NewTeammate(attributes Attrs) *Teammate {
	return newModel(&Teammate{}, &attributes).(*Teammate)
}

func CreateTeammate(attributes Attrs) (teammate *Teammate) {
	teammate = NewTeammate(attributes)
	teammate.Save()
	return teammate
}

func (team *Teammate) Name() string {
    return team.AttrName
}

func (team *Teammate) TeamUid() string {
    return team.AttrTeamUid
}
