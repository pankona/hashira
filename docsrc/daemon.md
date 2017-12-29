# daemon

## overview

* Daemon has responsibility for
  * Receive requests from front end.
  * Caches commands and reflect them to local document store for working without network.
  * Send chunk of commands to datastore for syncing.
  * Receive chunk of data to sync with datastore.

## sync with datastore

* Daemon has a document store to cache commands from front end, and store data from datastore.
* Periodically, daemon tries to perform syncing local data with datastore.
  * Send local commands, that are "not synced yet", to datastore.
  * Retrieve chunk of data. They are JSON formatted and represent "1 week of data", for example.
  * Daemon stores 10 weeks of data in local document store.

```uml
@startuml

hide footbox

participant "client 1"  as c1
participant "daemon 1"  as d1
participant "datastore" as ds
participant "daemon 2"  as d2
participant "client 2"  as c2

note right of d1 : assume daemon 1 is offline
c1 -> d1 : command-1 (cached)
c1 -> d1 : command-2 (cached)

note right of d1 : assume daemon 1 is online,\n then send all cached commands to datastore
d1 -> ds : command-1, 2 as JSON

ds -> ds : apply command-1, 2

d1 -> ds : fetch(10 weeks)
d1 <- ds : chunk of data as JSON\n 10 weeks of data for example
d1 -> d1 : store JSON to local document store

note right of d2 : assume daemon 2 is offline
d2 <- c2 : command-3 (cached)
d2 <- c2 : command-4 (cached)

note right of d2 : assume daemon 2 is online,\n then send all cached commands to datastore
ds <- d2 : command-3, 4 as JSON

ds -> ds : apply command-3, 4

ds <- d2 : fetch(10 weeks)
ds -> d2 : chunk of data as JSON
d2 -> d2 : store JSON to local document store

d1 -> ds : fetch(10 weeks)
d1 <- ds : chunk o data as JSON
d1 -> d1 : store JSON to local document store

@enduml
```

## gRPC API

* Hashira Service
  * send(array of command)
    * sends specified commands. 
  * retrieve(from, to (number of weeks)) array of task
    * returns array of task with specified term.

* command and related enumeration

what (enum)

| enum   | remarks                |
|--------|------------------------|
| new    | create a new task      |
| update | update state of a task |

command (structure)

| field   | type         | remarks                                  |
|---------|--------------|------------------------------------------|
| what    | enum of what | new, update, etc.                        |
| payload | string       | JSON formatted string how to treat what. |

## command handling

* When daemon receives commands, cache them and return ok immediately.
* If daemon is online, send cached commands to datastore.
* When datastore receives commands and succeed to apply them, datastore sends notification to daemon.
* At daemon receiving notification, retrieve chunk of data from datastore for syncing.
* When daemon succeed to apply them, daemon sends notification to front end.
* At front end receiving notification, retrieve chunk of data from daemon.
* When front end succeed to retrieve chunk of data, render them.

* If daemon is offline, postpone to send commands to datastore.
  * Instead, daemon applies the cached commands to local document store,
  and send notification to front end as same as written above.

```uml
@startuml

hide footbox

participant client    as c
participant daemon    as d
participant datastore as ds

c -> d : send(command-1)
d -> d : cache command-1
c <- d : ok

alt daemon is online

  note right of d : proceed on goroutine
  d -> ds : send(command-1)
  ds -> ds : apply(command-1)
  d <- ds : ok
  
  note right of ds: send notification\n because of update
  d <- ds : notify
  
  note right of d : fetch latest data\n from datastore
  d -> ds  : fetch(10 weeks)
  d <- ds  : chunk of data
  d -> d   : apply data to\n local document store

else daemon is offline

  d -> d : apply data to\n local document store\n (using cache)

end

note right of d : send notification\n because of update
c <- d   : notify
c -> d   : fetch(10 weeks)
c <- d   : chunk of data
c -> c   : render data

@enduml
```
