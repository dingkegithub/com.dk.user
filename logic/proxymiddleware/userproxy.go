package proxymiddleware

import (
	"context"
	"github.com/dingkegithub/com.dk.user/logic/common"
	"github.com/dingkegithub/com.dk.user/logic/service"
	"github.com/dingkegithub/com.dk.user/sidecar/discovery"
	"github.com/dingkegithub/com.dk.user/sidecar/discovery/kitnacos"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/lb"
	"time"
)

//
// 用户服务代理，在服务至上添加额外功能
//
func UserLogicProxy(ctx context.Context, logger log.Logger, cli discovery.RegisterCenterClient) service.UserLogicServiceMiddleware {
	var (
		maxAttempts = 3   // 失败尝试次数
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
	endpointer := sd.NewEndpointer(instance, NewServiceFactory, logger)

	// 将站点加入到负载均衡器
	balancer := lb.NewRoundRobin(endpointer)
	retry := lb.Retry(maxAttempts, maxTime, balancer)

	// 返回含有服务发现，负载均衡的service
	return func(svc service.UserLogicService) service.UserLogicService {
		return NewUsrLogicProxy(ctx, svc, logger, retry)
	}
}

type UsrLogicProxy struct {
	logger              log.Logger
	ctx                 context.Context
	nxt                 service.UserLogicService
	registerRpcEndpoint endpoint.Endpoint
}

func NewUsrLogicProxy(ctx context.Context, svc service.UserLogicService, logger log.Logger, ep endpoint.Endpoint) service.UserLogicService {
	return &UsrLogicProxy{
		logger:              logger,
		ctx:                 ctx,
		nxt:                 svc,
		registerRpcEndpoint: ep,
	}
}

//
// 用户注册接口
// Endpoint 调用
//
func (u *UsrLogicProxy) Register(ctx context.Context, request *service.RegisterUsrRequest) (*service.RegisterUsrResponse, error) {

	u.logger.Log("file", "userproxy.go", "function", "Register", "action", "invoke Register")
	// 服务一层一层，如同洋葱一般包裹，每一层接口相同
	rs := time.Now()
	_, err := u.nxt.Register(ctx, request)
	re := time.Since(rs)
	if err != nil {
		u.logger.Log("file", "userproxy.go", "function", "Register", "action", "invoke next Register", "error", err)
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

	return resp.(*service.RegisterUsrResponse), nil
}