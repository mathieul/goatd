# Go ATD #

Automated Task Distributor implemented in Go

## Requirements ##

  * launchpad.net/gocheck
  * github.com/sdegutis/fsm

## Architecture ##

The library is composed of the following modules:

  * event: local event BUS - manager running in a go routine
  * model: data models and collections
  * store: memory store to keep track of all the objects and their state
  * dispatch: distributors running in go routines and dispatch algorithms

The memory store runs in a go routine and offers a public API for CRUD
and query operations. the memory store will allow for using a redis
store instead in the future.

The memory store holds the "model truth". Model values returned via the store
public API are just copies, used to hold a temporary state of those models
and to interact with the public API to request state changes.

Model values hold properties and attributes. Attributes can be changed,
properties can only be accessed, and are updated using business rules
implemented in the model itself.

There is at most one distributor go routine running for each team. Distribution
is executed serially for a team - i.e.: only one task assignment is executed
in parallel, so there is no conflict between different distributor instances
trying to access the same teammate or the same task.

It is possible to subscribe to one, many or all events for 1, many or all
teams at a time. The use case to listen to the events for all teams is to publish
events to an external BUS in the future, and publish statistics for other
applications to consume.

The main routine or a go routine spawned by a zeromq value can create, read,
update or delete model values such as teams, teammates, queues, skills and
tasks. They can also enqueue / dequeue tasks, sign-in / sign-out and request
changing teammates status. And finally the ATD offers and assigns tasks to
teammate for completion.

The minimum number of go routines is 4:

  * event bus manager
  * memory store
  * distributor
  * requester

Eventually other entities will be added:

  * network API using zeromq
  * external event bus with a set of ATD public events
  * logging
  * statistics
