package impl

import (
	"com.dk.user/das/model"
	"com.dk.user/das/service"
	"context"
	"fmt"
	"github.com/astaxie/beego/orm"
	"time"
)

type BeeUserSvc struct {
	tbOrm orm.Ormer
	tbQs orm.QuerySeter
}


func NewUserSvc() service.UserSvc {
	tbOrm := orm.NewOrm()
	tbQs := tbOrm.QueryTable("user")
	return &BeeUserSvc{
		tbOrm: tbOrm,
		tbQs:  tbQs,
	}
}

func (s *BeeUserSvc) Update(ctx context.Context, uid string, data map[string]interface{}) (*model.User, error) {
	num, err := s.tbQs.Filter("uid", uid).Update(data)
	if err != nil {
		return nil, s.err(err)
	}

	fmt.Println("update rows: ", num)

	return s.Retrieve(ctx, uid)
}

func (s *BeeUserSvc) Retrieve(_ context.Context, uid string) (*model.User, error) {
	usr := &model.User{Uid: uid}
	err := s.tbQs.One(usr)
	if err != nil {
		return nil, s.err(err)
	}

	return usr, nil
}

func (s *BeeUserSvc) List(_ context.Context, data map[string]interface{}, limit, offset int64) ([]*model.User, error) {
	qs := s.tbQs
	for k, v := range data {
		qs = qs.Filter(k, v)
	}

	var users []*model.User
	num, err := qs.Offset(offset).Limit(limit).All(users)
	if err != nil {
		return nil, s.err(err)
	}

	fmt.Println("query record num: ", num)
	return users, nil
}

func (s *BeeUserSvc) Create(_ context.Context, usr *model.User) (*model.User, error) {
	var err error
	//
	//infoStr, err := json.Marshal(data)
	//if err != nil {
	//	return nil, logic.ErrParam
	//}
	//
	//usr := &model.User{}
	//err = json.Unmarshal(infoStr, &usr)
	//if err != nil {
	//	return nil, logic.ErrParam
	//}
	fmt.Println("file", "usersvc.go", "function", "Create", "action", "invoke create")
	start := time.Now()
	_, err = s.tbOrm.Insert(usr)
	end := time.Since(start)

	fmt.Println("file", "usersvc.go", "function", "Create", "action", "invoke create", "lost", end)
	if err != nil {
		fmt.Println("file", "usersvc.go", "function", "Create", "action", "invoke create", "error", err)
		return nil, s.err(err)
	}

	return usr, nil
}

func (s *BeeUserSvc) err(err error) error {
	switch err {
	case orm.ErrNoRows:
		return service.ErrNotFound
	case orm.ErrMultiRows:
		return service.ErrExist
	case orm.ErrArgs, orm.ErrMissPK:
		return service.ErrParam
	default:
		return service.ErrQuery
	}
}

