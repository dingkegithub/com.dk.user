package http

import (
	"github.com/dingkegithub/com.dk.user/sidecar/discovery"
	"github.com/dingkegithub/com.dk.user/utils/logging"
	"github.com/modern-go/concurrent"
	"strings"
	"sync"
	"time"
)

type SubscribeFunc func(uint64, string, []*discovery.ServiceMeta)

type Subscribe struct {
	interval  time.Duration
	listener  *concurrent.Map
	cacheMile *concurrent.Map
	cmd       *Cmd
	logger    logging.Logger
	signal    chan struct{}
	mutex     sync.Mutex
}

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

func (s *Subscribe) Register(name string, mils uint64, f SubscribeFunc) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.cacheMile.Store(name, mils)
	s.listener.Store(name, f)
}

func (s *Subscribe) Deregister(name string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.listener.Delete(name)
	s.cacheMile.Delete(name)
}

func (s *Subscribe) Close() {
	s.signal <- struct{}{}
	s.cmd.Close()
}

func (s *Subscribe) Cron() {
	ticker := time.Tick(s.interval * time.Second)

	for {
		select {
		case <-ticker:
			s.sync()

		case <-s.signal:
			s.logger.Log("file", "subscribe.go", "function", "Cron", "action", "signal break")
			break
		}
	}
}

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

		s.mutex.Lock()
		localCacheMillis, ok := s.cacheMile.Load(svcId)
		if !ok {
			s.cacheMile.Store(svcId, 0)
			localCacheMillis = 0
		}
		s.mutex.Unlock()

		instMetaList := make([]*discovery.ServiceMeta, 0, len(remoteList.Hosts))

		if len(remoteList.Hosts) == 0 {
			s.logger.Log("file", "subscribe.go",
				"function", "sync",
				"action", "service list is update to 0")

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

			s.cacheMile.Store(key, remoteList.CacheMillis)
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
