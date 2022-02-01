# Client and server for a process to request attach by gdlv.

These two packages allow a program to request that a debugger attach to it.
The motivating use case is debugging a compiler, because compilers are often run by other processes in ways that are tedious to capture and reproduce.
Adding a debugger-invocation when things go wrong, or just before the problematic file and line nmumber are processed, may help in understanding a bug and its cause.

## Debug server

* listen on a port
* receive a pid:binaryname request on a port, reply "1\n" (yes) or "0\n" (no) to signal intent to debug.
* start gdlv on process.

## Client of debug server

* if environment variable "DEBUG_SERVER" is defined, expect it to be a localhost port to connect to.
* assuming it is defined, open a connection, send "*pid*:*binary*\n", wait for response.
* if response is "1\n", sleep 60 seconds waiting for server to start gdlv and attach
