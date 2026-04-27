package option

import (
	"os"

	"github.com/sirupsen/logrus"
)

func TextFormatFn() {
	log := logrus.New()

	log.SetFormatter(&logrus.TextFormatter{
		ForceColors:            true,                      // 로그 색상
		ForceQuote:             true,                      // 터미널에서 색상 지원 하지 않아도 강제 적용
		TimestampFormat:        "2006-01-02 15:04:05.000", // 타임 포맷 형식
		FullTimestamp:          true,                      // 로그에 타임 스탬프 표시
		DisableLevelTruncation: false,                     // 로그 레벨 대문자 표기 여부
		DisableSorting:         false,                     // 필드 순서 정렬
	})

	// 출력 및 레벨 설정
	log.SetOutput(os.Stdout)
	log.SetLevel(logrus.TraceLevel)

	// 파일/라인 번호 출력 옵션
	log.SetReportCaller(true)

	// 테스트 로그
	entry := log.WithFields(logrus.Fields{
		"env":     "development",
		"service": "api-gateway",
	})

	entry.Info("TextFormatter 설정 테스트 중입니다.")
	entry.Warn("시간 포맷과 색상이 적용되었는지 확인하세요.")
	entry.Error("에러 발생 시 파일명과 라인 번호가 표시됩니다.")
}
