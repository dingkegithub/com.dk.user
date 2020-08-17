package setting

import (
	"github.com/dingkegithub/com.dk.user/sidecar/config"
	"sync"
)

var once sync.Once
var appCfg config.CfgCenterClient

func New(cli config.CfgCenterClient) {
	once.Do(func() {
		appCfg = cli
	})
}

func ApplicationConfig() config.CfgCenterClient  {
	return appCfg
}
