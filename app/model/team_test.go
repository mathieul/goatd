package model_test

import (
    . "launchpad.net/gocheck"
    "goatd/app/model"
    "strings"
)


type TeamSuite struct {
    store *model.Store
    team *model.Team
}

var _ = Suite(&TeamSuite{})

func (s *TeamSuite) SetUpTest(c *C) {
    s.store = model.NewStore(nil)
    s.team = s.store.Teams.Create(model.A{"Name": "Jon Snow & Egret"})
}

func (s *TeamSuite) TestCreateTeam(c *C) {
    c.Assert(s.team.Name(), Equals, "Jon Snow & Egret")
    c.Assert(s.team.Uid(), HasLen, 8 + 1 + 8)
}

func (s *TeamSuite) TestFindTeam(c *C) {
    uid := s.store.Teams.Create(model.A{"Name": "Metallica"}).Uid()
    s.store.Teams.Create(model.A{"Name": "Masada"})
    team := s.store.Teams.Find(uid)
    c.Assert(team.Name(), Equals, "Metallica")
    team = s.store.Teams.Find("unknown")
    c.Assert(team, IsNil)
}

func (s *TeamSuite) TestFindAllTeams(c *C) {
    uid1 := s.store.Teams.Create(model.A{"Name": "One"}).Uid()
    s.store.Teams.Create(model.A{"Name": "Two"})
    s.store.Teams.Create(model.A{"Name": "Three"})
    uid2 := s.store.Teams.Create(model.A{"Name": "Four"}).Uid()
    s.store.Teams.Create(model.A{"Name": "Five"})

    found := s.store.Teams.FindAll([]string{uid1, uid2})
    c.Assert(found[0].Name(), DeepEquals, "One")
    c.Assert(found[1].Name(), DeepEquals, "Four")

    notFound := s.store.Teams.FindAll([]string{"blahbalh"})
    c.Assert(notFound, HasLen, 0)
}

func (s *TeamSuite) TestDestroyTeam(c *C) {
    team := s.store.Teams.Create(model.A{"Name": "Destroy me"})
    c.Assert(team.Destroy().Uid(), Equals, team.Uid())
    c.Assert(s.store.Teams.Find(team.Uid()), IsNil)
}

func (s *TeamSuite) TestUpdateTeam(c *C) {
    team := s.store.Teams.Create(model.A{"Name": "Metallica"})
    team.Update("Name", "Iron Maiden")
    c.Assert(team.Name(), Equals, "Iron Maiden")
    found := s.store.Teams.Find(team.Uid())
    c.Assert(found.Name(), Equals, "Iron Maiden")
}

func (s *TeamSuite) TestSelectTeams(c *C) {
    tyrion := s.store.Teams.Create(model.A{"Name": "Tyrion Lannister"})
    s.store.Teams.Create(model.A{"Name": "Jon Snow"})
    jamie := s.store.Teams.Create(model.A{"Name": "Jamie Lannister"})
    c.Assert(s.store.Teams.Select(func (item interface{}) bool {
            team := item.(*model.Team)
            return strings.Contains(team.Name(), "Lannister")
        }),
        DeepEquals,
        []*model.Team{tyrion, jamie})
}

func (s *TeamSuite) TestCountAndEachTeams(c *C) {
    s.store.Teams.DestroyAll()
    c.Assert(s.store.Teams.Count(), Equals, 0)

    s.store.Teams.Create(model.A{"Name": "One"})
    s.store.Teams.Create(model.A{"Name": "Two"})

    c.Assert(s.store.Teams.Count(), Equals, 2)
    names := make([]string, 0)
    s.store.Teams.Each(func (item interface{}) {
        team := item.(*model.Team)
        names = append(names, team.Name())
    })
    c.Assert(names, DeepEquals, []string{"One", "Two"})
}
