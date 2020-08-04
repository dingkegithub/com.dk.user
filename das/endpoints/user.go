package endpoints

import (
	"context"
	"fmt"
	"github.com/dingkegithub/com.dk.user/das/model"
	"github.com/dingkegithub/com.dk.user/das/service"
	"github.com/go-kit/kit/endpoint"
	"time"
)

type UsrEndpoints struct {
	CreateEndpoint   endpoint.Endpoint
	UpdateEndpoint   endpoint.Endpoint
	RetrieveEndpoint endpoint.Endpoint
	ListEndpoint     endpoint.Endpoint
}

func NewUsrEndpoints(svc service.UserSvc) *UsrEndpoints {
	return &UsrEndpoints{
		CreateEndpoint:   MakeCreateEndpoint(svc),
		UpdateEndpoint:   MakeUpdateEndpoint(svc),
		RetrieveEndpoint: MakeRetrieveEndpoint(svc),
		ListEndpoint:     MakeListEndpoint(svc),
	}
}


func MakeCreateEndpoint(svc service.UserSvc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		fmt.Println("file", "user.go", "function", "MakeCreateEndpoint", "action", "invoke endpoint")
		req := request.(*model.User)
		start := time.Now()
		resp, err := svc.Create(ctx, req)
		end := time.Since(start)
		fmt.Println("file", "user.go", "function", "MakeCreateEndpoint", "action", "invoke endpoint", "lost", end)
		if err != nil {
			fmt.Println("file", "user.go", "function", "MakeCreateEndpoint", "action", "invoke endpoint", "error", err)
			return nil, err
		}
		fmt.Println("file", "user.go", "function", "MakeCreateEndpoint", "action", "invoke endpoint", "info", resp.String())
		return &UsrModelResponse{
			Usr: resp,
			Err: nil,
		}, nil
	}
}

func MakeUpdateEndpoint(svc service.UserSvc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*UsrUpdateRequest)
		resp, err := svc.Update(ctx, req.Uid, req.Data)
		if err != nil {
			return nil, err
		}

		return resp, nil
	}
}

func MakeRetrieveEndpoint(svc service.UserSvc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*UsrRetrieveRequest)
		resp, err := svc.Retrieve(ctx, req.Uid)
		if err != nil {
			return nil, err
		}
		return resp, nil
	}
}

func MakeListEndpoint(svc service.UserSvc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*UsrListRequest)
		resp, err := svc.List(ctx, req.Data, req.Limit, req.Offset)
		if err != nil {
			return nil, err
		}
		return resp, nil
	}
}
