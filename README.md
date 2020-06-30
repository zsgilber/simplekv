# simplekv
A simple in-memory key-value store. 

# Use

The server implements the RESP protocol, so redis clients can be used to send commands to the server:

```console
$ redis-cli -h localhost -p 3003 PING
PONG
```
