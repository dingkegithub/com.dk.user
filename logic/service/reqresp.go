package service

type RegisterUsrRequest struct {
	Uid      string `json:"uid"`
	Name     string `json:"Name"`
	Pwd      string `json:"Pwd"`
	PwdAgain string `json:"pwd_again"`
	Age      int    `json:"Age"`
}

type RegisterUsrDetail struct {
	Id string `json:"id"`
	Uid string `json:"uid"`
	Name string `json:"name"`
}

type RegisterUsrResponse struct {
	Err int64 `json:"err"`
	Msg string `json:"msg"`
	Data *RegisterUsrDetail `json:"data"`
}

type LoginRequest struct {
	Name     string `json:"Name"`
	Pwd      string `json:"Pwd"`
}

type LoginOutRequest struct {
}

type ResetPwdRequest struct {
	Uid    string `json:"uid"`
	OldPwd string `json:"old_pwd"`
	NewPwd string `json:"new_pwd"`
}

type UpdateUsrRequest struct {
	Uid string `json:"uid"`
	Data string `json:"data"`
}