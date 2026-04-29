# context/example

## Run

```bash
GOCACHE=/tmp/gocache go run ./context/example
```

## Request

```bash
curl -i localhost:8080/posts -H 'X-Trace-Id: REQ-1234'
```

## Request (success/timeout demo)

```bash
# 성공: delay가 timeout보다 짧음
curl -i 'localhost:8080/posts?delay_ms=100&timeout_ms=800' -H 'X-Trace-Id: REQ-1234'

# 타임아웃: delay가 timeout보다 김
curl -i 'localhost:8080/posts?delay_ms=3000&timeout_ms=800' -H 'X-Trace-Id: REQ-1234'
```

## Notes

- `X-Trace-Id`는 요청 스코프 메타데이터 예시로 `context.Value`에 담아 레이어(handler → service → repo)로 전파합니다.
- `timeout_ms`는 "요청 단위 타임아웃" 데모이고, `ctx.Done()`으로 repo가 빨리 중단되는 흐름을 로그로 확인할 수 있습니다.
