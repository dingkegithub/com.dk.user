package kitnacos

import (
	"com.dk.user/sidecar/discovery"
	nacoshttp "com.dk.user/sidecar/discovery/kitnacos/http"
	"encoding/json"
	"github.com/go-kit/kit/log"
	"github.com/nacos-group/nacos-sdk-go/model"
	"sync"
)

type Client interface {
	Register(svc *discovery.ServiceMeta) error

	Deregister(svc *discovery.ServiceMeta) error

	GetEntity(svcName string, clusterName string, group string) ([]string, error)

	WatchSvc(svcName, cluster, group string) <-chan *NacosEvent

	Stop()
}

type NacosEvent struct {
	Svc []string
}

type Listener struct {
	Cluster []string
	Group   string
	SvcName string
}

type defaultClient struct {
	cli      nacoshttp.NamingClient
	//cliCfg   constant.ClientConfig
	ev       chan *NacosEvent
	listener map[string]*Listener
	logger   log.Logger
	mutex sync.Mutex
	beatCh  map[string]chan struct{}
}

func (d *defaultClient) Register(svc *discovery.ServiceMeta) error {

	//
	//registerParam  := vo.RegisterInstanceParam{
	//	Ip:          svc.Ip,
	//	Port:        svc.Port,
	//	Weight:      0,
	//	Enable:      true,
	//	Healthy:     svc.Check,
	//	Metadata:    nil,
	//	ClusterName: svc.ClusterName,
	//	ServiceName: svc.SvcName,
	//	//GroupName:   svc.Group,
	//	Ephemeral:   false,
	//}
	//
	//ok, err := d.cli.RegisterInstance(registerParam)
	//if ok {
	//	d.logger.Log("file", "client.go", "function", "Register", "service", registerParam.ServiceName, "ip", svc.Ip, "error", err)
	//	return nil
	//}
	registerParam := &nacoshttp.RegisterInstanceRequest{
		Ip:          svc.Ip,
		Port:        svc.Port,
		NamespaceId: "",
		Weight:      svc.Weight,
		Enabled:     true,
		Healthy:     svc.Check,
		Metadata:    "",
		ClusterName: svc.ClusterName,
		ServiceName: svc.SvcName,
		GroupName:   svc.Group,
		Ephemeral:   false,
	}
	err := d.cli.RegisterInstance(registerParam)
	d.logger.Log("file", "client.go", "function", "Register", "service", registerParam.ServiceName, "ip", svc.Ip, "error", err)
	return err
}

func (d *defaultClient) Deregister(svc *discovery.ServiceMeta) (err error) {
	d.logger.Log("file", "client.go", "function", "Register", "service", svc.SvcName, "ip", svc.Ip)

	//if _, ok := d.listener[svc.SvcName]; ok {
	//	_ = d.cli.Unsubscribe(&vo.SubscribeParam{
	//		ServiceName:       svc.SvcName,
	//		Clusters:          []string{svc.ClusterName, },
	//		GroupName:         svc.Group,
	//	})
	//
	//	delete(d.listener, svc.SvcName)
	//}

	//_, err = d.cli.DeregisterInstance(vo.DeregisterInstanceParam{
	//	Ip:          svc.Ip,
	//	Port:        svc.Port,
	//	GroupName:   svc.Group,
	//	Cluster:     svc.ClusterName,
	//	ServiceName: svc.SvcName,
	//})
	//
	err = d.cli.DeregisterInstance(&nacoshttp.DeregisterInstanceRequest{
		Ip:          svc.Ip,
		Port:        svc.Port,
		NamespaceId: "",
		ClusterName: svc.ClusterName,
		ServiceName: svc.SvcName,
		GroupName:   svc.Group,
		Ephemeral:   false,
	})
	return
}

func (d *defaultClient) GetEntity(svcName string, clusterName string, group string) ([]string, error) {
	//insts, err := d.cli.SelectInstances(vo.SelectInstancesParam{
	//	ServiceName: svcName,
	//	Clusters:    []string{clusterName},
	//	GroupName:   group,
	//	HealthyOnly: true,
	//})

	lResp, err := d.cli.HealthyInstances(&nacoshttp.ListInstanceRequest{
		NamespaceId: "",
		ClusterName: clusterName,
		ServiceName: svcName,
		GroupName:   group,
		HealthyOnly: true,
	})
	d.logger.Log("file", "client.go", "function", "GetEntity", "svcName", svcName, "error", err)

	if err != nil {
		return nil, err
	}

	if len(lResp.Hosts) == 0 {
		return []string{}, nil
	}

	instList := make([]string, 0, len(lResp.Hosts))
	for _, ins := range lResp.Hosts {
		svcMeta := &discovery.ServiceMeta{
			Ip:          ins.Ip,
			Port:        ins.Port,
			SvcName:     svcName,
			SvcId:       "",
			Weight:      ins.Weight,
			Group:       "",
			ClusterName: clusterName,
			Check:       false,
			Healthy:     "",
		}
		metaStr, err := json.Marshal(svcMeta)
		if err != nil {
			d.logger.Log("file", "client.go", "function", "GetEntity", "action", "json marshal", "error", err)
			continue
		}
		instList = append(instList, string(metaStr))
	}

	return instList, nil
}

