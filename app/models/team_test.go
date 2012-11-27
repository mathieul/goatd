package models_test

import (
    . "launchpad.net/gocheck"
    "testing"
    "goatd/app/models"
)

func Test(t *testing.T) { TestingT(t) }

type TeamSuite struct{
}

var _ = Suite(&TeamSuite{})

func (s *TeamSuite) SetUpTest(c *C) {
}

func (s *TeamSuite) TearDownTest(c *C) {
}

func (s *TeamSuite) TestCreateTeamWithName(c *C) {
    team := models.CreateTeam(models.Attrs{"Name": "Metallica"})
    c.Assert(team.Name(), Equals, "Metallica")
    c.Assert(len(team.Uid()), Equals, 8 + 1 + 8)
}
