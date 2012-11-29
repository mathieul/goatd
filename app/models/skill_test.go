package models_test

import (
    . "launchpad.net/gocheck"
    "testing"
    "goatd/app/models"
)

func Test(t *testing.T) { TestingT(t) }

type SkillSuite struct{
}

var _ = Suite(&SkillSuite{})

func (s *SkillSuite) SetUpTest(c *C) {
}

func (s *SkillSuite) TearDownTest(c *C) {
}

func (s *SkillSuite) TestNewSkill(c *C) {
    skill := models.NewSkill(models.Attrs{"QueueUid": "0abc",
        "TeammateUid": "1def", "Level": models.LevelHigh, "Enabled": true})
    c.Assert(skill.QueueUid(), Equals, "0abc")
    c.Assert(skill.TeammateUid(), Equals, "1def")
    c.Assert(skill.Level(), Equals, models.LevelHigh)
    c.Assert(skill.Enabled(), Equals, true)
    c.Assert(skill.Persisted(), Equals, false)
}

func (s *SkillSuite) TestCreateSkill(c *C) {
    skill := models.CreateSkill(models.Attrs{"QueueUid": "0abc",
        "TeammateUid": "1def", "Level": models.LevelMedium, "Enabled": false})
    c.Assert(skill.Level(), Equals, models.LevelMedium)
    c.Assert(skill.Enabled(), Equals, false)
    c.Assert(skill.Persisted(), Equals, true)
}
