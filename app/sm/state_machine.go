package sm

import (
    "fmt"
    "log"
)

/*
 * Constants and initialized vars
 */
var NoAction func ([]interface{}) bool

func init() {
    NoAction = func (args []interface{}) bool { return true }
}

/*
 * Basic types and interfaces
 */
type Event uint
type Status uint
type statusTargetMap map[Status]statusTarget
type transitionMap map[Event]statusTargetMap

/*
 * Builder
 */
type statusTarget struct {
    value  Status
    action func ([]interface{}) bool
}

type Builder struct {
    stateMachine *StateMachine
    event Event
}

func newBuilder(stateMachine *StateMachine) *Builder {
    return &Builder{stateMachine, 0}
}

func (builder Builder) EventSingleTransition(event Event, from, to Status, action func ([]interface{}) bool) {
    newBuilder := Builder{builder.stateMachine, event}
    newBuilder.Transition(from, to, action)
}

func (builder Builder) EventMultiTransitions(event Event, callback func(Builder)) {
    newBuilder := Builder{builder.stateMachine, event}
    callback(newBuilder)
}

func (builder Builder) Transition(from, to Status, action func ([]interface{}) bool) {
    target := statusTarget{to, action}
    if transitions, found := builder.stateMachine.eventTransitions[builder.event]; found {
        fmt.Println("EventSingleTransition(", builder.event, "): found => set ", from)
        transitions[from] = target
    } else {
        fmt.Println("EventSingleTransition(", builder.event, "): not found - allocate")
        transitions = make(statusTargetMap, 1)
        transitions[from] = target
        builder.stateMachine.eventTransitions[builder.event] = transitions
    }
}

func (builder Builder) Event(event Event, args ...interface{}) {
    switch len(args) {
    case 1:
        builder.EventMultiTransitions(event, args[0].(func(Builder)))
    case 3:
        builder.EventSingleTransition(event, args[0].(Status), args[1].(Status), args[2].(func ([]interface{}) bool))
    default:
        log.Fatal(fmt.Errorf("sm.Event(): invalid call with %d arguments", len(args) + 1))
    }
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

func (stateMachine *StateMachine) Trigger(event Event, args ...interface{}) bool {
    if transitions, found := stateMachine.eventTransitions[event]; found {
        if target, listed := transitions[stateMachine.status]; listed {
            stateMachine.status = target.value
            if target.action != nil {
                target.action(args)
            }
            return true
        }
    }
    return false
}

func NewStateMachine(status Status, callback func (Builder)) (stateMachine *StateMachine) {
    eventTransitions := make(transitionMap, 5)
    stateMachine = &StateMachine{status, eventTransitions}
    builder := *newBuilder(stateMachine)
    callback(builder)
    return stateMachine
}
