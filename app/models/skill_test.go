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
