package sm

import (
    "fmt"
    "log"
)

/*
 * Constants and initialized vars
 */
var NoAction func ([]interface{}) bool
var eventToString map[Event]string
var statusToString map[Status]string

func init() {
    NoAction = func ([]interface{}) bool { return true }
    eventToString = make(map[Event]string, 20)
    statusToString = make(map[Status]string, 20)
}

/*
 * Event and Status
 */
type Event uint
type Status uint
type statusTargetMap map[Status]statusTarget
type transitionMap map[Event]statusTargetMap

func (event Event) String() string {
    if label, found := eventToString[event]; found {
        return fmt.Sprintf("Event{%s}", label)
    }
    return "Event{?Unknown?}"
}

func (status Status) String() string {
    if label, found := statusToString[status]; found {
        return fmt.Sprintf("Status{%s}", label)
    }
    return "Status{?Unknown?}"
}

func RegisterEvent(event Event, label string) {
    eventToString[event] = label
}

func RegisterStatus(status Status, label string) {
    statusToString[status] = label
}

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

func newBuilder(stateMachine *StateMachine, events ...Event) *Builder {
    builder := Builder{stateMachine, 0}
    if len(events) > 0 {
        builder.event = events[0]
    }
    return &builder
}

func (builder Builder) EventSingleTransition(event Event, from, to Status, action func ([]interface{}) bool) {
    newBuilder(builder.stateMachine, event).Transition(from, to, action)
}

func (builder Builder) EventMultiTransitions(event Event, callback func(Builder)) {
    callback(*newBuilder(builder.stateMachine, event))
}

func (builder Builder) Transition(from, to Status, action func ([]interface{}) bool) {
    target := statusTarget{to, action}
    if transitions, found := builder.stateMachine.eventTransitions[builder.event]; found {
        transitions[from] = target
    } else {
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
    status           Status
    eventTransitions transitionMap
    validator        func (...interface{}) bool
}

func (stateMachine StateMachine) Status() Status {
    return stateMachine.status
}

func (stateMachine *StateMachine) Trigger(event Event, args ...interface{}) bool {
    commit := true
    if transitions, found := stateMachine.eventTransitions[event]; found {
        if target, listed := transitions[stateMachine.status]; listed {
            if commit && target.action != nil {
                commit = target.action(args)
            }
            if commit && stateMachine.validator != nil {
                commit = stateMachine.validator(args...)
            }
            if commit { stateMachine.status = target.value }
            return commit
        }
    }
    return false
}

func (stateMachine *StateMachine) SetTriggerValidator(validator func (...interface{}) bool) {
    stateMachine.validator = validator
}

func NewStateMachine(status Status, callback func (Builder)) (stateMachine *StateMachine) {
    eventTransitions := make(transitionMap, 5)
    stateMachine = &StateMachine{status, eventTransitions, nil}
    builder := *newBuilder(stateMachine)
    callback(builder)
    return stateMachine
}
