package middleware


import (
	"com.dk.user/das/model"
	"com.dk.user/das/service"
	"context"
	"github.com/go-kit/kit/log"
	"time"
)

type loggingMiddleware struct {
	log log.Logger
	service.UserSvc
}

func (l loggingMiddleware) Update(ctx context.Context, uid string, data map[string]interface{}) (*model.User, error) {
	sT := time.Now()
	defer func(s time.Time) {
		_ = l.log.Log("function", "Update", "took", time.Since(s))
	}(sT)
	return l.UserSvc.Update(ctx, uid, data)
}

func (l loggingMiddleware) Retrieve(ctx context.Context, uid string) (*model.User, error) {
	sT := time.Now()
	defer func(s time.Time) {
		_ = l.log.Log("function", "Retrieve", "took", time.Since(s))
	}(sT)
	return l.UserSvc.Retrieve(ctx, uid)
}

func (l loggingMiddleware) List(ctx context.Context, data map[string]interface{}, limit, offset int64) ([]*model.User, error) {
	sT := time.Now()
	defer func(s time.Time) {
		_ = l.log.Log("function", "Retrieve", "took", time.Since(s))
	}(sT)
	return l.UserSvc.List(ctx, data, limit, offset)
}

func (l loggingMiddleware) Create(ctx context.Context, data *model.User) (*model.User, error) {
	sT := time.Now()
	defer func(s time.Time) {
		_ = l.log.Log("function", "Retrieve", "took", time.Since(s))
	}(sT)
	return l.Create(ctx, data)
}

func LoggingMiddleware(logger log.Logger) service.UserSvcMiddleware {
	return func(svc service.UserSvc) service.UserSvc {
		return loggingMiddleware{
			log:     logger,
			UserSvc: svc,
		}
	}
}
