# communication between daemon to datastore

## overview

* use gRPC for communication between daemon and datastore
  * daemon send chunk of transactions, and datastore will apply them
  * daemon retrieves updates from datastore periodically for syncing

## gRPC APIs

* HashiraStore service
  * update
    * send chunk of transactions and datastore will apply them
    * if some of transaction is older than "last update", older transactions will be discarded
  * diff
    * returns chunk of transactions they are applied since specified date/time

