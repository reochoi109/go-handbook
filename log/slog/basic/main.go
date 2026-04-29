package main

import (
	"context"
	"log/slog"
	"os"
)

/*
 slog는 내부적으로 로그의 중요도를 숫자로 관리합니다.
 숫자가 낮을수록 상세하고, 높을수록 치명적입니다.

 - LevelTrace (-8) : 가장 상세한 실행 흐름 추적 (Trace)
 - LevelDebug (-4) : 개발 단계의 상세 진단 정보 (Debug)
 - LevelInfo   (0) : 시스템의 주요 동작 정보 (Info)
 - LevelWarn   (4) : 주의가 필요한 비정상 징후 (Warn)
 - LevelError  (8) : 기능 장애 및 에러 상황 (Error)
 - LevelCritical (12) : 시스템 가동 불가능 상태 (Critical)

 설정한 Level 값보다 같거나 큰 로그만 출력됩니다.
*/

// 1. 커스텀 레벨 정의
const (
	LevelTrace = slog.Level(-8)
	LevelCritical = slog.Level(12)
)

func main() {
	opts := &slog.HandlerOptions{
		Level: LevelTrace,
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, opts))

	// 3. 6단계 로그 출력
	logger.Log(context.Background(), LevelTrace, "함수 진입: OrderService.Create()")
	logger.Debug("상세 개발 정보: DB 쿼리 실행 완료")
	logger.Info("시스템 정보: 서버가 8080 포트에서 시작됨")
	logger.Warn("잠재적 위험: 외부 API 응답이 2초 이상 지연됨")
	logger.Error("에러 발생: 주문 데이터 저장 실패")
	logger.Log(context.Background(), LevelCritical, "치명적 상태: 설정 파일(config.yaml)을 찾을 수 없음!")
}
