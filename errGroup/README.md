# errGroup

`golang.org/x/sync/errgroup`는 여러 작업을 병렬로 실행하고,

- 하나라도 실패하면 `Wait()`에서 에러를 반환하고
- `WithContext()`를 쓰면 실패 시 `Context` cancel이 전파되어 나머지 작업을 조기에 중단할 수 있습니다.

## Run

```bash
GOCACHE=/tmp/gocache go run ./errGroup/basic/f1
```

```bash
GOCACHE=/tmp/gocache go run ./errGroup/basic/f2
```

```bash
GOCACHE=/tmp/gocache go run ./errGroup/example
```

## Demo (cancel propagation)

```bash
FAIL_ORDERS=1 GOCACHE=/tmp/gocache go run ./errGroup/example
```
