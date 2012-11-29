package models

type Team struct {
	Storage
	Teammates
	Queues
    AttrName string
}

func NewTeam(attributes Attrs) (team *Team) {
	team = newModel(&Team{}, &attributes).(*Team)
	team.Teammates.SetOwner("Team", team.Uid())
	team.Queues.SetOwner("Team", team.Uid())
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

type Teammates struct {
	Owner
	items []*Teammate
}

func (teammates *Teammates) Create(attributes Attrs) (teammate *Teammate) {
	teammate = CreateTeammate(teammates.AddOwnerToAttributes(attributes))
	teammates.items = append(teammates.items, teammate)
	return teammate
}

type Queues struct {
	Owner
	items []*Queue
}

func (queues *Queues) Create(attributes Attrs) (queue *Queue) {
	queue = CreateQueue(queues.AddOwnerToAttributes(attributes))
	queues.items = append(queues.items, queue)
	return queue
}

