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
