package model

import (
	"encoding/json"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"time"
)

type User struct {
	Uid       uint64     `gorm:"primary_key" json:"Uid"`
	Name      string     `json:"Name"`
	Pwd       string     `json:"Pwd"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

func (u User) TableName() string {
	return "user"
}

func (u User) String() string {
	b, _ := json.Marshal(u)
	return string(b)
}
