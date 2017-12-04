# daemon

## overview

* Daemon has responsibility for
  * Receive requests from front end.
  * Caches commands and reflect them to local document store for working without network.
  * Send chunk of commands to datastore for syncing.
  * Receive chunk of data to sync with datastore.

## sync with datastore

* Daemon has a document store to cache commands from front end, and store data from datastore.
* Periodically, daemon tries to perform syncing local data with datastore.
  * Send local commands, that are "not synced yet", to datastore.
  * Retrieve chunk of data. They are JSON formatted and represent "1 week of data", for example.
  * Daemon stores 10 weeks of data in local document store.

![sync.svg](./uml/sync.svg)
