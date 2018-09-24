# communication between daemon to datastore

## overview

* Use gRPC for communication between daemon and datastore
  * Daemon send array of command, and datastore will apply them
  * Daemon retrieves updates from datastore periodically for syncing

## gRPC APIs

* Use Hashira service. See [daemon overview](../daemon/overview.md).

