# Usage:

This project is usable either by running the executable directly on two different shells:

```bash
    # Shell 1
    go run ./cmd/main.go -queue
    # Shell 2
    go run ./cmd/main.go -input in.txt -output out.txt
```

Or by executing docker for the queue service, and then running client locally:
```bash
    # Docker (on shell 1)
    docker build -t queue:local .
    docker run -ti --rm -p 8080:8080 queue:local
    # Client (on shell 2)
    go run ./cmd/main.go -url http://localhost:8080/ -input in.txt -output out.txt
```

Arguments allow different execution ports and URL's.
Docker container (and queue executable) honours `PORT` env variable.

## Comments on the design

### Queue

Because the task was not clear on the usage, I initially implemented two scaffold interfaces, one for
reactive queue (which can deliver messages by pushing them to a consumer) and one for static queue
(which must be polled for new messages).

In the trivial implementation, there is no way to have multiple consumers in the static queue case.
Changing that requires persisting named queue read markers, which may be problematic and requires ACID.

Reactive queue implementation requires a smarter client, in order to be useful. Due to time constraints
I did not explore http streaming, but that seems to be doable.

### Client

I've added a more interesting use case (interleaving) in the tests. I could not figure out any real world
applications for the client as described in the task.
