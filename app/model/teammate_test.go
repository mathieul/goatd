package model_test

import (
    . "launchpad.net/gocheck"
    "strings"
    "goatd/app/event"
    "goatd/app/model"
)

func teammateNames(teammates []*model.Teammate) (names []string) {
    for _, teammate := range teammates {
        names = append(names, teammate.Name())
    }
    return names
}

type TeammateOwner struct {
    *event.Identity
}

type TeammateSuite struct {
    store *model.Store
    teammate *model.Teammate
    owner *TeammateOwner
    busManager *event.BusManager
}

var _ = Suite(&TeammateSuite{})

func (s *TeammateSuite) SetUpTest(c *C) {
    s.busManager = event.NewBusManager()
    s.busManager.Start()
    s.store = model.NewStore(s.busManager)
    s.owner = &TeammateOwner{event.NewIdentity("Team")}
    s.teammate = s.store.Teammates.Create(model.A{"Name": "Agent"}, s.owner)
}

func (s *TeammateSuite) TearDownTest(c *C) {
    s.busManager.Stop()
}

func (s *TeammateSuite) TestCreateTeammate(c *C) {
    teammate := s.store.Teammates.Create(model.A{"Name": "Jon"}, s.owner)
    c.Assert(teammate.Name(), Equals, "Jon")
    c.Assert(teammate.TeamUid(), Equals, s.owner.Uid())
    c.Assert(teammate.Status(), Equals, model.StatusSignedOut)
}

func (s *TeammateSuite) TestCopyTeammate(c *C) {
    dog := model.NewTeammate(model.A{"Name": "The Hound"})
    c.Assert(dog.IsCopy(), Equals, false)
    rorge := dog.Copy().(*model.Teammate)
    c.Assert(rorge.IsCopy(), Equals, true)
    c.Assert(rorge.Name(), Equals, "The Hound")
    c.Assert(rorge.Status(), Equals, model.StatusSignedOut)
}

func (s *TeammateSuite) TestFindTeammate(c *C) {
    s.store.Teammates.Create(model.A{"Name": "Jon"}, s.owner)
    egret := s.store.Teammates.Create(model.A{"Name": "Egret"}, s.owner)
    found := s.store.Teammates.Find(egret.Uid())
    c.Assert(found.Name(), DeepEquals, "Egret")
    c.Assert(s.store.Teammates.Find("unknown"), IsNil)
}

func (s *TeammateSuite) TestUpdateTeammate(c *C) {
    teammate := s.store.Teammates.Create(model.A{"Name": "Jon"}, s.owner)
    teammate.Update("Name", "Egret")
    c.Assert(teammate.Name(), Equals, "Egret")
    found := s.store.Teammates.Find(teammate.Uid())
    c.Assert(found.Name(), Equals, "Egret")
}

func (s *TeammateSuite) TestFindAllTeammates(c *C) {
    t1 := s.store.Teammates.Create(model.A{"Name": "One"}, s.owner)
    s.store.Teammates.Create(model.A{"Name": "Two"}, s.owner)
    t3 := s.store.Teammates.Create(model.A{"Name": "Three"}, s.owner)
    foundTeammates := s.store.Teammates.FindAll([]string{t1.Uid(), t3.Uid()})
    c.Assert(teammateNames(foundTeammates), DeepEquals, []string{"One", "Three"})
}

func (s *TeammateSuite) TestSelectTeammates(c *C) {
    s.store.Teammates.Create(model.A{"Name": "Tyrion Lannister"}, s.owner)
    s.store.Teammates.Create(model.A{"Name": "Jon Snow"}, s.owner)
    s.store.Teammates.Create(model.A{"Name": "Jamie Lannister"}, s.owner)
    selectedTeammates := s.store.Teammates.Select(func (item interface{}) bool {
        teammate := item.(*model.Teammate)
        return strings.Contains(teammate.Name(), "Lannister")
    })
    c.Assert(teammateNames(selectedTeammates), DeepEquals, []string{"Tyrion Lannister", "Jamie Lannister"})
}

func (s *TeammateSuite) TestSignInSignOutTeammate(c *C) {
    c.Assert(s.teammate.Status(), Equals, model.StatusSignedOut)
    c.Assert(s.teammate.SignIn(), Equals, true)
    c.Assert(s.teammate.Status(), Equals, model.StatusOnBreak)
    c.Assert(s.teammate.SignOut(), Equals, true)
    c.Assert(s.teammate.Status(), Equals, model.StatusSignedOut)
}

func (s *TeammateSuite) TestChangingAvailability(c *C) {
    s.teammate.SignIn()
    c.Assert(s.teammate.MakeAvailable(), Equals, true)
    c.Assert(s.teammate.Status(), Equals, model.StatusWaiting)
    task := s.store.Tasks.Create(model.A{"Title": "Do It"}, s.owner)
    c.Assert(s.teammate.OfferTask(task.Uid()), Equals, true)
    c.Assert(s.teammate.Status(), Equals, model.StatusOffered)
    c.Assert(s.teammate.TaskUid(), DeepEquals, task.Uid())
}

func (s *TeammateSuite) TestAcceptFinishTask(c *C) {
    s.teammate.SignIn()
    s.teammate.MakeAvailable()
    task := s.store.Tasks.Create(model.A{"Title": "Do It"}, s.owner)
    s.teammate.OfferTask(task.Uid())
    c.Assert(s.teammate.AcceptTask(task.Uid()), Equals, true)
    c.Assert(s.teammate.Status(), Equals, model.StatusBusy)
    c.Assert(s.teammate.TaskUid(), DeepEquals, task.Uid())

    c.Assert(s.teammate.FinishTask(task.Uid()), Equals, true)
    c.Assert(s.teammate.Status(), Equals, model.StatusWrappingUp)
    c.Assert(s.teammate.TaskUid(), Equals, "")
}

func (s *TeammateSuite) TestOtherWorkOnBreakTask(c *C) {
    s.teammate.SignIn()
    s.teammate.MakeAvailable()
    task := s.store.Tasks.Create(model.A{"Title": "Do It"}, s.owner)
    s.teammate.OfferTask(task.Uid())
    s.teammate.AcceptTask(task.Uid())
    c.Assert(s.teammate.StartOtherWork(), Equals, false)
    s.teammate.FinishTask(task.Uid())

    c.Assert(s.teammate.StartOtherWork(), Equals, true)
    c.Assert(s.teammate.Status(), Equals, model.StatusOtherWork)
    c.Assert(s.teammate.GoOnBreak(), Equals, true)
    c.Assert(s.teammate.Status(), Equals, model.StatusOnBreak)
    c.Assert(s.teammate.StartOtherWork(), Equals, true)
    c.Assert(s.teammate.Status(), Equals, model.StatusOtherWork)
}

func (s *TeammateSuite) TestStatusUpdateIsPersistent(c *C) {
    s.teammate.SignIn()
    s.teammate.MakeAvailable()
    found := s.store.Teammates.Find(s.teammate.Uid())
    c.Assert(found.Status(), Equals, model.StatusWaiting)
}
