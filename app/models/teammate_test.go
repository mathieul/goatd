package models_test

import (
    . "launchpad.net/gocheck"
    "testing"
    "goatd/app/models"
)

func Test(t *testing.T) { TestingT(t) }

type TeammateSuite struct{
    teammates *models.Teammates
}

var _ = Suite(&TeammateSuite{})

func (s *TeammateSuite) SetUpTest(c *C) {
    s.teammates = models.NewTeammates("Team", "zxcvbnm")
}

func (s *TeammateSuite) TestCreateTeammate(c *C) {
    teammate := s.teammates.Create(models.Attrs{"Name": "Egret"})
    c.Assert(teammate.Name(), Equals, "Egret")
    c.Assert(teammate.TeamUid(), Equals, "zxcvbnm")
    c.Assert(teammate.Persisted(), Equals, true)
}
