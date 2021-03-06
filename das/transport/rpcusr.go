package transport

import (
	"context"
	"github.com/dingkegithub/com.dk.user/das/endpoints"
	"github.com/dingkegithub/com.dk.user/das/proto/userpb"
	"github.com/dingkegithub/com.dk.user/utils/logging"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"time"
)

type rpcUsrTransport struct {
	logger logging.Logger

	create grpctransport.Handler

	retrieve grpctransport.Handler

	list grpctransport.Handler

	update grpctransport.Handler
}

func (u *rpcUsrTransport) Create(ctx context.Context, r *userpb.RegisterRequest) (*userpb.RegisterResponse, error) {
	start := time.Now()
	_, resp, err := u.create.ServeGRPC(ctx, r)
	lost := time.Since(start)
	u.logger.Log("file", "rpcusr.go",
		"func", "Create",
		"msg", "rpc ceate invoked",
		"lost", lost)

	if err != nil {
		return nil, err
	}

	return resp.(*userpb.RegisterResponse), nil
}

func (u *rpcUsrTransport) Retrieve(ctx context.Context, r *userpb.RetrieveRequest) (*userpb.RetrieveResponse, error) {
	start := time.Now()
	_, resp, err := u.retrieve.ServeGRPC(ctx, r)
	lost := time.Since(start)

	u.logger.Log("file", "rpcusr.go",
		"func", "Create",
		"msg", "rpc retrieve invoked",
		"lost", lost)

	if err != nil {
		return nil, err
	}

	return resp.(*userpb.RetrieveResponse), nil
}

func (u *rpcUsrTransport) List(ctx context.Context, r *userpb.ListRequest) (*userpb.ListResponse, error) {
	start := time.Now()
	_, resp, err := u.list.ServeGRPC(ctx, r)
	lost := time.Since(start)

	u.logger.Log("file", "rpcusr.go",
		"func", "Create",
		"msg", "rpc retrieve invoked",
		"lost", lost)

	if err != nil {
		return nil, err
	}

	return resp.(*userpb.ListResponse), nil
}

func (u *rpcUsrTransport) Update(ctx context.Context, r *userpb.UpdateRequest) (*userpb.UpdateResponse, error) {
	start := time.Now()
	_, resp, err := u.update.ServeGRPC(ctx, r)
	lost := time.Since(start)

	u.logger.Log("file", "rpcusr.go",
		"func", "Create",
		"msg", "rpc retrieve invoked",
		"lost", lost)
	if err != nil {
		return nil, err
	}

	return resp.(*userpb.UpdateResponse), nil
}

func NewRpcUsrSvc(ctx context.Context, logger logging.Logger, endpoints *endpoints.UsrEndpoints) userpb.UserDasServiceServer {
	return &rpcUsrTransport{
		logger: logger,

		create: grpctransport.NewServer(
			endpoints.CreateEndpoint,
			decodeCreateRequest,
			encodeModelResponse,
		),

		retrieve: grpctransport.NewServer(
			endpoints.RetrieveEndpoint,
			decodeRetrieveRequest,
			encodeModelResponse,
		),

		list: grpctransport.NewServer(
			endpoints.ListEndpoint,
			decodeListRequest,
			encodeModelListResponse,
		),

		update: grpctransport.NewServer(
			endpoints.UpdateEndpoint,
			decodeUpdateRequest,
			encodeModelResponse,
		),
	}
}
