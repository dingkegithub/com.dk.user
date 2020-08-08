package main

import (
	"context"
	"fmt"
)

const (
	UserName = "bql"
	Password = "123456"
)

type LoginSvc interface {
	Login(ctx context.Context, usr, pwd string) (string, int, error)

	Logout(ctx context.Context, usr string) string
}


type defaultLoginSvc struct {

}

func (d *defaultLoginSvc) Login(ctx context.Context, usr, pwd string) (string, int, error) {
	if usr != UserName {
		return "", 401, fmt.Errorf("%s not register", usr)
	}

	if pwd != Password {
		return "", 402, fmt.Errorf("pwd not match with usr %s", usr)
	}
    return fmt.Sprintf("token-%s-%s", usr, pwd), 200, nil
}

func (d *defaultLoginSvc) Logout(ctx context.Context, usr string) string {
	return fmt.Sprintf("%s loginout success", usr)
}

func NewDefaultService() LoginSvc {
	return &defaultLoginSvc{}
}

