package endpoints

import (
	"context"
	"fmt"
	"github.com/dingkegithub/com.dk.user/logic/service"
	"github.com/go-kit/kit/endpoint"
	"time"
)

//
// 统一管理站点
//
type UserEndpoints struct {
	RegisterEndpoint endpoint.Endpoint
	LoginEndpoint endpoint.Endpoint
}

//
// @param svc 占掉调用的服务
//
func NewUserLogicEndpoints(svc service.UserLogicService) *UserEndpoints {
	return &UserEndpoints{
		RegisterEndpoint:MakeRegisterEndpoint(svc),
		LoginEndpoint: MakeLoginEndpoint(svc),
	}
}

//
// 每个接口提供一个 Endpoint 工transport调用
//
func MakeRegisterEndpoint(svc service.UserLogicService) endpoint.Endpoint  {
	// 方法实际有 transport 调用
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*service.RegisterUsrRequest)

		// Endpoint 调用实际service的接口，这里首先调用代理
		start := time.Now()
		resp, err := svc.Register(ctx, req)
		lost := time.Since(start)
		fmt.Println("file", "userendpoints.go", "function", "MakeRegisterEndpoint", "lost", lost)

		if err != nil {
			return nil, err
		}

		return resp, nil
	}
}

func MakeLoginEndpoint(srv service.UserLogicService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*service.LoginRequest)

		start := time.Now()
		resp, err := srv.Login(ctx, req)
		lost := time.Since(start)

		fmt.Println("file", "userendpoints.go",
			"func", "MakeLoginEndpoint",
			"msg", "login request",
			"lost", lost)
		if err != nil {
			return nil, err
		}

		return resp, err
	}
}