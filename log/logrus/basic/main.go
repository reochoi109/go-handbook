package main

import (
	"github.com/sirupsen/logrus"
)

func main() {
	// 로그 레벨을 Trace로 설정해야 모든 로그가 출력됨
	logrus.SetLevel(logrus.TraceLevel)

	// 1. Trace: 상세 실행 경로
	logrus.Trace("데이터베이스 연결 설정 정보 확인중...")

	// 2. Debug: 개발 시 진단 정보
	logrus.Debug("API 요청 헤더: Content-Type=application/json")

	// 3. Info: 시스템 주요 흐름
	logrus.Info("주문 서비스가 정상적으로 시작되었습니다.")

	// 4. Warn: 비정상 징후 (당장 죽진 않음)
	logrus.Warn("외부 결제 서비스 응답이 2초 이상 지연되었습니다.")

	// 5. Error: 특정 작업 실패
	logrus.Error("주문 데이터 저장에 실패했습니다. (DB Timeout)")

	// 6. Fatal: 시스템 종료 직전
	// logrus.Fatal("설정 파일이 없습니다. 시스템을 종료합니다.")

	// 7. Panic: 복구 불가능한 상태
	// logrus.Panic("스택 오버플로우 발생! 시스템 복구 불가능.")
}
