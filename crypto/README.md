# crypto

표준 라이브러리 기반으로 `rand` / `hash` / `HMAC` 기본을 다룹니다.

## Run

```bash
GOCACHE=/tmp/gocache go run ./crypto/basic
```

```bash
GOCACHE=/tmp/gocache go run ./crypto/example
```

## Notes

- 보안 목적의 랜덤은 `crypto/rand`를 사용(`math/rand` 아님).
- 해시는 “비밀번호 저장용”이 아닙니다(비밀번호는 `bcrypt/scrypt/argon2` 같은 KDF가 필요).
