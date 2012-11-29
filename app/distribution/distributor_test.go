package distribution_test

import (
    . "launchpad.net/gocheck"
    "testing"
    "goatd/app/distribution"
    "goatd/app/models"
)

func Test(t *testing.T) { TestingT(t) }

type DistirbutorSuite struct{
}

var _ = Suite(&DistirbutorSuite{})

func (s *DistirbutorSuite) SetUpTest(c *C) {
}

func (s *DistirbutorSuite) TearDownTest(c *C) {
}

func (s *DistirbutorSuite) TestNewDistributor(c *C) {
    team := models.NewTeam(models.Attrs{"Name": "Get To Work"})
    distributor := distribution.NewDistributor(team)
    c.Assert(distributor.Team(), DeepEquals, team)
}
