package user

import "github.com/sirupsen/logrus"

type UserService struct {
	log *logrus.Entry
}

// 생성자 주입
func NewUserService(baseLogger *logrus.Logger) *UserService {
	return &UserService{
		log: baseLogger.WithField("domain", "user"),
	}
}

func (s *UserService) SignUp(username string) {
	s.log.WithField("username", username).Info("회원 가입 프로세스 시작")
}
