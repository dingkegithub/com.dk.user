package http

import (
	"encoding/json"
	"github.com/dingkegithub/com.dk.user/sidecar/discovery"
	"github.com/dingkegithub/com.dk.user/utils/osutils"
	"github.com/go-kit/kit/log"
	"io/ioutil"
	"os"
	"path"
	"sync"
)

type CacheInstance struct {
	Mils      uint64                  `json:"mils"`
	Instances []*discovery.ServiceMeta `json:"instances"`
}

type LocalCache struct {
	cacheFile string                    // local cache file
	logger    log.Logger                // log interface, need implement Log(kv... interface{})
	mutex     sync.RWMutex              // protect memory cache instance
	instance  map[string]*CacheInstance // memory cache
}

// read service instance list from cache
// @param name service tag
// @return []*Instance instance list
func (lc *LocalCache) Instance(name string) *CacheInstance {

	lc.mutex.RLock()
	defer lc.mutex.RUnlock()

	insts, ok := lc.instance[name]
	if !ok {
		return nil
	}

	c := &CacheInstance{
		Mils:      insts.Mils,
		Instances: nil,
	}

	b, err := json.Marshal(insts.Instances)
	if err != nil {
		return nil
	}

	err = json.Unmarshal(b, &c.Instances)
	if err != nil {
		return nil
	}

	return c
}

// cache remote server list into local disk and local memory
// @param name service tag
// @param inst instance list of service
// @return error store status
func (lc *LocalCache) Store(name string, inst *CacheInstance) error {
	lc.mutex.Lock()
	defer lc.mutex.Unlock()

	s, ok := lc.instance[name]
	if ok {
		if s.Mils == inst.Mils {
			return nil
		}
	}

	c := &CacheInstance{
		Mils:      inst.Mils,
		Instances: nil,
	}

	b, err := json.Marshal(inst.Instances)
	if err != nil {
		return nil
	}
	err = json.Unmarshal(b, &c.Instances)

	lc.instance[name] = c

	body, err := json.Marshal(lc.instance)
	if err != nil {
		return err
	}

	if !(osutils.Exists(lc.cacheFile)) {
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

	lc.logger.Log("file", "cache.go", "function", "update", "action", "write", "size", size)
	return nil
}

//
// read service list from local disk
//
func (lc *LocalCache) Load() error {
	if !osutils.Exists(lc.cacheFile) {
		return ErrCacheFileExist
	}

	f, err := os.Open(lc.cacheFile)
	if err != nil {
		return err
	}

	content, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	instanceList := make(map[string]*CacheInstance)
	err = json.Unmarshal(content, &instanceList)
	if err != nil {
		return err
	}

	lc.mutex.Lock()
	defer lc.mutex.Unlock()
	lc.instance = instanceList

	return nil
}

func NewLocalCache(dir string, logger log.Logger) (*LocalCache, error) {
	localCache := &LocalCache{
		instance: make(map[string]*CacheInstance),
		logger:   logger,
	}

	if !(osutils.Exists(dir) && osutils.IsDir(dir)) {
		err := osutils.Mkdir(dir, true)
		if err != nil {
			return nil, err
		}
	}

	cacheFile := path.Join(dir, "instance")
	if !(osutils.Exists(cacheFile) && osutils.IsFile(cacheFile)) {
		err := osutils.Touch(cacheFile)
		if err != nil {
			return nil, err
		}
	}

	localCache.cacheFile = cacheFile
	return localCache, nil
}
