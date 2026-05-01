# patterns/workerpool

## Run

```bash
GOCACHE=/tmp/gocache go run ./patterns/workerpool/basic
```

## Benchmark

```bash
GOCACHE=/tmp/gocache go test ./patterns/workerpool/basic -run '^$' -bench .
```

## Flow (visual)

```mermaid
flowchart LR
  PROD[Producer] --> J[jobs chan]
  J --> W1[Worker 1]
  J --> W2[Worker 2]
  J --> W3[Worker 3]
  J --> W4[Worker 4]

  W1 --> R[results chan]
  W2 --> R
  W3 --> R
  W4 --> R
  R --> CONS[Consumer<br/>range results]

  W1 --> CXL[cancel ctx<br/>fail-fast]
  W2 --> CXL
  W3 --> CXL
  W4 --> CXL
  CXL -. stop .-> W1
  CXL -. stop .-> W2
  CXL -. stop .-> W3
  CXL -. stop .-> W4
```

## Notes

- worker pool은 "동시성 제한"과 "작업 큐"를 만들 때 자주 씁니다.
- fail-fast가 필요하면 cancel 전파를 설계하고, 결과/에러 채널 close 책임을 명확히 두세요.
