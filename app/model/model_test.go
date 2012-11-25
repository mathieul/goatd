package model_test

import (
	. "launchpad.net/gocheck"
	"testing"
	"goatd/app/model"
)

func Test(t *testing.T) { TestingT(t) }

type Person struct {
	model.Model
}

func NewPerson(attributes model.Attributes) (*Person) {
	fields := []string{"firstName", "lastName"}
	return &Person{model.NewModel(fields, attributes)}
}

type ModelSuite struct{
	person *Person
}

var _ = Suite(&ModelSuite{})

func (s *ModelSuite) SetUpTest(c *C) {
	s.person = NewPerson(model.Attributes{"firstName", "John", "lastName", "Doe"})
}

func (s *ModelSuite) TearDownTest(c *C) {
}

func (s *ModelSuite) TestItReadsAttributesWithGet(c *C) {
    c.Assert("John", Equals, s.person.Get("firstName"))
    c.Assert("Doe", Equals, s.person.Get("lastName"))
}
