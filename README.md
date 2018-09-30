# hashira

Application to manage **today's** tasks. Establish a Hashira for the day.  
design documents are available at [here](https://pankona.github.io/hashira/)

## Concepts

* Makes tasks clear for the day and concentrate to eliminate them.
* Records time consuming for each tasks to reveal differences between estimations and results.
  * May help increasing accuracy of our work-load estimation.
  * May help us to notice our waste of time.

## Features

### Manage today's tasks

* "Backlogs" to add any miscellaneous tasks.
* Move tasks for today to "ToDo".
* Starts a task, then the task moves to "Doing".
* Finishes the task, then the task moves to "Done".
* If the "Doing" task is interrupted by another unexpected task, then the "Doing" task moves back to "ToDo" and new task is placed on "Doing".
* At the end of the day, **(Not implemented yet)**
  * "Done" field is archived automatically and new one is created for new day.
  * "Doing" task moves back to "ToDo" automatically to intent to start again next day.

### Calculate consumed times **(Not implemented yet)**

* Calculates how many times is consumed for each "Doing" task.
* Consumed time is measured only for one task, which is placed on top of "Doing".
* Show them in graph. 

## Available Keybindings of hashira-cui

* Ordinal use

| Key   | Action                                                            | Remarks                                           |
|-------|-------------------------------------------------------------------|---------------------------------------------------|
| Enter | Show input window for register a new task                         |                                                   |
| e     | Show input window for editing focused task                        |                                                   |
| Space | Select focused task                                               |                                                   |
| j / k | Up/Down cursor<br>(change priority if a task is selected)         |                                                   |
| h / l | Change focused pane<br>(change task's pane if a task is selected) |                                                   |
| i / I | Move focused task to left/right pane                              |                                                   |
| x     | Move focused task to Done                                         | If focused task is already on Done, it is deleted |

* While input

| Key      | Action      | Remarks                 |
| Ctrl-h/l | Move cursor | Same as using arrow key |


## Notes

* hashira generates its configuration file under `$HOME/.config/hashira`
  * Remove them for re-initializing hashira or leaving from using hashira...

## Limitation

### hashira-cui

* Multi-byte characters are not supported

## LICENSE

MIT

## Author

Yosuke Akatsuka (@pankona)
* [Twitter](https://twitter.com/pankona)
* [GitHub](https://github.com/pankona)
* [Qiita](https://qiita.com/pankona)
