package model_test

import (
	. "launchpad.net/gocheck"
	"testing"
	"goatd/app/model"
)

func Test(t *testing.T) { TestingT(t) }

type Person struct {
	model.AttributeOwner
}

func NewPerson(attributes model.Attributes) (*Person) {
	fields := model.Fields{"firstName", "lastName"}
	owner := model.NewAttributeOwner(&fields, &attributes)
	return &Person{*owner}
}

type ModelSuite struct{
	person *Person
}

var _ = Suite(&ModelSuite{})

func (s *ModelSuite) SetUpTest(c *C) {
	s.person = NewPerson(model.Attributes{"firstName": "John", "lastName": "Doe"})
}

func (s *ModelSuite) TearDownTest(c *C) {
}

func (s *ModelSuite) TestItReadsAttributesWithGet(c *C) {
	firstName, found := s.person.Get("firstName")
	c.Assert(true, Equals, found)
    c.Assert("John", Equals, firstName)

    firstName, found = s.person.Get("unknown")
	c.Assert(false, Equals, found)
    c.Assert(firstName, IsNil)
}
