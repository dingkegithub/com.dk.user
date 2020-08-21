package impl

import (
	"context"
	"github.com/dingkegithub/com.dk.user/logic/service"
	"github.com/go-kit/kit/log"
	"golang.org/x/crypto/bcrypt"
	"sync/atomic"
	"time"
)


var tmpIdGenerate uint64 = 0

type userLogicSvc struct {
	logger log.Logger
}

func NewUserLogicService(logger log.Logger) service.UserLogicService  {
	return &userLogicSvc{
		logger: logger,
	}
}

func (u *userLogicSvc) Login(ctx context.Context, request *service.LoginRequest) (*service.Response, error) {
	if request.Name == "" || request.Pwd == "" {
		return nil, service.ErrorParaPassword
	}

	return nil, nil
}

func (u *userLogicSvc) Register(ctx context.Context, request *service.RegisterUsrRequest) (*service.Response, error) {
	rs := time.Now()
	defer func() {
		u.logger.Log("file", "usersvc.go",
			"func", "Register",
			"msg", "register user",
			"lost", time.Since(rs))
	}()

	if request.Pwd != request.PwdAgain {
		return nil, service.ErrorParaPassword
	}

	s := time.Now()
	//pwd, err := bcrypt.GenerateFromPassword([]byte(request.Pwd), bcrypt.DefaultCost)
	pwd, err := bcrypt.GenerateFromPassword([]byte(request.Pwd), 4)
	e := time.Since(s)
	u.logger.Log("file", "usersvc.go",
		"func", "Register",
		"msg", "generate password",
		"lost", e)

	if err != nil {
		return nil, service.ErrorServerInternal
	}

	return &service.Response{
		Error: 0,
		Msg:   "ok",
		Data:  &service.RegisterUsrRequest{
			Uid:      atomic.AddUint64(&tmpIdGenerate, 1),
			Name:     request.Name,
			Pwd:      string(pwd),
			PwdAgain: string(pwd),
			Age:      0,
		},
	}, nil
}
