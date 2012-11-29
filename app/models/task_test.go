package models_test

import (
    . "launchpad.net/gocheck"
    "testing"
    "goatd/app/models"
)

func Test(t *testing.T) { TestingT(t) }

type TaskSuite struct{
}

var _ = Suite(&TaskSuite{})

func (s *TaskSuite) SetUpTest(c *C) {
}

func (s *TaskSuite) TearDownTest(c *C) {
}

func (s *TaskSuite) TestNewTaskWithTitle(c *C) {
    task := models.NewTask(models.Attrs{"Title": "Buy the milk"})
    c.Assert(task.Title(), Equals, "Buy the milk")
    c.Assert(task.Persisted(), Equals, false)
}

func (s *TaskSuite) TestCreateTaskWithTitle(c *C) {
    task := models.CreateTask(models.Attrs{"Title": "Call back the milkman"})
    c.Assert(task.Title(), Equals, "Call back the milkman")
    c.Assert(task.Persisted(), Equals, true)
}
