# Client of debug server

* if environment variable "DEBUG_SERVER" is defined, expect it to be a localhost port to connect to.
* assuming it is defined, open a connection, send "*pid*:*binary*\n", wait for response.
* if response is "1\n", sleep 60 seconds waiting for server to start gdlv and attach
