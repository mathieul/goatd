package models

type Attrs []string

type Team struct {
    uid string
    name string
}

func CreateTeam(attributes Attrs) (team *Team) {
    team = new(Team)
    team.uid = "TODO"
    team.name = "Metallica"
    return team
}

func (team *Team) Uid() string {
    return team.uid
}

func (team *Team) Name() string {
    return team.name
}
