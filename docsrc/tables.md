
# Database tables

## User

| column    | type   | remarks                       |
|-----------|--------|-------------------------------|
| user_id   | number | primary key, non-null, unique |
| user_name | string | non-null, unique              |
| password  | string | non-null                      |

## Task

| column    | type   | remarks                       |
|-----------|--------|-------------------------------|
| task_id   | number | primary key, non-null, unique |
| task_name | string | non-null                      |
| user_id   | number | non-null                      | 
| status_id | number | non-null                      | 
| label_id  | number | non-null                      | 
| done_at   | time   | non-null                      | 

* task represents tasks and its status.

## Status

| column      | type   | remarks                       |
|-------------|--------|-------------------------------|
| status_id   | number | primary key, non-null, unique |
| status_name | string | non-null, unique              |

* status represents task's status. like "Backlog", "ToDo", "Doing" and "Done" will be inserted here.

## Consume

| column      | type   | remarks                       |
|-------------|--------|-------------------------------|
| consume_id  | number | primary key, non-null, unique |
| task_id     | number | non-null                      |
| started_at  | time   | non-null                      | 
| finished_at | time   | non-null                      | 
| consumed    | time   | non-null                      | 

* consume represents how many times are consumed for specified task.

## Label

| column     | type   | remarks                       |
|------------|--------|-------------------------------|
| label_id   | number | primary key, non-null, unique |
| label_name | string | non-null                      |

* label represents task's label. this is used to task classification.

## ER diagram

```uml
@startuml

entity task {
    + task_id
    --
    task_name
    status_id
    label_id
    done_at
}

entity status {
    + status_id
    --
    status_name
}

entity consume {
  + consume_id
  --
  task_id
  started_at
  finished_at
  consumed
}

entity label {
    + label_id
    --
    label_name
}

task }---- status
task }---- label
task -ri-{ consume

@enduml
```
