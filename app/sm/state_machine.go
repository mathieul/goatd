package sm

/*
 * Constants and initialized vars
 */
var NoAction Action

func init() {
    NoAction = func (args []interface{}) bool { return true }
}

/*
 * Basic types and interfaces
 */
type Event uint
type Status uint
type Action func ([]interface{}) bool

/*
 * StateMachine
 */
type StateMachine struct {}
type Builder struct {}
type Transitioner struct {}

func (builder Builder) Event(event Event, args ...interface{}) {
}

func (transitioner Transitioner) Transition(from, to Status, callback Action) {
}

func (stateMachine StateMachine) Status() Status {
    return 1
}

func (stateMachine *StateMachine) Trigger(event Event) bool {
    return true
}


func NewStateMachine(status Status, callback func (Builder)) *StateMachine {
    return new(StateMachine)
}