# hashira

Application to manage **today's** tasks. Establish a Hashira for the day.
design documents are available at [here](https://pankona.github.io/hashira/)

## concepts

* Makes tasks clear for the day and concentrate to eliminate them.
* Records time consuming for each tasks to reveal differences between estimations and results.
  * May help increasing accuracy of our work-load estimation.
  * May help us to notice our waste of time.

## features

### manage today's tasks

* "backlogs" to add any miscellaneous tasks.
* Move tasks for today to "ToDo".
* Starts a task, then the task moves to "Doing".
  * Only one task can be placed on "Doing".
* Finishes the task, then the task moves to "Done".
* If the "Doing" task is interrupted by another unexpected task, then the "Doing" task moves back to "ToDo" and new task is placed on "Doing".
* At the end of the day,
  * "Done" field is archived automatically and new one is created for new day.
  * "Doing" task moves back to "ToDo" automatically to intent to start again next day.

### calculate consumed times

* Calculates how many time is consumed for each "Doing" task.
* Show them in graph. 

## LICENSE

MIT
