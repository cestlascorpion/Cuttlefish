# Cuttlefish

OnlineSvr: Fix your servers nintendo or whatever.

```txt
message Tentacle {
  uint32 uid = 1; // cuttlefish's id
  string key = 2; // which defines the tentacle
  map<string, string> val = 3; // info of the tentacle
}

A cuttle has one or more tentacle(which means connection) of course.
```
