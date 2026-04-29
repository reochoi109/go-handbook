# errors

## Run

```bash
GOCACHE=/tmp/gocache go run ./errors/basic
```

```bash
GOCACHE=/tmp/gocache go run ./errors/example
```

```bash
curl -i localhost:8080/users/me -H 'X-User-Id: user_123'
curl -i localhost:8080/users/me -H 'X-User-Id: user_999'
curl -i localhost:8080/users/me
```
