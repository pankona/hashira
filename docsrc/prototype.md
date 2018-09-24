# Prototype

Here's design memorandum of hashira prototyping, for checking usability.

## Feature notice of prototype version

* No performance/resource consideration
* front end is for linux PC only
  * CUI and GUI (may use astilectron)
* No cloud syncing
* No "Doing" time calculation
* Daemon uses MySQL as DB

## Sequence of tweaking tasks

```uml
@startuml

hide footbox

actor user           as user
participant "cui"    as cui
participant "daemon" as d
participant "DB"     as db

== initialize daemon ==

user -> d : launch
d -> db : create config and DB if not exist
d <- db : ok
d -> d : start gRPC server
user <- d : ok

== add a new task ==

cui -> d : command(new)\n via gRPC
note right of d : with task's\n name\n status\n label
d -> db : execute insert query
d <- db : ok
cui <- d : ok

== retrieve list of tasks (ordinal use) ==

cui -> d : command(list)\n via gRPC
note right of d : with\n Backlog\n ToDo\n Doing\n Today's Done
d -> db : execute select query
d <- db : ok (with task list)
cui <- d : ok (with task list)
cui -> cui : render using\n received tasks

== retrieve list of tasks (long history) ==

cui -> d : command(list)\n via gRPC
note right of d : with\n Done of this week\n or this month\n or last 2 months
d -> db : execute select query
d <- db : ok (with task list)
cui <- d : ok (with task list)
cui -> cui : render using\n received tasks

== move a task to Doing ==

cui -> d : command(update)\n via gRPC
note right of d : with\n task id\n new status (Doing)
d -> db : execute update query
note right of db: execute\n move current Doing to ToDo\n move specified task to Doing
d <- db : ok (with updated tasks)
cui <- d : ok (with updated tasks)
cui -> d : retrieve list of tasks
cui <- d : list of tasks
cui -> cui : render using\n received tasks

== move a task to Done ==

cui -> d : command(update)\n via gRPC
note right of d : with\n task id\n new status (Done)
d -> db : execute update query
note right of db: execute\n move current Doing to Done
d <- db : ok (with updated task)
cui <- d : ok (with updated task)
cui -> d : retrieve list of tasks
cui <- d : list of tasks
cui -> cui : render using\n received tasks

@enduml
```
