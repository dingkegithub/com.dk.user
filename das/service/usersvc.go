package service

import (
	"context"
	"github.com/dingkegithub/com.dk.user/das/model"
)

type UserSvc interface {

	Update(ctx context.Context, uid uint64, data *model.User) (*model.User, error)

	Retrieve(ctx context.Context, uid uint64) (*model.User, error)

	List(ctx context.Context, data map[string]interface{}, limit, offset int64) ([]*model.User, error)

	Create(ctx context.Context, data *model.User) (*model.User, error)
}

type UserSvcMiddleware func(UserSvc) UserSvc
