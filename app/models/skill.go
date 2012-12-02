package models

import (
    "goatd/app/event"
)

/*
 * Skill
 */

type Skill struct {
    Storage
    AttrQueueUid string
    AttrTeammateUid string
    AttrLevel int
    AttrEnabled bool
}

func NewSkill(attributes Attrs) *Skill {
    return newModel(&Skill{}, &attributes).(*Skill)
}

func CreateSkill(attributes Attrs) (skill *Skill) {
    skill = NewSkill(attributes)
    skill.Save()
    return skill
}

func (team *Skill) QueueUid() string {
    return team.AttrQueueUid
}

func (team *Skill) TeammateUid() string {
    return team.AttrTeammateUid
}

func (team *Skill) Level() int {
    return team.AttrLevel
}

func (team *Skill) Enabled() bool {
    return team.AttrEnabled
}

/*
 * Skills
 */

type Skills struct {
    owner event.Identity
    items []*Skill
}

func NewSkills(owner event.Identity) (skills *Skills) {
    skills = new(Skills)
    skills.owner = owner
    return skills
}

func (skills *Skills) Create(attributes Attrs) (skill *Skill) {
    skill = CreateSkill(skills.owner.AddToAttributes(attributes))
    skills.items = append(skills.items, skill)
    return skill
}
