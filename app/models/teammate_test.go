package models_test

import (
    . "launchpad.net/gocheck"
    "testing"
    "goatd/app/models"
)

func Test(t *testing.T) { TestingT(t) }

type TeammateSuite struct{
}

var _ = Suite(&TeammateSuite{})

func (s *TeammateSuite) SetUpTest(c *C) {
}

func (s *TeammateSuite) TearDownTest(c *C) {
}

func (s *TeammateSuite) TestNewTeammateWithNameAndTeamUid(c *C) {
    teammate := models.NewTeammate(models.Attrs{"Name": "Jon", "TeamUid": "zxcvbnm"})
    c.Assert(teammate.Name(), Equals, "Jon")
    c.Assert(teammate.TeamUid(), Equals, "zxcvbnm")
    c.Assert(teammate.Persisted(), Equals, false)
}

func (s *TeammateSuite) TestCreateTeammateWithNameAndTeamUid(c *C) {
    teammate := models.CreateTeammate(models.Attrs{"Name": "Egret", "TeamUid": "zxcvbnm"})
    c.Assert(teammate.Name(), Equals, "Egret")
    c.Assert(teammate.TeamUid(), Equals, "zxcvbnm")
    c.Assert(teammate.Persisted(), Equals, true)
}
