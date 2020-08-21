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
	"github.com/go-kit/kit/ratelimit"
	kitrpc "github.com/go-kit/kit/transport/grpc"
	"github.com/sony/gobreaker"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"io"
	"time"
)

type SvcLoginFactory struct {
	conn *grpc.ClientConn
	meta *discovery.ServiceMeta
	qps int
}

func NewSvcLoginFactory(instance string) (endpoint.Endpoint, io.Closer, error) {
	meta := &discovery.ServiceMeta{}

	err := json.Unmarshal([]byte(instance), meta)
	if err != nil {
		return nil, nil, err
	}

	addr := fmt.Sprintf("%s:%d", meta.Ip, meta.Port)
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, nil, err
	}

	factory := &SvcLoginFactory{
		conn: conn,
		meta: meta,
		qps: 1000,
	}

	return factory.Endpoint(), factory, nil
}

func (sf *SvcLoginFactory) Endpoint() endpoint.Endpoint {
	ep := kitrpc.NewClient(
		sf.conn,
		sf.meta.SvcName,
		"Login",
		encodeLoginRpcRequest,
		decodeLoginRpcRequest,
		userpb.ListResponse{},
	).Endpoint()

	ep = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(ep)
	return ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), sf.qps))(ep)

}

func (sf *SvcLoginFactory) Close() error {
	return sf.conn.Close()
}

func decodeLoginRpcRequest(ctxt context.Context, request interface{}) (interface{}, error)  {
	resp := request.(*userpb.ListResponse)
	return &service.Response{
		Error: resp.Err,
		Msg:   resp.Msg,
		Data:  resp.Data,
	}, nil
}

func encodeLoginRpcRequest(ctxt context.Context, request interface{}) (interface{}, error)  {
	req := request.(*service.LoginRequest)

	return &userpb.ListRequest{
		Limit:  1,
		Offset: 0,
		Qs:     []*userpb.UserData{
			&userpb.UserData{
				Uid:  0,
				Name: req.Name,
			},
		},
	}, nil
}