package sm

import (
    "fmt"
    "log"
)

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
type statusTargetMap map[Status]statusTarget
type transitionMap map[Event]statusTargetMap

/*
 * Builder
 */
type statusTarget struct {
    value  Status
    action Action
}

type Builder struct {
    stateMachine *StateMachine
}

func newBuilder(stateMachine *StateMachine) *Builder {
    return &Builder{stateMachine}
}

func (builder Builder) EventSingleTransition(event Event, from, to Status, action Action) {
    target := statusTarget{to, action}
    if transitions, found := builder.stateMachine.eventTransitions[event]; found {
        fmt.Println("EventSingleTransition(", event, "): found => set ", from)
        transitions[from] = target
    } else {
        fmt.Println("EventSingleTransition(", event, "): not found - allocate")
        transitions = make(statusTargetMap, 1)
        transitions[from] = target
        builder.stateMachine.eventTransitions[event] = transitions
    }
}

func (builder Builder) EventMultiTransitions(event Event, callback func(Transitioner)) {
}

func (builder Builder) Event(event Event, args ...interface{}) {
    switch len(args) {
    case 1:
        builder.EventMultiTransitions(event, args[0].(func(Transitioner)))
    case 3:
        builder.EventSingleTransition(event, args[0].(Status), args[1].(Status), args[2].(Action))
    default:
        log.Fatal(fmt.Errorf("sm.Event(): invalid call with %d arguments", len(args) + 1))
    }
}

/*
 * Transitioner
 */
type Transitioner struct {}

func (transitioner Transitioner) Transition(from, to Status, callback Action) {
}

/*
 * StateMachine
 */
type StateMachine struct {
    status Status
    eventTransitions transitionMap
}

func (stateMachine StateMachine) Status() Status {
    return stateMachine.status
}

func (stateMachine *StateMachine) Trigger(event Event) bool {
    if transitions, found := stateMachine.eventTransitions[event]; found {
        if target, listed := transitions[stateMachine.status]; listed {
            stateMachine.status = target.value
            if target.action != nil {
                target.action([]interface{}{})
            }
            return true
        }
    }
    return false
}

func NewStateMachine(status Status, callback func (Builder)) (stateMachine *StateMachine) {
    eventTransitions := make(transitionMap, 5)
    stateMachine = &StateMachine{status, eventTransitions}
    callback(*newBuilder(stateMachine))
    return stateMachine
}
