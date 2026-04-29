# errors

## Run

```bash
GOCACHE=/tmp/gocache go run ./error/basic
```

```bash
GOCACHE=/tmp/gocache go run ./error/example
```

```bash
curl -i localhost:8080/users/me -H 'X-User-Id: user_123'
curl -i localhost:8080/users/me -H 'X-User-Id: user_999'
curl -i localhost:8080/users/me
```
