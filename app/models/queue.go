package models

type Queue struct {
    Storage
    AttrName string
    AttrTeamUid string
}

func NewQueue(attributes Attrs) *Queue {
    return newModel(&Queue{}, &attributes).(*Queue)
}

func CreateQueue(attributes Attrs) (queue *Queue) {
    queue = NewQueue(attributes)
    queue.Save()
    return queue
}

func (team *Queue) Name() string {
    return team.AttrName
}

func (team *Queue) TeamUid() string {
    return team.AttrTeamUid
}
