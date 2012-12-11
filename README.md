# Go ATD #

Automated Task Distributor implemented in Go

## Requirements ##

  * launchpad.net/gocheck
  * github.com/sdegutis/fsm

## Architecture ##

The library is composed of the following modules:

  * event: local event BUS with BUS manager running in a go routine
  * model: data models and collections
  * store: memory store to keep track of all the objects and their state
  * dispatch: distributors running in go routines and dispatch algorithms

The memory store runs in a go routine and offers a public API for CRUD
and query operations. the memory store will allow for using a redis
store instead in the future.

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

Evenutually other entities will be added:

  * network API using zeromq
  * external event bus with a set of ATD public events
  * logging
