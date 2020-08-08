package main

import (
	"context"
	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	LoginEndpoint endpoint.Endpoint
	LogoutEndpoint endpoint.Endpoint
}

func NewEps(svc LoginSvc) *Endpoints {
	return &Endpoints{
		LoginEndpoint: func(ctx context.Context, request interface{}) (response interface{}, err error) {
			req := request.(*LoginRequest)
			token, code, err  := svc.Login(ctx, req.Name, req.Pwd)
			if err != nil {
				return &LoginResponse{
					Token: "",
					Code:  code,
					Err:   err.Error(),
					Msg:   "kitlogin failed",
				}, nil
			}

			return &LoginResponse{
				Token: token,
				Code:  200,
				Err:   "",
				Msg:   "kitlogin ok",
			}, nil
		},

		LogoutEndpoint: func(ctx context.Context, request interface{}) (response interface{}, err error) {
			req := request.(*LogoutRequest)
			msg := svc.Logout(ctx, req.Name)
			return &LogoutResponse{
				Err: "",
				Msg: msg,
			}, nil
		},
	}
}

type LoginRequest struct {
	Name string `json:"name"`
	Pwd string `json:"pwd"`
}

type LoginResponse struct {
	Token string `json:"token"`
	Code int `json:"code"`
	Err string `json:"err"`
	Msg string `json:"msg"`
}

type LogoutRequest struct {
	Name string `json:"name"`
}

type LogoutResponse struct {
	Code int `json:"code"`
	Err string `json:"err"`
	Msg string `json:"msg"`
}
