package user

import (
	"context"

	"github.com/sirupsen/logrus"
)

type UserService struct {
	log *logrus.Entry
}

func NewUserService(baseLogger *logrus.Logger) *UserService {
	return &UserService{
		log: baseLogger.WithField("domain", "user"),
	}
}

func (s *UserService) SignUpWithContext(ctx context.Context, username string) {
	traceID, _ := ctx.Value("trace_id").(string)

	s.log.WithFields(logrus.Fields{
		"username": username,
		"trace_id": traceID,
	}).Info("User sign-up process started")
}

func (s *UserService) SignUp(username string) {
	s.log.WithField("username", username).Info("User sign-up process started")
}
