package proxymiddleware

import (
	"context"
	"github.com/dingkegithub/com.dk.user/das/proto/userpb"
	"github.com/dingkegithub/com.dk.user/logic/common"
	"github.com/dingkegithub/com.dk.user/logic/service"
	"github.com/dingkegithub/com.dk.user/sidecar/discovery"
	"github.com/dingkegithub/com.dk.user/sidecar/discovery/kitnacos"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/lb"
	"golang.org/x/crypto/bcrypt"
	"time"
)

//
// 用户服务代理，在服务至上添加额外功能
//
func UserLogicProxy(ctx context.Context, logger log.Logger, cli discovery.RegisterCenterClient) service.UserLogicServiceMiddleware {
	var (
		maxAttempts = 3 // 失败尝试次数
		maxTime     = 250 * time.Millisecond
	)

	// 正常情况下 注册服务有异常，应当直接退出，服务不可用
	if cli == nil {
		return func(svc service.UserLogicService) service.UserLogicService {
			return svc
		}
	}

	// 服务发现
	instance, err := kitnacos.NewInstancer(
		cli,
		&discovery.ServiceMeta{
			Group:   "default",
			Cluster: "default",
			SvcName: common.ServiceUsrDasSrv,
		},
		logger)

	if err != nil {
		logger.Log("file", "userproxy.go", "error", err)
		panic(err)
	}

	// 可用服务站点
	registerEps := sd.NewEndpointer(instance, NewRegisterFactory, logger)
	loginEps := sd.NewEndpointer(instance, NewSvcLoginFactory, logger)

	// 将站点加入到负载均衡器
	registerBalancer := lb.NewRoundRobin(registerEps)
	registerRetry := lb.Retry(maxAttempts, maxTime, registerBalancer)

	loginBalance := lb.NewRoundRobin(loginEps)
	loginRetry := lb.Retry(maxAttempts, maxTime, loginBalance)

	// 返回含有服务发现，负载均衡的service
	return func(svc service.UserLogicService) service.UserLogicService {
		return NewUsrLogicProxy(ctx, svc, logger, registerRetry, loginRetry)
	}
}

type UsrLogicProxy struct {
	logger              log.Logger
	ctx                 context.Context
	nxt                 service.UserLogicService
	registerRpcEndpoint endpoint.Endpoint
	loginRpcEndpoint    endpoint.Endpoint
}

func NewUsrLogicProxy(ctx context.Context, svc service.UserLogicService, logger log.Logger, registerEps endpoint.Endpoint, loginEps endpoint.Endpoint) service.UserLogicService {
	return &UsrLogicProxy{
		logger:              logger,
		ctx:                 ctx,
		nxt:                 svc,
		registerRpcEndpoint: registerEps,
		loginRpcEndpoint:    loginEps,
	}
}

//
// 用户注册接口
// Endpoint 调用
//
func (u *UsrLogicProxy) Register(ctx context.Context, request *service.RegisterUsrRequest) (*service.Response, error) {

	rs := time.Now()
	_, err := u.nxt.Register(ctx, request)
	re := time.Since(rs)
	if err != nil {
		u.logger.Log("file", "userproxy.go",
			"function", "Register", "action", "invoke next Register", "error", err)
		return nil, err
	}

	u.logger.Log("file", "userproxy.go", "function", "Register", "action", "rpc invoke", "lost", re)
	// 调用 rpc 接口, 这里有loadbalance 包裹
	start := time.Now()
	resp, err := u.registerRpcEndpoint(ctx, request)
	end := time.Since(start)
	u.logger.Log("file", "userproxy.go", "function", "Register", "action", "rpc invoke", "lost", end)
	if err != nil {
		u.logger.Log("file", "userproxy.go", "function", "Register", "action", "rpc invoke", "error", err)
		return nil, err
	}

	return resp.(*service.Response), nil
}

func (u *UsrLogicProxy) Login(ctx context.Context, request *service.LoginRequest) (*service.Response, error) {
	// 调用 rpc 接口, 这里有loadbalance 包裹
	start := time.Now()
	resp, err := u.loginRpcEndpoint(ctx, request)
	end := time.Since(start)
	u.logger.Log("file", "userproxy.go", "func", "Login", "action", "rpc invoke", "lost", end)

	if err != nil {
		u.logger.Log("file", "userproxy.go", "function", "Login", "action", "rpc invoke", "error", err)
		return nil, err
	}

	response := resp.(*service.Response)

	data := response.Data.([]*userpb.UserData)
	if len(data) > 1 {
		return &service.Response{
			Error: service.ErrorCodeParaPassword,
			Msg:   service.ErrMsg(service.ErrorCodeParaPassword).Error(),
			Data:  nil,
		}, nil
	}

	if len(data) <= 0 {
		return &service.Response{
			Error: service.ErrorCodeParaPassword,
			Msg:   service.ErrMsg(service.ErrorCodeParaPassword).Error(),
			Data:  nil,
		}, nil
	}

	userInfo := data[0]
	pwd, err := bcrypt.GenerateFromPassword([]byte(request.Pwd), 4)
	if err != nil {
		return &service.Response{
			Error: service.ErrorCodeServerInternal,
			Msg:   service.ErrMsg(service.ErrorCodeServerInternal).Error(),
			Data:  nil,
		}, nil
	}

	if string(pwd) != userInfo.Pwd {
		return &service.Response{
			Error: service.ErrorCodeParaPassword,
			Msg:   service.ErrMsg(service.ErrorCodeParaPassword).Error(),
			Data:  nil,
		}, nil
	}

	token := ""

	return &service.Response{
		Error: 0,
		Msg:   "ok",
		Data:  service.LoginToken{
			Uid:   userInfo.Uid,
			Token: token,
		},
	}, nil
}
