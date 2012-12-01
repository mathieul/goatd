package models_test

import (
    . "launchpad.net/gocheck"
    "testing"
    "goatd/app/models"
)

func Test(t *testing.T) { TestingT(t) }

type SkillSuite struct {
    skills *models.Skills
}

var _ = Suite(&SkillSuite{})

func (s *SkillSuite) SetUpTest(c *C) {
    s.skills = models.NewSkills("Team", "inc")
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
