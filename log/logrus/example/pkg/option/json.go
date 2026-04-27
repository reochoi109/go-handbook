package option

import (
	"os"

	"github.com/sirupsen/logrus"
)

func JsonFormatFn() {
	log := logrus.New()

	log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05.000", // 시간 포맷 정의
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "timestamp", // JSON 내 키 이름 변경 가능 (예: time -> timestamp)
			logrus.FieldKeyLevel: "level",
			logrus.FieldKeyMsg:   "message",
			logrus.FieldKeyFunc:  "caller",
		},

		PrettyPrint: true, // 로그 정렬 여부
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

	entry.Info("JSONFormatter 설정 테스트 중입니다.")
	entry.Warn("시간 포맷과 구조화된 필드가 적용되었습니다.")
	entry.Error("에러 발생 시 파일명과 라인 번호가 'caller' 필드에 담깁니다.")
}
