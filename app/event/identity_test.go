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
    identity := event.NewIdentity("Team", "")
    c.Assert(identity.Uid(), HasLen, 17)
    c.Assert(identity.Uid()[8:9], Equals, "-")
}

func (s *IdentitySuite) TestCopyIdentity(c *C) {
    // TODO
}