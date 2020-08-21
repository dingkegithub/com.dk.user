package service

import "context"

type UserLogicService interface {
	Register(ctx context.Context, request *RegisterUsrRequest) (*Response, error)

	Login(ctx context.Context, request *LoginRequest) (*Response, error)
	//
	//Logout()
	//
	//ResetPwd()
	//
	//UpdateInfo()
}

type UserLogicServiceMiddleware func(service UserLogicService) UserLogicService