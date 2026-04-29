# context/example

## Run

```bash
GOCACHE=/tmp/gocache go run ./context/example
```

## Request

```bash
curl -i localhost:8080/posts -H 'X-Trace-Id: REQ-1234'
```
