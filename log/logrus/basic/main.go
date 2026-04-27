package main

import (
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.Info("간단한 정보 로그")
	logrus.Warn("주의가 필요한 상황")
	logrus.Error("에러 발생!")

	logrus.WithFields(logrus.Fields{
		"action": "login",
		"user":   "admin",
	}).Info("로그인 시도")
}
