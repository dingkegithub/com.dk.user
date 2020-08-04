package http

import (
	"fmt"
	"github.com/dingkegithub/com.dk.user/utils/netutils"
	"github.com/go-kit/kit/log"
	"sync"
	"time"
)



type namingClient struct {
	svcName string
	logger log.Logger
	nm *netutils.ClusterNodeManager
	cache *LocalCache
	mutex sync.Mutex
	beatCh map[string]chan struct{}
}

func (n *namingClient) buildUrl() string {
	nodeUrl, err := n.nm.Random()
	if err != nil {
		n.logger.Log("file", "httpclient.go", "function", "buildUrl", "action", "random cluster node", "error", err)
		return ""
	}

	return fmt.Sprintf("http://%s", nodeUrl)
}


func (n *namingClient) heartbeat(svc, group, cluster string, ip string, port uint16) {
	if _, ok := n.beatCh[svc]; ok {
		return
	}

	instanceId := fmt.Sprintf("%s-%s-%d", svc, ip, port)
	ch := make(chan struct{})
	n.mutex.Lock()
	n.beatCh[instanceId] = ch
	n.mutex.Unlock()

	go func() {
		ticker := time.Tick(3 * time.Second)
		for true {
			select {
			case <-n.beatCh[instanceId]:
				break

			case <-ticker:
				bUrl := n.buildUrl()
				if bUrl == "" {
					continue
				}

				err := CmdHeartbeatInstance(bUrl, &HeartbeatRequest{
					ServiceName: svc,
					GroupName:   group,
					Ephemeral:   false,
					Beat:        &Beat{
						Cluster:     cluster,
						Ip:          ip,
						Metadata:    nil,
						Port:        port,
						Scheduled:   false,
						ServiceName: "",
						Weight:      0,
					},
				})
				if err != nil {
					continue
				}
			}
		}
	}()
}

func (n *namingClient) RegisterInstance(req *RegisterInstanceRequest) error {
	reqUrl := n.buildUrl()
	if reqUrl == "" {
		return ErrNotFoundHealthyNode
	}

	err := CmdCreateInstance(reqUrl, req)
	if err != nil {
		return err
	}

	if req.Healthy {
		n.heartbeat(req.ServiceName, req.GroupName, req.ClusterName, req.Ip, req.Port)
	}

	return nil
}

func (n *namingClient) DeregisterInstance(req *DeregisterInstanceRequest) error {
	reqUrl := n.buildUrl()
	if reqUrl == "" {
		return ErrNotFoundHealthyNode
	}

	err := CmdDeleteInstance(reqUrl, req)
	if err != nil {
		return err
	}

	instanceId := fmt.Sprintf("%s-%s-%d", req.ServiceName, req.ServiceName, req.Port)
	if ch, ok := n.beatCh[instanceId]; ok {
		ch <- struct{}{}
	}

	return nil
}

func (n *namingClient) HealthyInstances(req *ListInstanceRequest) (*ListInstanceResponse, error) {
	reqUrl := n.buildUrl()
	if reqUrl == "" {
		return nil, ErrNotFoundHealthyNode
	}

	req.HealthyOnly = true

	return CmdListInstance(reqUrl, req)
}

func (n *namingClient) DetailOfInstance(req *DetailInstanceRequest) (*DetailInstanceResponse, error) {
	reqUrl := n.buildUrl()
	if reqUrl == "" {
		return nil, ErrNotFoundHealthyNode
	}

	req.HealthyOnly = true

	return CmdDetailInstance(reqUrl, req)
}

func NewDefaultNamingClient(cacheDir string, interval uint64, logger log.Logger, servers... string) (NamingClient, error) {
	clusterNodeManager, err := netutils.NewClusterNodeManager(interval, servers...)
	if err != nil {
		return nil, err
	}

	cache, err := NewLocalCache(cacheDir)
	if err != nil {
		fmt.Println("file", "httpclient.go", "function", "NewDefaultNamingClient", "action", "new local cache", "error", err)
		return nil, err
	}

	namingClient := &namingClient{
		nm: clusterNodeManager,
		logger: logger,
		cache: cache,
		beatCh: make(map[string]chan struct{}),
	}

	return namingClient, nil
}