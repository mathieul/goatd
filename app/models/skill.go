package models

import (
    "goatd/app/identification"
    "goatd/app/event"
)

/*
 * Skill
 */

type Skill struct {
    Storage
    team *Team
    identity *identification.Identity
    AttrQueueUid string
    AttrTeammateUid string
    AttrLevel int
    AttrEnabled bool
}

func NewSkill(attributes Attrs) *Skill {
    skill := newModel(&Skill{}, &attributes).(*Skill)
    skill.identity = identification.NewIdentity("Skill", skill.Uid(), skill)
    return skill
}

func CreateSkill(attributes Attrs) (skill *Skill) {
    skill = NewSkill(attributes)
    skill.Save()
    event.Manager().PublishEvent(event.KindSkillCreated, *skill.identity,
        []string{
            skill.QueueUid(),
            skill.TeammateUid(),
            levelToString[skill.Level()],
            boolToString[skill.Enabled()],
        })
    return skill
}

func (skill *Skill) QueueUid() string {
    return skill.AttrQueueUid
}

func (skill *Skill) TeammateUid() string {
    return skill.AttrTeammateUid
}

func (skill *Skill) Level() int {
    return skill.AttrLevel
}

func (skill *Skill) Enabled() bool {
    return skill.AttrEnabled
}

func (skill *Skill) SetTeam(team *Team) {
    skill.team = team
}

func (skill Skill) Team() (team *Team) {
    return skill.team
}


/*
 * Skills
 */

type Skills struct {
    Collection
}

func toSkillSlice(source []interface{}) []*Skill {
    skills := make([]*Skill, 0, len(source))
    for _, skill := range source {
        skills = append(skills, skill.(*Skill))
    }
    return skills
}

func NewSkills(owner identification.Identity) (skills *Skills) {
    skills = new(Skills)
    skills.Collection = NewCollection(func(attributes Attrs, owner interface{}) interface{} {
        skill := CreateSkill(attributes)
        skill.SetTeam(owner.(*Team))
        return skill
    }, owner)
    return skills
}

func (skills *Skills) Create(attributes Attrs) (skill *Skill) {
    return skills.Collection.Create(attributes).(*Skill)
}

func (skills Skills) Slice() []*Skill {
    return toSkillSlice(skills.Collection.Slice())
}

func (skills Skills) Find(uid string) *Skill {
    if found := skills.Collection.Find(uid); found != nil {
        return found.(*Skill)
    }
    return nil
}

func (skills Skills) FindAll(uids []string) []*Skill {
    return toSkillSlice(skills.Collection.FindAll(uids))
}

func (skills Skills) Select(tester func(interface{}) bool) (result []*Skill) {
    return toSkillSlice(skills.Collection.Select(tester))
}
