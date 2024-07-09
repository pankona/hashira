# Communication between front to daemon

## Overview

- hashira uses GRPC for communication between front and daemon

## PC

### CLI and Daemon

- CLI to Daemon
  - Add a new task to Backlog
  - Change task's status
  - Show task list
- Daemon to CLI
  - None

### GUI and Daemon

- Assume to use Electron
- GUI to daemon
  - Add a new task to Backlog, ToDo, Doing, Done
  - Change task's status
  - Show task list on each status
  - Show consume of each task
- Daemon to GUI
  - notify any update of tasks

- TODO: Show GUI picture

## Android

### application and daemon

- GUI to Daemon
  - Add a new task to Backlog, ToDo, Doing, Done
  - Change task's status
  - Show task list on each status
  - Show consume of each task
- Daemon to GUI
  - Notify any update of tasks

- TODO: Show GUI picture

### Widget and Daemon

- Widget to Daemon
  - Add a new task to Backlog
  - Change task's status to Done
  - Show task list
- Daemon to widget
  - Notify any update of tasks

- TODO: Show GUI picture
