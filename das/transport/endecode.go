package transport

import (
	"com.dk.user/das/endpoints"
	"com.dk.user/das/model"
	"com.dk.user/das/proto/userpb"
	"com.dk.user/das/service"
	"context"
	"fmt"
)

func decodeCreateRequest(_ context.Context, r interface{}) (interface{}, error) {
	fmt.Println("file", "endecode.go", "function", "decodeCreateRequest", "action", "invoke")
	req := r.(*userpb.RegisterRequest)

	usrModel := &model.User{
		Uid:     req.Uid,
		Name:    req.Name,
		Pwd:     req.Pwd,
	}
	fmt.Println("file", "endecode.go", "function", "decodeCreateRequest", "action", "decode over")
	return usrModel, nil
}

func decodeRetrieveRequest(_ context.Context, r interface{}) (interface{}, error) {
	req, ok := r.(*userpb.RetrieveRequest)
	if !ok {
		return nil, service.ErrParam
	}

	return &endpoints.UsrRetrieveRequest{Uid:req.Uid}, nil
}

func decodeUpdateRequest(_ context.Context, r interface{}) (interface{}, error)  {
	req, ok := r.(*userpb.UpdateRequest)
	if !ok {
		return nil, service.ErrParam
	}

	updateData := make(map[string]interface{})
	for _, v := range req.Data {
		updateData[v.Key] = v.Val
	}

	return &endpoints.UsrUpdateRequest{
		Uid:  req.Uid,
		Data: updateData,
	}, nil

}

func decodeListRequest(_ context.Context, r interface{}) (interface{}, error) {
	req, ok := r.(*userpb.ListRequest)
	if !ok {
		return nil, service.ErrParam
	}

	qs := make(map[string]interface{})
	for _, lq := range req.Qs {
		qs[lq.Key] = lq.Val
	}

	return &endpoints.UsrListRequest{
		Limit:  req.Limit,
		Offset: req.Offset,
		Data:   qs,
	}, nil
}

func encodeModelResponse(_ context.Context, r interface{}) (interface{}, error) {
	fmt.Println("file", "endecode.go", "function", "encodeModelResponse", "action", "invoke")
	resp, ok := r.(*endpoints.UsrModelResponse)
	if !ok {
		fmt.Println("file", "endecode.go", "function", "encodeModelResponse", "action", "interface insert", "error", ok)
		return nil, service.ErrSvcInner
	}

	return &userpb.RegisterResponse{
		//Err: service.ErrMapToCode(resp.Err),
		//Msg: resp.Err.Error(),
		Err: 20000,
		Msg: "ok",
		Data: &userpb.UserData{
			Uid:  resp.Usr.Uid,
			Name: resp.Usr.Name,
		},
	}, nil
}

func encodeModelListResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp, ok := r.(*endpoints.UsrListResponse)
	if !ok {
		return nil, service.ErrSvcInner
	}

	usrData := make([]*userpb.UserData, 0, len(resp.Usr))
	for _, usr := range resp.Usr {
		usrData = append(usrData, &userpb.UserData{
			Uid:                  usr.Uid,
			Name:                 usr.Name,
		})
	}

	return &userpb.ListResponse{
		Err:                  service.ErrMapToCode(resp.Err),
		Msg:                  resp.Err.Error(),
		Data:                 usrData,
	}, nil
}
