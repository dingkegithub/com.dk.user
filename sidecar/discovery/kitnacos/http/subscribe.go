package http

import (
	"github.com/dingkegithub/com.dk.user/sidecar/discovery"
	"github.com/dingkegithub/com.dk.user/utils/logging"
	"github.com/modern-go/concurrent"
	"strings"
	"sync"
	"time"
)

// listener function interface
type SubscribeFunc func(uint64, string, []*discovery.ServiceMeta)

type Subscribe struct {
	interval  time.Duration   // timer interval second
	listener  *concurrent.Map // store registered listener
	cacheMile *concurrent.Map // stoer service latest mils
	cmd       *Cmd            // register compoment client
	logger    logging.Logger  // log interface
	signal    chan struct{}   // close signal
	mutex     sync.Mutex      // protect linstener and cache Mile
}

// @param interval interval of timer, second
// @param logger need implement interface Log(kv ...interface{})
// @param cmd nacos client api
func NewSubscribe(interval uint64, logger logging.Logger, cmd *Cmd) *Subscribe {
	s := &Subscribe{
		time.Duration(interval),
		concurrent.NewMap(),
		concurrent.NewMap(),
		cmd,
		logger,
		make(chan struct{}),
		sync.Mutex{},
	}
	go s.Cron()
	return s
}

// user register listener
func (s *Subscribe) Register(name string, mils uint64, f SubscribeFunc) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.cacheMile.Store(name, mils)
	s.listener.Store(name, f)
}

// user deregister listener
func (s *Subscribe) Deregister(name string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.listener.Delete(name)
	s.cacheMile.Delete(name)
}

// close client
func (s *Subscribe) Close() {
	s.signal <- struct{}{}
}

// timer check remote list and close signal
func (s *Subscribe) Cron() {
	ticker := time.Tick(s.interval * time.Second)

	for {
		select {
		case <-ticker:
			s.sync()

		case <-s.signal:
			s.logger.Log("file", "subscribe.go",
				"func", "Cron",
				"action", "signal break")
			break
		}
	}
}

// read remote service's instance list
// keep local cache version is latest
func (s *Subscribe) sync() {

	s.listener.Range(func(key, f interface{}) bool {
		svcId := key.(string)

		keys := strings.Split(svcId, "-")
		if len(keys) != 3 {
			s.logger.Log("file", "subscribe.go",
				"function", "sync",
				"action", "split",
				"svcId", svcId)
			return true
		}

		svcName := keys[0]
		cluster := keys[1]
		group := keys[2]

		// 1. pull remote instance list
		remoteList, err := s.cmd.CmdListInstance(&ListInstanceRequest{
			NamespaceId: "",
			ClusterName: cluster,
			ServiceName: svcName,
			GroupName:   group,
			HealthyOnly: true,
		})

		if err != nil {
			s.logger.Log("file", "subscribe.go",
				"function", "sync",
				"svcId", svcId,
				"action", "request service list",
				"error", err.Error())
			return true
		}

		// 2. read latest cache time
		s.mutex.Lock()
		localCacheMillis, ok := s.cacheMile.Load(svcId)
		if !ok {
			s.cacheMile.Store(svcId, 0)
			localCacheMillis = 0
		}
		s.mutex.Unlock()

		instMetaList := make([]*discovery.ServiceMeta, 0, len(remoteList.Hosts))

		// 2. Is all of instance of service down?
		if len(remoteList.Hosts) == 0 {
			s.logger.Log("file", "subscribe.go",
				"function", "sync",
				"action", "service list is update to 0")

			// 2.1. notify blank list to listener
			f.(SubscribeFunc)(0, svcId, instMetaList)
			s.mutex.Lock()
			s.cacheMile.Store(svcId, uint64(0))
			s.mutex.Unlock()
			return true
		}

		s.logger.Log("file", "subscribe.go",
			"function", "sync",
			"action", "service list is update",
			"remote_millis", remoteList.CacheMillis,
			"local_millis", localCacheMillis)

		// 3. whether remote list newer than local cache
		if remoteList.CacheMillis > localCacheMillis.(uint64) {
			for _, v := range remoteList.Hosts {
				m := &discovery.ServiceMeta{
					Ver:     "",
					Group:   group,
					Cluster: cluster,
					Idc:     "",
					Weight:  0,
					Tag:     "",
					Ip:      v.Ip,
					Port:    v.Port,
					SvcName: svcName,
					SvcId:   "",
					Check:   true,
					Healthy: "",
					Meta:    v.Metadata,
				}
				instMetaList = append(instMetaList, m)
			}

			// 3.1. update local time
			s.cacheMile.Store(key, remoteList.CacheMillis)

			// 3.2. notify latest instance list to listener
			f.(SubscribeFunc)(remoteList.CacheMillis, svcId, instMetaList)
		} else {
			s.logger.Log("file", "subscribe.go",
				"function", "sync",
				"action", "instance list doesn't need update",
				"svcId", svcId,
				"remote_millis", remoteList.CacheMillis,
				"local_millis", localCacheMillis)
		}

		return true
	})
}