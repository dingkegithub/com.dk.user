package endpoints

import "github.com/dingkegithub/com.dk.user/das/model"

type UsrCreateRequest struct {
	Data *model.User `json:"data"`
}

type UsrUpdateRequest struct {
	Uid  uint64                 `json:"uid"`
	Data map[string]interface{} `json:"data"`
}

type UsrRetrieveRequest struct {
	Uid uint64 `json:"uid"`
}

type UsrModelResponse struct {
	Usr *model.User `json:"usr"`
	Err error       `json:"err"`
}

type UsrListRequest struct {
	Limit  int64                  `json:"limit"`
	Offset int64                  `json:"offset"`
	Data   map[string]interface{} `json:"data"`
}

type UsrListResponse struct {
	Usr []*model.User `json:"usr"`
	Err error         `json:"err"`
}
