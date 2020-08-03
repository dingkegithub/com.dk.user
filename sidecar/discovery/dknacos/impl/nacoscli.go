package impl

import (
	"com.dk.user/sidecar/discovery"
	"com.dk.user/sidecar/discovery/dknacos"
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/model"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"sync"
)

type NacosNativeCli struct {
	mutex    sync.Mutex
	cli      naming_client.INamingClient
	services sync.Map
}

func NewNacosNativeClient() dknacos.Client {
	cliCfg := constant.ClientConfig{
		TimeoutMs:      10 * 1000,
		ListenInterval: 30 * 1000,
		BeatInterval:   5 * 1000,
		NamespaceId:    "public",
		CacheDir:       "/Users/dk/github/com.dk.micro/com.gz.heartbeat/kitnacos/cache",
		LogDir:         "/Users/dk/github/com.dk.micro/com.gz.heartbeat/kitnacos/log",
		Username:       "",
		Password:       "",
	}

	serverConfigs := []constant.ServerConfig{
		{
			IpAddr:      "127.0.0.1",
			ContextPath: "/nacos",
			Port:        8848,
		},
	}

	namingClient, err := clients.CreateNamingClient(map[string]interface{}{
		"serverConfigs": serverConfigs,
		"clientConfig":  cliCfg,
	})

	if err != nil {
		panic("naming client create failed")
	}

	return &NacosNativeCli{
		cli:      namingClient,
	}
}

func (n *NacosNativeCli) Register(meta *discovery.ServiceMeta) bool {
	_, err := n.cli.RegisterInstance(vo.RegisterInstanceParam{
		Ip:          meta.Ip,
		Port:        uint64(meta.Port),
		Weight:      float64(meta.Weight),
		Enable:      true,
		Healthy:     meta.Check,
		ClusterName: meta.ClusterName,
		ServiceName: meta.SvcName,
	})

	if err != nil {
		fmt.Println("register server failed ", err.Error())
		return false
	}

	return true
}

func (n *NacosNativeCli) Deregister(meta *discovery.ServiceMeta) bool {
	_, err := n.cli.DeregisterInstance(vo.DeregisterInstanceParam{
		Ip:          meta.Ip,
		Port:        uint64(meta.Port),
		Cluster:     meta.ClusterName,
		ServiceName: meta.SvcName,
	})

	if err != nil {
		fmt.Println("deregister logic failed: ", err.Error())
		return false
	}

	return true
}

func (n *NacosNativeCli) Service(meta *discovery.ServiceMeta) []interface{} {
	if svcList, ok := n.services.Load(meta.SvcName); ok {
		return svcList.([]interface{})
	}

	n.mutex.Lock()
	defer n.mutex.Unlock()

	if svcList, ok := n.services.Load(meta.SvcName); ok {
		return svcList.([]interface{})
	}

	go func() {

	}()

	svcList, err := n.cli.SelectInstances(vo.SelectInstancesParam{
		Clusters:    []string{meta.ClusterName,},
		ServiceName: meta.SvcName,
		GroupName:   meta.Group,
		HealthyOnly: true,
	})

	instances := make([]interface{}, 0, len(svcList))
	for _, svc := range svcList {
		instances = append(instances, svc)
	}

	n.services.Store(meta.SvcName, instances)

	err = n.cli.Subscribe(&vo.SubscribeParam{
		ServiceName:       meta.SvcName,
		Clusters:          []string{meta.ClusterName, },
		GroupName:         meta.Group,
		SubscribeCallback: n.Watch,
	})
	if err != nil {
		fmt.Println("subscribe kitnacos failed: ", err.Error())
		return nil
	}

	return instances
}

func (n *NacosNativeCli) Watch(svcList []model.SubscribeService, err error) {
	if err != nil {
		fmt.Println("watch callback err: ", err.Error())
	}

	groupSvc := make(map[string][]interface{})

	for _, svc := range svcList {
		instance := &model.Instance{
			Valid:       svc.Valid,
			InstanceId:  svc.InstanceId,
			Port:        svc.Port,
			Ip:          svc.Ip,
			Weight:      svc.Weight,
			Metadata:    svc.Metadata,
			ClusterName: svc.ClusterName,
			ServiceName: svc.ServiceName,
			Enable:      svc.Enable,
		}
		groupSvc[svc.ServiceName] = append(groupSvc[svc.ServiceName], instance)
	}

	n.mutex.Lock()
	defer n.mutex.Unlock()
	for svcName, svcList := range groupSvc {
		n.services.Store(svcName, svcList)
	}
}
