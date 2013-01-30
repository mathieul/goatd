package model_test

import (
    . "launchpad.net/gocheck"
    "goatd/app/event"
    "goatd/app/model"
)

func skillUids(skills []*model.Skill) (uids []string) {
    for _, skill := range skills {
        uids = append(uids, skill.Uid())
    }
    return uids
}

type SkillOwner struct {
    *event.Identity
}

type SkillSuite struct {
    store *model.Store
    owner *SkillOwner
}

var _ = Suite(&SkillSuite{})

func (s *SkillSuite) SetUpTest(c *C) {
    s.store = model.NewStore(nil)
    s.owner = &SkillOwner{event.NewIdentity("Team")}
}

func (s *SkillSuite) TestCreateSkill(c *C) {
    skill := s.store.Skills.Create(model.A{"QueueUid": "queue123", "TeammateUid": "teammate456"}, s.owner)
    c.Assert(skill.QueueUid(), Equals, "queue123")
    c.Assert(skill.TeammateUid(), Equals, "teammate456")
    c.Assert(skill.TeamUid(), Equals, s.owner.Uid())

    queue := SkillOwner{event.NewIdentity("Queue")}
    teammate := SkillOwner{event.NewIdentity("Teammate")}
    alternative := s.store.Skills.Create(model.A{}, s.owner, queue, teammate)
    c.Assert(alternative.QueueUid(), Equals, queue.Uid())
    c.Assert(alternative.TeammateUid(), Equals, teammate.Uid())
    c.Assert(alternative.TeamUid(), Equals, s.owner.Uid())
}

func (s *SkillSuite) TestFindSkill(c *C) {
    s.store.Skills.Create(model.A{"QueueUid": "queue123", "TeammateUid": "teammate456"}, s.owner)
    s2 := s.store.Skills.Create(model.A{"QueueUid": "queue123", "TeammateUid": "teammate789"}, s.owner)
    found := s.store.Skills.Find(s2.Uid())
    c.Assert(found.TeammateUid(), DeepEquals, "teammate789")
    c.Assert(s.store.Skills.Find("unknown"), IsNil)
}

func (s *SkillSuite) TestFindAllSkills(c *C) {
    s1 := s.store.Skills.Create(model.A{"QueueUid": "queue123", "TeammateUid": "teammate456"}, s.owner)
    s.store.Skills.Create(model.A{"QueueUid": "queue123", "TeammateUid": "teammate789"}, s.owner)
    s3 := s.store.Skills.Create(model.A{"QueueUid": "queue456", "TeammateUid": "teammate007"}, s.owner)
    foundSkills := s.store.Skills.FindAll([]string{s1.Uid(), s3.Uid()})
    c.Assert(skillUids(foundSkills), DeepEquals, []string{s1.Uid(), s3.Uid()})
}

func (s *SkillSuite) TestUpdateSkill(c *C) {
    skill := s.store.Skills.Create(model.A{"QueueUid": "queue123", "TeammateUid": "teammate456"}, s.owner)
    skill.Update("QueueUid", "queue007")
    c.Assert(skill.QueueUid(), Equals, "queue007")
    found := s.store.Skills.Find(skill.Uid())
    c.Assert(found.QueueUid(), Equals, "queue007")
}

func (s *SkillSuite) TestDestroySkill(c *C) {
    skill := s.store.Skills.Create(model.A{"Name": "Destroy me"}, s.owner)
    c.Assert(skill.Destroy().Uid(), Equals, skill.Uid())
    c.Assert(s.store.Skills.Find(skill.Uid()), IsNil)
}

func (s *SkillSuite) TestSelectSkills(c *C) {
    s1 := s.store.Skills.Create(model.A{"QueueUid": "queue007", "TeammateUid": "teammate_james"}, s.owner)
    s.store.Skills.Create(model.A{"QueueUid": "queue123", "TeammateUid": "teammate456"}, s.owner)
    s3 := s.store.Skills.Create(model.A{"QueueUid": "queue007", "TeammateUid": "teammate_bond"}, s.owner)
    selectedSkills := s.store.Skills.Select(func (item interface{}) bool {
        skill := item.(*model.Skill)
        return skill.QueueUid() == "queue007"
    })
    c.Assert(skillUids(selectedSkills), DeepEquals, []string{s1.Uid(), s3.Uid()})
}
