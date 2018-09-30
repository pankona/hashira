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

## LICENSE

MIT

## Author

Yosuke Akatsuka (@pankona)
* [Twitter](https://twitter.com/pankona)
* [GitHub](https://github.com/pankona)
* [Qiita](https://qiita.com/pankona)
