
# protocol buffers

```proto
syntax = "proto3";

package hashira;

// hashira service definition.
service Hashira {
  rpc create   (CreateCommand)   returns (Result) {}
  rpc update   (UpdateCommand)   returns (Result) {}
  rpc delete   (DeleteCommand)   returns (Result) {}
  rpc retrieve (RetrieveCommand) returns (Result) {}
}

message CreateCommand {
}

message UpdateCommand {
}

message DeleteCommand {
}

message RetrieveCommand {
}
```