func (d *defaultClient) WatchSvc(svcName, cluster, group string) <-chan *NacosEvent {

	//if _, ok := d.listener[svcName]; !ok {
	//	d.logger.Log("file", "client.go", "function", "WatchSvc", "action", "subscribe", "service", svcName, "cluster", cluster, "group", group)
	//	_ = d.cli.Subscribe(&vo.SubscribeParam{
	//		ServiceName:       svcName,
	//		Clusters:          []string{cluster, },
	//		//GroupName:         group,
	//		SubscribeCallback: d.Watch,
	//	})
	//
	//	d.listener[svcName] = &Listener{
	//		Cluster: []string{cluster},
	//		//Group:   group,
	//		SvcName: svcName,
	//	}
	//}

	return d.ev
}

func (d *defaultClient) Watch(svcList []model.SubscribeService, err error) {
	//if err != nil {
	//	d.logger.Log("file", "client.go", "function", "Watch", "action", "active watch", "error", err)
	//	return
	//}
	//
	//instList := make([]string, 0, len(svcList))
	//d.logger.Log("file", "client.go", "function", "Watch", "action", "active watch", "size", len(svcList))
	//
	//for _, ins := range svcList {
	//	svcMeta := &discovery.ServiceMeta{
	//		Ip:          ins.Ip,
	//		Port:        ins.Port,
	//		SvcName:     ins.ServiceName,
	//		SvcId:       "",
	//		Weight:      ins.Weight,
	//		//Group:       "",
	//		ClusterName: ins.ClusterName,
	//		Check:       false,
	//		Healthy:     "",
	//	}
	//	metaStr, err := json.Marshal(svcMeta)
	//	if err != nil {
	//		d.logger.Log("file", "client.go", "function", "watch", "action", "json marshal", "err", err)
	//		continue
	//	}
	//	d.logger.Log("file", "client.go", "function", "watch", "status", "service update", "info", string(metaStr))
	//	instList = append(instList, string(metaStr))
	//}
	//
	//d.logger.Log("file", "client.go", "function", "watch", "status", "notify update by chan", "info", len(instList))
	//d.ev <- &NacosEvent{Svc: instList}
}

func (d *defaultClient) Stop() {
	d.logger.Log("kitnacos", "stop")

	for _, v := range d.listener {
		d.logger.Log("svc_status", "unsubscribe", "svc", v.SvcName, "cluster", v.Cluster, "group", v.Group)
		//d.cli.Unsubscribe(&vo.SubscribeParam{
		//	ServiceName: v.SvcName,
		//	Clusters:    v.Cluster,
		//	//GroupName:         v.Group,
		//	SubscribeCallback: nil,
		//})
	}
}

func NewDefaultClient(interval uint64, catchLogBase string, logger log.Logger, servers ...string) (Client, error) {
	//func NewDefaultClient(namespace string, catchLogBase string, host string, port uint64, logger log.Logger) (Client, error)  {
	//ns := "public"
	//if namespace != "" {
	//	ns = namespace
	//}

	//var logDir string
	//var catchDir string
	//{
	//	logDir = path.Join(catchLogBase, "log")
	//	catchDir = path.Join(catchLogBase, "catch")
	//
	//	if ok := osutils.IsDir(catchLogBase); !ok {
	//		err := osutils.Mkdir(catchLogBase, true)
	//		if err != nil {
	//			return nil, err
	//		}
	//	}
	//
	//	_ = osutils.Mkdir(logDir, false)
	//	_ = osutils.Mkdir(catchDir, false)
	//}

	//cliCfg := constant.ClientConfig{
	//	TimeoutMs:      10 * 1000,
	//	ListenInterval: 30 * 1000,
	//	BeatInterval:   5 * 1000,
	//	NamespaceId:    ns,
	//	CacheDir:       catchDir,
	//	LogDir:         logDir,
	//	Username:       "nacos",
	//	Password:       "nacos",
	//}
	//logger.Log("file", "client.go", "function", "NewDefaultClient", "ns", cliCfg.NamespaceId, "logdir", cliCfg.LogDir)

	if len(servers) <= 0 {
		return nil, ErrorParam
	}
	//
	//nacosHost := "127.0.0.1"
	//if host != "" {
	//	nacosHost = host
	//}
	//
	//var nacosPort uint64 = 8848
	//if port <= 0 {
	//	nacosPort = port
	//}
	//
	//serverConfigs := []constant.ServerConfig{
	//	{
	//		IpAddr:      nacosHost,
	//		ContextPath: "/nacos",
	//		Port:        nacosPort,
	//	},
	//}
	//
	//namingClient, err := clients.CreateNamingClient(map[string]interface{} {
	//	constant.KEY_SERVER_CONFIGS: serverConfigs,
	//	constant.KEY_CLIENT_CONFIG: cliCfg,
	//})
	//
	namingClient, err := nacoshttp.NewDefaultNamingClient(catchLogBase, interval, logger, servers...)
	if err != nil {
		return nil, err
	}

	return &defaultClient{
		cli:      namingClient,
		//cliCfg:   cliCfg,
		listener: make(map[string]*Listener),
		ev:       make(chan *NacosEvent),
		logger:   logger,
		beatCh: make(map[string]chan struct{}),
	}, nil
}
