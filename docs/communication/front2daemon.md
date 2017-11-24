# communication between front to daemon

## overview

* use gRPC for communication between front and daemon

## PC

### cli and daemon

* cli to daemon
  * add a new task to Backlog
  * change task's status
  * show task list
* daemon to cli
  * none

### gui and daemon

* assume to use Electron
* gui to daemon
  * add a new task to Backlog, ToDo, Doing, Done
  * change task's status
  * show task list on each status
  * show consume of each task
* daemon to gui
  * notify any update of tasks

* TODO: show GUI picture

## Android

### application and daemon

* gui to daemon
  * add a new task to Backlog, ToDo, Doing, Done
  * change task's status
  * show task list on each status
  * show consume of each task
* daemon to gui
  * notify any update of tasks

* TODO: show GUI picture

### widget and daemon

* widget to daemon
  * add a new task to Backlog
  * change task's status to Done
  * show task list
* daemon to widget
  * notify any update of tasks

* TODO: show GUI picture

## gRPC APIs

* Hashira Service
  * new
  * changeStatus
  * list
  * subscribe


