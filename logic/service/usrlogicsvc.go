package service

import "context"

type UserLogicService interface {
	Register(ctx context.Context, request *RegisterUsrRequest) (*RegisterUsrResponse, error)

	//Login()
	//
	//Logout()
	//
	//ResetPwd()
	//
	//UpdateInfo()
}

type UserLogicServiceMiddleware func(service UserLogicService) UserLogicService