package event_test

import (
    . "launchpad.net/gocheck"
    "goatd/app/event"
)

type IdentitySuite struct {}

var _ = Suite(&IdentitySuite{})

func (s *IdentitySuite) TestNewIdentity(c *C) {
    identity := event.NewIdentity("Team", "blah")
    c.Assert(identity.Kind(), Equals, "Team")
    c.Assert(identity.Uid(), Equals, "blah")
}

func (s *IdentitySuite) TestGenerateUid(c *C) {
    identity := event.NewIdentity("Team")
    c.Assert(identity.Uid(), HasLen, 17)
    c.Assert(identity.Uid()[8:9], Equals, "-")
}

func (s *IdentitySuite) TestCopyIdentity(c *C) {
    original := event.NewIdentity("Team")
    copy := original.Copy(original)
    c.Assert(copy.Kind(), Equals, original.Kind())
    c.Assert(copy.Uid(), Equals, original.Uid())
    copy.SetKind("Teammate")
    copy.SetUid("blah-blah")
    c.Assert(copy.Kind(), Not(Equals), original.Kind())
    c.Assert(copy.Uid(), Not(Equals), original.Uid())
}

func (s *IdentitySuite) TestNoIdentity(c *C) {
    identity := event.NewIdentity()
    c.Assert(event.NoIdentity(identity), Equals, true)
    identity.SetUid("123abc")
    c.Assert(event.NoIdentity(identity), Equals, false)
}
