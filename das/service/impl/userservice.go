package impl

import (
	"context"
	"github.com/dingkegithub/com.dk.user/das/model"
	"github.com/dingkegithub/com.dk.user/das/service"
	"github.com/dingkegithub/com.dk.user/utils/logging"
	"github.com/jinzhu/gorm"
)

type UserSrv struct {
	logger logging.Logger
}

func NewUserSrv(logger logging.Logger) service.UserSvc {
	return &UserSrv{
		logger: logger,
	}
}

func (s *UserSrv) Update(ctx context.Context, uid uint64, data *model.User) (*model.User, error) {
	db := model.DM().Db()

	errObj := db.Model(&model.User{}).Where(&model.User{Uid: uid}).Update(data)
	if errObj.Error != nil {
		s.logger.Log("file", "userservice.go",
			"func", "update",
			"msg", "update record failed",
			"uid", uid,
			"error", errObj.Error)
		return nil, service.ErrQuery
	}

	return s.Retrieve(ctx, uid)
}

func (s *UserSrv) Retrieve(_ context.Context, uid uint64) (*model.User, error) {
	db := model.DM().Db()

	var user *model.User
	err := db.Model(&model.User{}).Where(&model.User{Uid: uid}).First(&user).Error
	if err != nil {
		s.logger.Log("file", "userservice.go",
			"func", "Retrieve",
			"msg", "query record failed",
			"error", err)
		return nil, service.ErrQuery
	}

	return user, nil
}

func (s *UserSrv) List(_ context.Context, data *model.User, limit, offset int64) ([]*model.User, error) {
	db := model.DM().Db()

	var users []*model.User
	err := db.Model(&model.User{}).Where(data).Offset(offset).Limit(limit).Find(&users).Error
	if err != nil {
		s.logger.Log("file", "userservice.go",
			"func", "List",
			"msg", "query users failed",
			"offset", offset,
			"limit", limit,
			"error", err)
		return nil, service.ErrQuery
	}

	return users, nil
}

func (s *UserSrv) Create(_ context.Context, usr *model.User) (*model.User, error) {
	var err error

	db := model.DM().Db()

	newUsr := &model.User{}
	err = db.Model(&model.User{}).Where(&model.User{Name: usr.Name}).First(&newUsr).Error
	if err == gorm.ErrRecordNotFound {
		err = db.Create(usr).Error
		if err != nil {
			s.logger.Log("file", "userservice.go",
				"func", "Create",
				"msg", "insert new user failed",
				"error", err)
			return nil, service.ErrUnknown
		}

		return usr, nil
	}

	s.logger.Log("file", "userservice.go",
		"func", "Create",
		"msg", "insert new user failed",
		"error", err)
	return nil, service.ErrExist
}