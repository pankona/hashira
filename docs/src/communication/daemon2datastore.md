# Communication between daemon to datastore

## Overview

- Use GRPC for communication between daemon and datastore
  - Daemon send array of command, and datastore will apply them
  - Daemon retrieves updates from datastore periodically for syncing

## GRPC APIs

- Use Hashira service. See [daemon overview](../daemon/overview.md).
