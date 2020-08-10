package apollo

import (
	"encoding/json"
	"github.com/dingkegithub/com.dk.user/sidecar/config"
	"github.com/dingkegithub/com.dk.user/utils/logging"
	"github.com/modern-go/concurrent"
	"github.com/shima-park/agollo"
)

type logAdapter struct {
	logger logging.Logger
}

func (cl logAdapter) Log(kvs ...interface{}) {
	_ = cl.logger.Log(kvs...)
}

type NotifyFunc func(key string, value interface{})


type apolloCfgCenterClient struct {
	opts          *apolloOptions
	cfUrl         string
	appId         string
	signal        chan struct{}
	logger        logging.Logger
	apolloCli     agollo.Agollo
	registerTable *concurrent.Map
}

func (s *apolloCfgCenterClient) GetNamespace(ns string) map[string]interface{} {
	return s.apolloCli.GetNameSpace(ns)
}

func (s *apolloCfgCenterClient) Get(key string, ns string) string {
	if ns == "" {
		return s.apolloCli.Get(key)
	}

	return s.apolloCli.Get(key, agollo.WithNamespace(ns))
}

func NewApolloCfgCenterClient(
	addr string,
	appId string,
	logger logging.Logger,
	opts... Option) (config.CfgCenterClient, error) {

	apolloOptions := newApolloOptions()
	apolloOptions.apply(opts...)

	cfgLog := &logAdapter{
		logger: logger,
	}

	apolloCli, err := agollo.New(addr, appId,
		agollo.Cluster(apolloOptions.Cluster),
		agollo.WithLogger(cfgLog),
		agollo.AutoFetchOnCacheMiss(),
		agollo.BackupFile(apolloOptions.CacheFile))

	if err != nil {
		logger.Log("file", "apolloclient.go",
			"func", "NewApolloCfgCenterClient",
			"msg", "init apollo client",
			"error", err)

		return nil, err
	}

	errCh := apolloCli.Start()

	watchCh := apolloCli.Watch()

	cli := &apolloCfgCenterClient{
		opts:          apolloOptions,
		cfUrl:         addr,
		signal:        make(chan struct{}),
		appId:         appId,
		logger:        logger,
		apolloCli:     apolloCli,
		registerTable: concurrent.NewMap(),
	}
	go cli.listener(errCh, watchCh)

	return cli, nil
}



func (s *apolloCfgCenterClient) Close() {
	s.signal <- struct{}{}
	<- s.signal
	s.apolloCli.Stop()
}

func (s *apolloCfgCenterClient) Register(key string, f config.CfgListenerFunc) {
	s.registerTable.Store(key, f)
}

func (s *apolloCfgCenterClient) Deregister(key string) {
	s.registerTable.Delete(key)
}

func (s *apolloCfgCenterClient) errPollerHandler(err *agollo.LongPollerError) {
	b, er := json.Marshal(err.Notifications)
	if er != nil {
		s.logger.Log("file", "apolloclient.go",
			"func", "errPoller",
			"msg", "marshal Notifications",
			"error", er)
	}

	s.logger.Log("file", "apolloclient.go",
		"func", "listener",
		"ConfigServerURL", err.ConfigServerURL,
		"AppId", err.AppID,
		"Cluster", err.Cluster,
		"Namespace", err.Namespace,
		"Notifications", string(b),
		"Err", err.Err.Error(),
	)
}

func (s *apolloCfgCenterClient) listener(errch <-chan *agollo.LongPollerError, wCh <-chan *agollo.ApolloResponse) {
	for true {
		select {
		case err := <-errch:
			s.errPollerHandler(err)

		case resp := <-wCh:
			s.watcher(resp)

		case <- s.signal:
			s.logger.Log("file", "apolloclient.go",
				"func", "listener",
				"msg", "received close signal")
			break
		}
	}
	s.signal <- struct{}{}
}

func (s *apolloCfgCenterClient) watcher(response *agollo.ApolloResponse) {
	s.logger.Log("file", "apolloclient.go",
		"func", "watcher",
		"msg", "config is update",
		"namespace", response.Namespace,
		"old", response.OldValue,
		"new", response.NewValue,
		"changes", response.Changes,
		"Error", response.Error,
	)

	if response.Error == nil {
		for _, change := range response.Changes {

			if f, ok := s.registerTable.Load(change.Key); ok {
				s.logger.Log("file", "apolloclient.go",
					"func", "watcher",
					"msg", "config changed notify",
					"key", change.Key,
					"val", change.Value,
					"stat", change.Type)
				f.(NotifyFunc)(change.Key, change.Value)
			}
		}
	}
}
