package model

import (
	"encoding/json"
	"fmt"
	"github.com/dingkegithub/com.dk.user/das/setting"
	"github.com/dingkegithub/com.dk.user/utils/logging"
	"github.com/jinzhu/gorm"
	"sync"
)

type mysqlConfig struct {
	Ip       string `json:"ip"`
	Port     uint16 `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Db       string `json:"db"`
}

func (m *mysqlConfig) url() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4",
		m.User, m.Password, m.Ip, m.Port, m.Db)
}

func (m *mysqlConfig) valid() bool {
	if m.Ip == "" {
		return false
	}

	if m.Port <= 0 {
		return false
	}

	if m.User == "" {
		return false
	}

	return true
}

var onece sync.Once
var dbManager *DbManager

func DM() *DbManager {
	return dbManager
}

type DbManager struct {
	cfgMutex sync.RWMutex
	mysqlCfg *mysqlConfig
	logger   logging.Logger
	dbMutex sync.RWMutex
	dbConn   *gorm.DB
}


func New(logger logging.Logger) {
	onece.Do(func() {
		dm := &DbManager{
			logger: logger,
		}

		ap := setting.ApplicationConfig().Get("db", "application")
		err := json.Unmarshal([]byte(ap), &dm.mysqlCfg)
		if err != nil {
			logger.Log("file", "user.go",
				"func", "New",
				"msg", "unmarshal db config failed",
				"error", err)
			panic(err)
		}

		if dm.mysqlCfg.Ip == "" {
			dm.mysqlCfg.Ip = "127.0.0.1"
		}

		if dm.mysqlCfg.Port <= 0 {
			dm.mysqlCfg.Port = 3306
		}

		if dm.mysqlCfg.User == "" {
			dm.mysqlCfg.User = "root"
		}

		dbManager = dm
		setting.ApplicationConfig().Register("db", dbManager.listener)
	})
}

func (dm *DbManager) Close() {
	setting.ApplicationConfig().Deregister("db")
	dm.dbMutex.Lock()
	defer dm.dbMutex.Unlock()
	dm.dbConn.Close()
	dm.dbConn = nil
}

func (dm *DbManager) reconnect(url string) *gorm.DB {
	mysqlUrl := fmt.Sprintf("%s?charset=utf8&parseTime=True&loc=Local", url)
	conn, err := gorm.Open("mysql", mysqlUrl)
	if err != nil {
		dm.logger.Log("file", "dbmanager.go",
			"func", "reconnect",
			"msg", "connect db failed",
			"url", mysqlUrl,
			"error", err)
		return nil
	}

	err = dm.dbConn.Close()
	if err != nil {
		dm.logger.Log("file", "dbmanager.go",
			"func", "reconnect",
			"msg", "close db failed",
			"url", mysqlUrl,
			"error", err)
		return nil
	}

	return conn
}

func (dm *DbManager) listener(key string, value interface{}) {
	dbCfgStr := value.(string)
	dbCfg := &mysqlConfig{}
	err := json.Unmarshal([]byte(dbCfgStr), &dbCfg)
	if err != nil {
		dm.logger.Log("file", "user.go",
			"func", "listener",
			"msg", "unmarshal db config failed",
			"error", err)
		return
	}

	if ! dbCfg.valid() {
		dm.logger.Log("file", "user.go",
			"func", "listener",
			"msg", "check config format invalid")
		return
	}

	newUrl := dbCfg.url()
	conn := dm.reconnect(newUrl)
	if conn == nil {
		return
	}

	dm.cfgMutex.Lock()
	dm.mysqlCfg = dbCfg
	dm.cfgMutex.Unlock()

	dm.dbMutex.Lock()
	defer dm.dbMutex.Unlock()
	dm.dbConn = conn
}

func (dm *DbManager) Db() *gorm.DB {
	dm.dbConn.RLock()
	defer dm.dbConn.Unlock()
	return dm.dbConn
}