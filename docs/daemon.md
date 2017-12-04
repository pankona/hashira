# daemon

## overview

* daemon has responsibility for
  * receive requests from front end
  * caches commands and reflect them to local database for working without network
  * send chunk of commands to datastore for syncing
  * receive chunk of commands to sync with datastore 
  * notify updates to front end

## cache mechanism

* Daemon has document store to cache commands from front end.
* Daemon has a database to store current status of tasks.
  * Commands are used to apply modification to database.
  * They are also used to apply modification to datastore.

## sync with datastore

* Periodically, daemon tries to perform syncing local data with datastore. 
  * Send local commands, that are "not synced yet", to datastore.
  * Retrieve commands from datastore that are "not applied to local database yet".
  * Apply all retrieved commands to local database.

![sync.svg](./uml/sync.svg)
