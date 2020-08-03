package model

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	Id      int64    `json:"Id"`
	Uid     string   `json:"Uid"`
	Name    string   `json:"Name"`
	Pwd     string   `json:"Pwd"`
	//Profile *Profile `json:"Profile" orm:"rel(one)"`
}

func (u User) String() string {
	b, _ := json.Marshal(u)
	return string(b)
}

type Profile struct {
	Id   string `json:"Id"`
	Pid  string `json:"Pid" orm:"pk"`
	//User *User  `json:"User" orm:"reverse(one)"`
	Age  int    `json:"Age"`
}

func Init(user, pwd, db string) {
	fmt.Println("0, err: ", user, pwd, db)

	dbInfo := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4",
		user, pwd, "127.0.0.1", 3306, db)

	err := orm.RegisterDataBase("default", "mysql", dbInfo, 30)

	fmt.Println("1, err: ", err)
	orm.RegisterModel(new(User))
	err = orm.RunSyncdb("default", false, true)
	fmt.Println("2, err: ", err)
}
