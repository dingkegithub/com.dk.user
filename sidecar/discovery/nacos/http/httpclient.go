package http

import (
	"encoding/json"
	"fmt"
	"github.com/dingkegithub/com.dk.user/sidecar/discovery"
	"github.com/dingkegithub/com.dk.user/utils/netutils"
	"github.com/go-kit/kit/log"
	"sync"
	"time"
)

// nacos http client
type nacosRegisterCenterClient struct {
	cmd       *Cmd
	logger    log.Logger
	cache     *LocalCache
	beatMutex sync.Mutex
	beatCh    map[string]chan struct{}
	subscribe *Subscribe
	subMutex  sync.Mutex
	subChan   map[string]chan []*discovery.ServiceMeta
}

// subscribe callback, subscribe check status of Register Center
// subscribe will invoke callback listener when it found data update
func (n *nacosRegisterCenterClient) listener(mils uint64, name string, m []*discovery.ServiceMeta) {

	ci := &CacheInstance{
		Mils:      mils,
		Instances: nil,
	}

	b, err := json.Marshal(m)
	if err != nil {
		n.logger.Log("file", "httpclient.go",
			"func", "listener",
			"msg", "json marshal meta",
			"error", err)
		return
	}

	err = json.Unmarshal(b, &ci.Instances)
	if err != nil {
		n.logger.Log("file", "httpclient.go",
			"func", "listener",
			"msg", "json unmarshal data",
			"error", err)
		return
	}

	err = n.cache.Store(name, ci)
	if err != nil {
		n.logger.Log("file", "httpclient.go",
			"func", "listener",
			"msg", "store instance list into cache",
			"error", err)
		return
	}

	n.logger.Log("file", "httpclient.go",
		"func", "listener",
		"msg", "push instance",
		"svc", name)
	n.subMutex.Lock()
	defer n.subMutex.Unlock()
	n.subChan[name] <- m
}

// register callback for service
// @param svc service info, need SvcName, Cluster, Group info
// @return chan receive service update info
func (n *nacosRegisterCenterClient) Watch(svc *discovery.ServiceMeta) <-chan []*discovery.ServiceMeta {
	svcId := fmt.Sprintf("%s-%s-%s", svc.SvcName, svc.Cluster, svc.Group)

	n.subMutex.Lock()
	defer n.subMutex.Unlock()

	if _, ok := n.subChan[svcId]; !ok {
		ch := make(chan []*discovery.ServiceMeta)
		n.subChan[svcId] = ch
		n.subscribe.Register(svcId, 0, n.listener)
		n.logger.Log("file", "httpclient.go",
			"func", "Watch",
			"msg", "add watcher",
			"svc", svcId)
	}

	return n.subChan[svcId]
}

// cancel watcher
// @param svc service info, need SvcName, Cluster, Group info
func (n *nacosRegisterCenterClient) CancelWatch(svc *discovery.ServiceMeta) {

	svcId := fmt.Sprintf("%s-%s-%s", svc.SvcName, svc.Cluster, svc.Group)
	n.subMutex.Lock()
	defer n.subMutex.Unlock()
	close(n.subChan[svcId])
	delete(n.subChan, svcId)
	n.subscribe.Deregister(svcId)

}

func (n *nacosRegisterCenterClient) closeSubscribe()  {
	n.subMutex.Lock()
	defer n.subMutex.Unlock()

	for k, v := range n.subChan {
		close(v)
		delete(n.subChan, k)
		n.subscribe.Deregister(k)
	}
	n.subscribe.Close()
}

func (n *nacosRegisterCenterClient) deregister()  {

	n.beatMutex.Lock()
	defer n.beatMutex.Unlock()
	for svcId, ch := range n.beatCh {
		ch <- struct{}{}
		<-ch
		close(ch)
		delete(n.beatCh, svcId)
	}
}

// close register center's client
func (n *nacosRegisterCenterClient) Close() {
	n.closeSubscribe()
	n.deregister()
}

//
// keep service healthy
// @param svc service name
// @param group group name of service
// @param cluster cluster name
// @param ip service ip address
// @param port service port
//
func (n *nacosRegisterCenterClient) heartbeat(svc, cluster, group string, ip string, port uint16) {
	if _, ok := n.beatCh[svc]; ok {
		return
	}

	instanceId := fmt.Sprintf("%s-%s-%d", svc, ip, port)
	ch := make(chan struct{})
	n.beatMutex.Lock()
	n.beatCh[instanceId] = ch
	n.beatMutex.Unlock()

	go func() {
		ticker := time.Tick(3 * time.Second)
		for true {
			select {
			case <-n.beatCh[instanceId]:
				n.beatCh[instanceId] <- struct{}{}
				break

			case <-ticker:
				err := n.cmd.CmdHeartbeatInstance(&HeartbeatRequest{
					ServiceName: svc,
					GroupName:   group,
					Ephemeral:   false,
					Beat: &Beat{
						Cluster:     cluster,
						Ip:          ip,
						Metadata:    nil,
						Port:        port,
						Scheduled:   false,
						ServiceName: svc,
						Weight:      0,
					},
				})
				if err != nil {
					n.logger.Log("file", "httpclient.go",
						"func", "heartbeat",
						"msg", "send heartbeat exception",
						"error", err)
					continue
				}
			}
		}
	}()
}

