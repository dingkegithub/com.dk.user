package endpoints

import "com.dk.user/das/model"


type UsrCreateRequest struct {
	Data *model.User
}

type UsrUpdateRequest struct {
	Uid string `json:"uid"`
	Data map[string]interface{} `json:"data"`
}

type UsrRetrieveRequest struct {
	Uid string `json:"uid"`
}

type UsrModelResponse struct {
	Usr *model.User
	Err error
}

type UsrListRequest struct {
	Limit int64 `json:"limit"`
	Offset int64 `json:"offset"`
	Data map[string]interface{} `json:"data"`
}

type UsrListResponse struct {
	Usr []*model.User
	Err error
}


