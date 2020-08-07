package localcfg


import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dingkegithub/com.dk.user/utils/osutils"
	"io/ioutil"
)


type CfgLoader struct {
	cfgFile string
	localCfg *LocalCfg
}

var (
	cfgLoader *CfgLoader
)

func GetCfg() *CfgLoader {
	return cfgLoader
}

func NewCfgLoader(f string) (*CfgLoader, error) {

	if ! (osutils.Exists(f) && osutils.IsFile(f)) {
		return nil, errors.New(fmt.Sprintf("file not exist: %s", f))
	}

	data, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, err
	}

	localCfg := &LocalCfg{}
	err = json.Unmarshal(data, localCfg)
	if err != nil {
		return nil, err
	}

	cfgLoader = &CfgLoader {
		cfgFile:  f,
		localCfg: localCfg,
	}

	return cfgLoader, nil
}

func (cl CfgLoader) GetLogCfg() *Log {
	return cl.localCfg.Log
}

func (cl CfgLoader) GetApolloCfg() *ApolloParam {
	return cl.localCfg.Apollo
}