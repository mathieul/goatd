package model

import (
    "goatd/app/event"
    "goatd/app/sm"
)


/*
 * Global
 */
const (
    LevelNone int = iota
    LevelLow
    LevelMedium
    LevelHigh
)


/*
 * Skill
 */

type Skill struct {
    *event.Identity
    busManager *event.BusManager
    store *Store
    AttrQueueUid string
    AttrTeammateUid string
    AttrLevel int
    AttrTeamUid string
}

func NewSkill(attributes A) (skill *Skill) {
    skill = newModel(&Skill{}, &attributes).(*Skill)
    if skill.AttrLevel == LevelNone { skill.AttrLevel = LevelMedium }
    skill.Identity = event.NewIdentity("Skill")
    return skill
}

func (skill *Skill) Copy() Model {
    identity := skill.Identity.Copy()
    return &Skill{identity, nil, nil, skill.AttrQueueUid,
        skill.AttrTeammateUid, skill.AttrLevel, skill.AttrTeamUid}
}

func (skill *Skill) SetupComs(busManager *event.BusManager, store *Store) {
    skill.busManager = busManager
    skill.store = store
}

func (skill *Skill) Update(name string, value interface{}) bool {
    setAttributeValue(skill, name, value)
    return skill.store.Update(KindSkill, skill.Uid(), name, value)
}

func (skill Skill) Reload() *Skill {
    if found := skill.store.Skills.Find(skill.Uid()); found != nil {
        return found
    }
    return nil
}

func (skill Skill) Status(_ ...sm.Status) sm.Status {
    return StatusNone
}

func (skill Skill) QueueUid() string {
    return skill.AttrQueueUid
}

func (skill Skill) TeammateUid() string {
    return skill.AttrTeammateUid
}

func (skill Skill) Level() int {
    return skill.AttrLevel
}

func (skill Skill) TeamUid() string {
    return skill.AttrTeamUid
}


/*
 * SkillStoreProxy
 */

type SkillStoreProxy struct {
    store *Store
}

func toSkillSlice(source []Model) []*Skill {
    skills := make([]*Skill, 0, len(source))
    for _, skill := range source {
        skills = append(skills, skill.(*Skill))
    }
    return skills
}

func (proxy *SkillStoreProxy) Create(attributes A, owners ...event.Identified) *Skill {
    for _, owner := range owners { attributes = owner.AddToAttributes(attributes) }
    return proxy.store.Create(KindSkill, attributes).(*Skill)
}

func (proxy *SkillStoreProxy) Find(uid string) *Skill {
    if value := proxy.store.Find(KindSkill, uid); value != nil { return value.(*Skill) }
    return nil
}

func (proxy *SkillStoreProxy) FindAll(uids []string) []*Skill {
    values := proxy.store.FindAll(KindSkill, uids)
    return toSkillSlice(values)
}

func (proxy *SkillStoreProxy) Select(tester func(interface{}) bool) []*Skill {
    values := proxy.store.Select(KindSkill, tester)
    return toSkillSlice(values)
}
