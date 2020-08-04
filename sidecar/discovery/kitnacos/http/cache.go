package http

import (
	"github.com/dingkegithub/com.dk.user/utils/osutils"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"sync"
)

type Instance struct {
	Ip          string            `json:"ip"`
	Port        uint64            `json:"port"`
	ServiceName string            `json:"serviceName"`
	ClusterName string            `json:"clusterName"`
	Enable      bool              `json:"enable"`
	InstanceId  string            `json:"instanceId"`
	Metadata    map[string]string `json:"metadata"`
	Weight      float64           `json:"weight"`
}

type LocalCache struct {
	cacheFile string
	mutex sync.RWMutex
	instance map[string][]*Instance
}

func (lc *LocalCache) Instance(name string) []*Instance {

	lc.mutex.RLock()
	defer lc.mutex.RUnlock()
	insts, ok := lc.instance[name]
	if !ok {
		return nil
	}

	res := make([]*Instance, len(insts))
	copy(res, insts)

	return res
}

func (lc *LocalCache) Store(name string, inst []*Instance) error  {
	lc.mutex.Lock()
	defer lc.mutex.Unlock()
	newInsts := make([]*Instance, len(inst))
	copy(newInsts, inst)
	lc.instance[name] = newInsts
	body, err := json.Marshal(lc.instance)
	if err != nil {
		return err
	}

	if ! (osutils.Exists(lc.cacheFile)) {
		err = osutils.Touch(lc.cacheFile)
		if err != nil {
			return err
		}
	}

	f, err := os.OpenFile(lc.cacheFile, os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	size, err := f.Write(body)
	if err != nil {
		return err
	}

	fmt.Println("file", "cache.go", "function", "update", "action", "write", "size", size)
	return nil
}

func (lc *LocalCache) Load() error {
	if ! osutils.Exists(lc.cacheFile) {
		return fmt.Errorf("not found cache file")
	}

	f, err := os.Open(lc.cacheFile)
	if err != nil {
		return err
	}

	content, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	instanceList := make(map[string][]*Instance)
	err = json.Unmarshal(content, &instanceList)
	if err != nil {
		return err
	}

	lc.mutex.Lock()
	defer lc.mutex.Unlock()
	lc.instance = instanceList

	return nil
}

func NewLocalCache(dir string) (*LocalCache, error) {
	localCache := &LocalCache{
		instance: make(map[string][]*Instance),
	}

	if ! (osutils.Exists(dir) && osutils.IsDir(dir)) {
		err := osutils.Mkdir(dir, true)
		if err != nil {
			return nil, err
		}
	}

	cacheFile := path.Join(dir, "instance")
	if ! (osutils.Exists(cacheFile) && osutils.IsFile(cacheFile)) {
		err := osutils.Touch(cacheFile)
		if err != nil {
			return nil, err
		}
	}

	localCache.cacheFile = cacheFile
	return localCache, nil
}

