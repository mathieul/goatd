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

type Teammates struct {
	items []*Teammate
	ownerUid string
	ownerName string
}

func (teammates *Teammates) Create(attributes Attrs) (teammate *Teammate) {
	attributes[teammates.ownerName + "Uid"] = teammates.ownerUid
	teammate = CreateTeammate(attributes)
	teammates.items = append(teammates.items, teammate)
	return teammate
}

func (teammates *Teammates) SetOwner(name, uid string) {
	teammates.ownerName = name
	teammates.ownerUid = uid
}
