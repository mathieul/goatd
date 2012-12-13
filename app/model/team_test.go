package model_test

import (
    . "launchpad.net/gocheck"
    "goatd/app/model"
    "strings"
)


type TeamSuite struct{
    teams *model.Teams
    team *model.Team
}

var _ = Suite(&TeamSuite{})

func (s *TeamSuite) SetUpTest(c *C) {
    s.teams = model.NewTeams()
    s.team = s.teams.New(model.A{"Name": "Jon Snow & Egret"})
}

func (s *TeamSuite) TestNewTeam(c *C) {
    c.Assert(s.team.Name(), Equals, "Jon Snow & Egret")
    c.Assert(s.team.Uid(), HasLen, 8 + 1 + 8)
}

func (s *TeamSuite) TestReturnsSlice(c *C) {
    t1 := s.teams.New(model.A{"Name": "Lannister"})
    t2 := s.teams.New(model.A{"Name": "Stark"})
    t3 := s.teams.New(model.A{"Name": "Baratheon"})
    c.Assert(s.teams.Slice(), DeepEquals, []*model.Team{s.team, t1, t2, t3})
}

func (s *TeamSuite) TestFindTeam(c *C) {
    uid := s.teams.New(model.A{"Name": "Metallica"}).Uid()
    s.teams.New(model.A{"Name": "Masada"})
    team := s.teams.Find(uid)
    c.Assert(team.Name(), Equals, "Metallica")
    team = s.teams.Find("unknown")
    c.Assert(team, IsNil)
}

func (s *TeamSuite) TestFindAllTeams(c *C) {
    uid1 := s.teams.New(model.A{"Name": "One"}).Uid()
    s.teams.New(model.A{"Name": "Two"})
    s.teams.New(model.A{"Name": "Three"})
    uid2 := s.teams.New(model.A{"Name": "Four"}).Uid()
    s.teams.New(model.A{"Name": "Five"})

    found := s.teams.FindAll([]string{uid1, uid2})
    c.Assert(found[0].Name(), DeepEquals, "One")
    c.Assert(found[1].Name(), DeepEquals, "Four")
}

func (s *TeamSuite) TestSelectTeams(c *C) {
    tyrion := s.teams.New(model.A{"Name": "Tyrion Lannister"})
    s.teams.New(model.A{"Name": "Jon Snow"})
    jamie := s.teams.New(model.A{"Name": "Jamie Lannister"})
    c.Assert(s.teams.Select(func (item interface{}) bool {
            team := item.(*model.Team)
            return strings.Contains(team.Name(), "Lannister")
        }),
        DeepEquals,
        []*model.Team{tyrion, jamie})
}

// func (s *TeamSuite) TestNewTeammate(c *C) {
//     teammate := s.team.Teammates.New(model.A{"Name": "Kirk Hammett"})
//     c.Assert(teammate.Name(), Equals, "Kirk Hammett")
//     c.Assert(teammate.TeamUid(), Equals, s.team.Uid())
//     c.Assert(teammate.Team(), DeepEquals, s.team)
// }

// func (s *TeamSuite) TestNewQueue(c *C) {
//     queue := s.team.Queues.New(model.A{"Name": "Support"})
//     c.Assert(queue.Name(), Equals, "Support")
//     c.Assert(queue.TeamUid(), Equals, s.team.Uid())
// }

// func (s *TeamSuite) TestNewSkill(c *C) {
//     skill := s.team.Skills.New(model.A{"QueueUid": "0abc", "TeammateUid": "1def"})
//     c.Assert(skill.QueueUid(), Equals, "0abc")
//     c.Assert(skill.TeammateUid(), Equals, "1def")
// }

// func (s *TeamSuite) TestNewTask(c *C) {
//     task := s.team.Tasks.New(model.A{"Title": "Buy the milk"})
//     c.Assert(task.Title(), Equals, "Buy the milk")
// }
