package service

import (
	"context"
	"github.com/dingkegithub/com.dk.user/das/model"
)

type UserSvc interface {

	Update(ctx context.Context, uid string, data map[string]interface{}) (*model.User, error)

	Retrieve(ctx context.Context, uid string) (*model.User, error)

	List(ctx context.Context, data map[string]interface{}, limit, offset int64) ([]*model.User, error)

	Create(ctx context.Context, data *model.User) (*model.User, error)
}

type UserSvcMiddleware func(UserSvc) UserSvc
