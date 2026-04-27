package main

import (
	"flag"
	"fmt"
)

// AppConfig는 애플리케이션의 모든 설정을 담습니다.
type AppConfig struct {
	Port  int
	Env   string
	Debug bool
}

func main() {
	// 1. flag 정의
	// flag.X(이름, 기본값, 설명)
	port := flag.Int("port", 8080, "서버 포트 번호")
	env := flag.String("env", "development", "실행 환경 (development, production)")
	debug := flag.Bool("debug", false, "디버그 모드 활성화 여부")

	// 2. flag 파싱
	flag.Parse()

	// 3. 구조체에 매핑
	cfg := AppConfig{
		Port:  *port,
		Env:   *env,
		Debug: *debug,
	}

	// 4. 비즈니스 로직 적용
	if cfg.Debug {
		fmt.Println("디버그 모드가 활성화되었습니다.")
	}

	fmt.Printf("서버 시작: 환경[%s], 포트[%d]\n", cfg.Env, cfg.Port)
}