//
// register service instance
// @param req  service instance information
// @return error register status
//
func (n *nacosRegisterCenterClient) Register(svc *discovery.ServiceMeta) error {

	if err := n.cmd.CmdCreateInstance(&RegisterInstanceRequest{
		Ip:          svc.Ip,
		Port:        svc.Port,
		NamespaceId: "",
		Weight:      svc.Weight,
		Enabled:     true,
		Healthy:     svc.Check,
		Metadata:    svc.MetaString(),
		ClusterName: svc.Cluster,
		ServiceName: svc.SvcName,
		GroupName:   svc.Group,
		Ephemeral:   false,
	}); err != nil {
		return err
	}

	if svc.Check {
		n.heartbeat(svc.SvcName, svc.Group, svc.Cluster, svc.Ip, svc.Port)
	}

	return nil
}

// deregister service instance
// @param req service instance information
// @return error deregister status
func (n *nacosRegisterCenterClient) Deregister(svc *discovery.ServiceMeta) error {
	if err := n.cmd.CmdDeleteInstance(&DeregisterInstanceRequest{
		Ip:          svc.Ip,
		Port:        svc.Port,
		NamespaceId: "",
		ClusterName: svc.Cluster,
		ServiceName: svc.SvcName,
		GroupName:   svc.Group,
		Ephemeral:   false,
	}); err != nil {
		return err
	}

	svcId := fmt.Sprintf("%s-%s-%d", svc.SvcName, svc.Ip, svc.Port)
	if ch, ok := n.beatCh[svcId]; ok {
		ch <- struct{}{}
		<-ch
		n.beatMutex.Lock()
		delete(n.beatCh, svcId)
		n.beatMutex.Unlock()
	}

	return nil
}

// query healthy service list
func (n *nacosRegisterCenterClient) GetServices(svc *discovery.ServiceMeta) ([]*discovery.ServiceMeta, error) {

	n.logger.Log("file", "httpclient.go",
		"func", "GetServices",
		"msg", "service list retrieve",
		"service", svc.SvcName,
		"cluster", svc.Cluster,
		"group", svc.Group)

	l, err := n.cmd.CmdListInstance(&ListInstanceRequest{
		NamespaceId: "",
		ClusterName: svc.Cluster,
		ServiceName: svc.SvcName,
		GroupName:   svc.Group,
		HealthyOnly: true,
	})

	serviceId := fmt.Sprintf("%s-%s-%s", svc.SvcName, svc.Cluster, svc.Group)
	instMetaList := make([]*discovery.ServiceMeta, 0, len(l.Hosts))

	// register center down
	if err != nil {
		if ! n.cache.HasMemCache(serviceId) {
			if err := n.cache.Load(); err != nil {
				n.logger.Log("file", "httpclient.go",
					"func", "GetServices",
					"msg", "load service from local cache",
					"error", err)
				return instMetaList, nil
			}
		}

		instanceList := n.cache.Instance(serviceId)
		if instanceList == nil {
			return instMetaList, nil
		}
		return instanceList.Instances, nil
	}

	ci := &CacheInstance{
		Mils:      l.CacheMillis,
		Instances: nil,
	}

	for _, v := range l.Hosts {
		i := &discovery.ServiceMeta{
			Ver:     svc.Ver,
			Group:   svc.Group,
			Cluster: svc.Cluster,
			Idc:     svc.Idc,
			Weight:  svc.Weight,
			Tag:     svc.Tag,
			Ip:      v.Ip,
			Port:    v.Port,
			SvcName: svc.SvcName,
			SvcId:   svc.SvcId,
			Check:   svc.Check,
			Healthy: svc.Healthy,
			Meta:    v.Metadata,
		}
		instMetaList = append(instMetaList, i)
	}

	ci.Instances = instMetaList

	svcId := fmt.Sprintf("%s-%s-%s", svc.SvcName, svc.Cluster, svc.Group)
	n.cache.Store(svcId, ci)

	return ci.Instances, nil
}

func NewDefaultClient(cacheDir string, logger log.Logger, cm *netutils.ClusterNodeManager) (discovery.RegisterCenterClient, error) {

	cache, err := NewLocalCache(cacheDir, logger)
	if err != nil {
		logger.Log("file", "httpclient.go",
			                "func", "NewDefaultNamingClient",
			                "msg", "new local cache",
			                "error", err)
		return nil, err
	}

	cmd := &Cmd{
		logger: logger,
		nm:     cm,
	}

	s := NewSubscribe(10, logger, cmd)
	namingClient := &nacosRegisterCenterClient{
		logger:    logger,
		cache:     cache,
		beatCh:    make(map[string]chan struct{}),
		subscribe: s,
		cmd:       cmd,
		subChan:   make(map[string]chan []*discovery.ServiceMeta, 10),
	}

	return namingClient, nil
}
