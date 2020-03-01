# Design notes for api

## Basic usecases

### Clients sync tasks by uploading tasks to server

- Client is assumed to have attribute "last sync date/time"
- Client's tasks has attribute to indicate "dirty or not"
  - Dirty means "the task is not synced since its last update"
- Client can extract "unsync tasks" by filtering by "dirty task"
- Client sends "unsync tasks" to server to sync their latest status
- Status of tasks that are synced correctly becomes "not dirty"

### Clients sync tasks by downloading tasks from server

- Client is assumed to have attribute "last sync date/time" (already written in previous section)
- Client downloads "unsync tasks on server" after uploading "unsync tasks" on client (written in previous section)

### When does syncing of tasks run?

- On start
- On terminate (if possible) 
- Periodically (like each 10 seconds)

## Components

### Task ball: Chunk of tasks to upload/download

- (Maybe) tasks to upload/download is formed as a large JSON object
  - Name it as "Task ball" in this document.

## Other notes

- If client fails to upsync, then client does NOT attempt to downsync not to loose local modifications
- If client fails to downsync, then client does NOT attempt to retry 
  - If there's some modifications on client's local, the downsync may overwrite client's local tasks unintentionally
  - Always "first upsync, then downsync"
