package models_test

import (
    . "launchpad.net/gocheck"
    "testing"
    "goatd/app/models"
)

func Test(t *testing.T) { TestingT(t) }

type SkillSuite struct {
    team *models.Team
    skills *models.Skills
}

var _ = Suite(&SkillSuite{})

func (s *SkillSuite) SetUpTest(c *C) {
    s.team = models.NewTeam(models.Attrs{"Name": "Hello, inc."})
    s.skills = s.team.Skills
}

func (s *SkillSuite) TestCreateSkill(c *C) {
    skill := s.skills.Create(models.Attrs{"QueueUid": "0abc",
        "TeammateUid": "1def", "Level": models.LevelMedium, "Enabled": false})
    c.Assert(skill.QueueUid(), Equals, "0abc")
    c.Assert(skill.TeammateUid(), Equals, "1def")
    c.Assert(skill.Level(), Equals, models.LevelMedium)
    c.Assert(skill.Enabled(), Equals, false)
    c.Assert(skill.Persisted(), Equals, true)
}

func (s *SkillSuite) TestReturnsSlice(c *C) {
    s1 := s.skills.Create(models.Attrs{"QueueUid": "01", "TeammateUid": "t1", "Level": models.LevelLow})
    s2 := s.skills.Create(models.Attrs{"QueueUid": "01", "TeammateUid": "t2", "Level": models.LevelMedium})
    s3 := s.skills.Create(models.Attrs{"QueueUid": "01", "TeammateUid": "t3", "Level": models.LevelHigh})
    c.Assert(s.skills.Slice(), DeepEquals, []*models.Skill{s1, s2, s3})
}

func (s *SkillSuite) TestFindSkill(c *C) {
    s.skills.Create(models.Attrs{"QueueUid": "01", "TeammateUid": "t1"})
    s2 := s.skills.Create(models.Attrs{"QueueUid": "01", "TeammateUid": "t2"})
    c.Assert(s.skills.Find(s2.Uid()), DeepEquals, s2)
}

func (s *SkillSuite) TestFindSkillSlice(c *C) {
    s1 := s.skills.Create(models.Attrs{"QueueUid": "01", "TeammateUid": "t1"})
    s.skills.Create(models.Attrs{"QueueUid": "01", "TeammateUid": "t2"})
    s3 := s.skills.Create(models.Attrs{"QueueUid": "01", "TeammateUid": "t3"})
    c.Assert(s.skills.FindAll([]string{s1.Uid(), s3.Uid()}),
             DeepEquals,
             []*models.Skill{s1, s3})
}

func (s *SkillSuite) TestSelectSkills(c *C) {
    s.skills.Create(models.Attrs{"QueueUid": "01", "TeammateUid": "t1", "Level": models.LevelMedium})
    matches := s.skills.Create(models.Attrs{"QueueUid": "02", "TeammateUid": "t1", "Level": models.LevelHigh})
    s.skills.Create(models.Attrs{"QueueUid": "03", "TeammateUid": "t2", "Level": models.LevelMedium})
    c.Assert(s.skills.Select(func (item interface{}) bool {
            skill := item.(*models.Skill)
            return skill.Level() == models.LevelHigh && skill.TeammateUid() == "t1"
        }),
        DeepEquals,
        []*models.Skill{matches})
}
