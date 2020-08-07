package proxymiddleware

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dingkegithub/com.dk.user/das/proto/userpb"
	"github.com/dingkegithub/com.dk.user/logic/service"
	"github.com/dingkegithub/com.dk.user/sidecar/discovery"
	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/ratelimit"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"github.com/sony/gobreaker"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"io"
	"time"
)

type ServiceFactory struct {
	instance string
	meta *discovery.ServiceMeta
	conn *grpc.ClientConn
	logger log.Logger
}


//
// @param instance 服务发现随机取得的实例
//
// @return endpoint 站点信息
// @return io.Closer 站点失效后如何关闭
//
func NewServiceFactory(instance string) (endpoint.Endpoint, io.Closer, error) {
	meta := &discovery.ServiceMeta{}
	err := json.Unmarshal([]byte(instance), &meta)
	if err != nil {
		return nil, nil, err
	}

	hostPort := fmt.Sprintf("%s:%d", meta.Ip, meta.Port)
	fmt.Println("file", "svcfactory.go",
		"function", "NewServiceFactory",
		"action", "dialling",
		"addr", hostPort)

	conn, err := grpc.Dial(hostPort, grpc.WithInsecure())
	if err != nil {
		fmt.Println("file", "svcfactory.go",
			"function", "NewServiceFactory",
			"action", "dial",
			"error", err)
		return nil, nil, err
	}

	svcFactory := &ServiceFactory{
		instance: instance,
		conn: conn,
		meta: meta,
	}

	return svcFactory.Endpoint(), svcFactory, nil
}

//
// 实际对应站点
//
func (sf *ServiceFactory) Endpoint() endpoint.Endpoint {
	var registerEp endpoint.Endpoint
	{
		registerEp = kitgrpc.NewClient(
			sf.conn,
			"userpb.UserDasService",
			"Create",
			encodeRegisterRpcRequest,
			decodeRegisterRpcResponse,
			&userpb.RegisterResponse{},
		).Endpoint()
		fmt.Println("file", "svcfactory.go",
			"function", "Endpoint",
			"action", "NewClient",
			"svc", sf.meta.SvcName)
	}

	qps  := 100
	registerEp = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(registerEp)
	return ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), qps))(registerEp)
}

//
// 站点关闭
//
func (sf *ServiceFactory) Close() error {
	fmt.Println("file", "svcfactory.go",
		"function", "close",
		"action", "close endpoint",
		"endpoint", sf.instance)
	return sf.conn.Close()
}

func encodeRegisterRpcRequest(ctx context.Context, req interface{}) (interface{}, error)  {
	fmt.Println("file", "svcfactory.go", "function", "encodeRegisterRpcRequest", "action", "cvt request to pb requext")
	svcRequest := req.(*service.RegisterUsrRequest)

	u := &userpb.RegisterRequest{
		Uid:  svcRequest.Uid,
		Name: svcRequest.Name,
		Pwd:  svcRequest.Pwd,
	}

	return u, nil
}

func decodeRegisterRpcResponse(ctx context.Context, resp interface{}) (interface{}, error)  {
	fmt.Println("file", "svcfactory.go", "function", "decodeRegisterRpcRequest", "action", "cvt pb response to response")
	dasResp, ok := resp.(*userpb.RegisterResponse)
	if !ok {
		fmt.Println("file", "svcfactory.go", "function", "decodeRegisterRpcRequest", "action", "insert pb type failed")
		return nil, service.ErrorServerInternal
	}

	if dasResp.Data == nil {
		fmt.Println("file", "svcfactory.go", "function", "decodeRegisterRpcRequest", "action", "das resp error")
		return &service.RegisterUsrResponse{
			Err:  dasResp.Err,
			Msg:  dasResp.Msg,
			Data: nil,
		}, nil
	}

	return &service.RegisterUsrResponse{
		Err:  service.ErrorCodeSuccess,
		Msg:  service.ErrorSuccess.Error(),
		Data: &service.RegisterUsrDetail{
			Id:   dasResp.Data.Uid,
			Uid:  dasResp.Data.Uid,
			Name: dasResp.Data.Name,
		},
	}, nil
}