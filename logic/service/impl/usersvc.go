package impl

import (
	"com.dk.user/logic/service"
	"context"
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/hashicorp/go-uuid"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type userLogicSvc struct {
	logger log.Logger
}

func NewUserLogicService(logger log.Logger) service.UserLogicService  {
	return &userLogicSvc{
		logger: logger,
	}
}

func (u *userLogicSvc) Register(ctx context.Context, request *service.RegisterUsrRequest) (*service.RegisterUsrResponse, error) {
	u.logger.Log("file", "usersvc.go", "function", "Register", "action", "check password")
	rs := time.Now()
	defer func() {
		u.logger.Log("file", "usersvc.go", "function", "Register", "defer", time.Since(rs))
	}()
	if request.Pwd != request.PwdAgain {
		return nil, service.ErrorParaPassword
	}

	u.logger.Log("file", "usersvc.go", "function", "Register", "action", "generate password")
	s := time.Now()
	//pwd, err := bcrypt.GenerateFromPassword([]byte(request.Pwd), bcrypt.DefaultCost)
	pwd, err := bcrypt.GenerateFromPassword([]byte(request.Pwd), 4)
	e := time.Since(s)
	u.logger.Log("file", "usersvc.go", "function", "Register", "action", "generate password", "lost", e)
	if err != nil {
		return nil, service.ErrorServerInternal
	}

	request.Pwd = string(pwd)
	request.PwdAgain = string(pwd)

	u.logger.Log("file", "usersvc.go", "function", "Register", "action", "generate user id")
	uid, err := uuid.GenerateUUID()
	if err != nil {
		fmt.Println("generate uid error")
		return nil, service.ErrorServerInternal
	}
	request.Uid = uid

	u.logger.Log("file", "usersvc.go", "function", "Register", "action", "return")
	return nil, nil
}
