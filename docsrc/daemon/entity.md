# Data entities for local document store

* Two data entities are declared for local document store.
  * Tasks
  * Priorities

* Refer [proto file](https://github.com/pankona/hashira/blob/master/proto/hashira.proto) for latest data entity declaration.

## Tasks

* Tasks represents each ToDo items.

```proto
message Task {
    string id        = 1;
    string name      = 2;
    Place  place     = 3;
    bool   isDeleted = 4;
}
```

## Priorities

* Priorities represents priority of each task.
  * Priorities are represented as its place and array of task ID.
  * The ID placed on lower index means higher priority.

```proto
message Priority {
    Place place         = 1;
    repeated string ids = 2;
}
```
