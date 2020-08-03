package kitnacos

import (
	"com.dk.user/sidecar/discovery/kitnacos/internal/instance"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"
)

type Instancer struct {
	cli Client
	logger log.Logger
	cache *instance.Cache
	quitC chan struct{}
}

func (i *Instancer) Register(events chan<- sd.Event) {
	i.cache.Register(events)
}

func (i *Instancer) Deregister(events chan<- sd.Event) {
	i.cache.Deregister(events)
}

func (i *Instancer) Stop() {
	close(i.quitC)
}

func (i *Instancer) loop(ev <- chan *NacosEvent) {
	for {
		select {
		case msg := <- ev:
			i.logger.Log("file", "instancer.go", "function", "loop", "action", "update")
			i.cache.Update(sd.Event{
				Instances: msg.Svc,
				Err:       nil,
			})

		case <-i.quitC:
			i.logger.Log("file", "instancer.go", "function", "loop", "action", "quit")
			return
		}
	}
}

func NewInstancer(cli Client, svcName string, clusterName, group string, logger log.Logger) (sd.Instancer, error) {
	inst := &Instancer{
		cache: instance.NewCache(),
		cli: cli,
		logger: logger,
		quitC: make(chan struct{}),
	}

	instList, err := cli.GetEntity(svcName, clusterName, group)
	if err != nil {
		return nil, err
	}

	logger.Log("file", "instancer.go", "function", "NewInstancer", "action", "GetEntity", "instance", len(instList), "error", err)
	inst.cache.Update(sd.Event{
		Instances: instList,
		Err:       err,
	})

	evC := cli.WatchSvc(svcName, clusterName, group)
	go inst.loop(evC)

	return inst, nil
}