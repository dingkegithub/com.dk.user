package service

type LoginToken struct {
	Uid uint64 `json:"uid"`
	Token string `json:"token"`
}

type Response struct {
	Error int64       `json:"error"`
	Msg   string      `json:"msg"`
	Data  interface{} `json:"data"`
}

type RegisterUsrDetail struct {
	Uid  uint64 `json:"uid"`
	Name string `json:"name"`
}

type RegisterUsrRequest struct {
	Uid      uint64 `json:"uid"`
	Name     string `json:"Name"`
	Pwd      string `json:"Pwd"`
	PwdAgain string `json:"pwd_again"`
	Age      int    `json:"Age"`
}

type LoginRequest struct {
	Name string `json:"Name"`
	Pwd  string `json:"Pwd"`
}

type LoginOutRequest struct {
	Uid   uint64 `json:"uid"`
	Token string `json:"token"`
}

type ResetPwdRequest struct {
	Uid    uint64 `json:"uid"`
	Token  string `json:"token"`
	OldPwd string `json:"old_pwd"`
	NewPwd string `json:"new_pwd"`
}

type UpdateUsrRequest struct {
	Uid  uint64 `json:"uid"`
	Name string `json:"Name"`
}