package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"
)

func main() {
	// durationcheck: time.Second * 60 보다는 time.Minute 처럼 쓰는 게 안전/명확함
	d := time.Second * 60
	_ = d

	// staticcheck: time.Since(time.Now()) 는 거의 항상 의미 없는 코드(항상 거의 0)
	_ = time.Since(time.Now())

	// noctx: http 요청을 만들 때 context를 붙이지 않으면 취소/타임아웃 전파가 안 됨
	req, _ := http.NewRequest(http.MethodGet, "https://example.com", nil)
	_ = req

	// bodyclose: http.Get() 결과의 Body는 꼭 닫아야 함(리소스 누수)
	resp, _ := http.Get("https://example.com")
	_ = resp

	// errorlint: error 비교는 == 대신 errors.Is 를 써서 래핑된 에러도 잡기
	err := context.DeadlineExceeded
	if err == context.DeadlineExceeded {
		fmt.Println("deadline")
	}

	// ineffassign: 값을 할당해놓고 사용하기 전에 바로 덮어씀(무의미한 할당)
	x := 1
	x = 2
	fmt.Println(x)

	// exportloopref: 루프 변수를 고루틴에서 캡처하면 의도치 않게 같은 값이 찍힐 수 있음
	for _, v := range []int{1, 2, 3} {
		value := v
		go func() {
			fmt.Println(value)
		}()
	}

	// errorlint: errors.New() 로 만든 에러는 매번 다른 값이라 직접 비교(==)가 의미 없음
	if errors.New("x") == errors.New("x") {
		fmt.Println("절대 실행되지 않음")
	}
}
