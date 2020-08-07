package kitnacos

import (
	"github.com/dingkegithub/com.dk.user/sidecar/discovery"
	"github.com/dingkegithub/com.dk.user/sidecar/discovery/kitnacos/internal/instance"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"
	"time"
)

type Instancer struct {
	cli    discovery.Client
	logger log.Logger
	cache  *instance.Cache
	quitC  chan struct{}
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

func (i *Instancer) loop(ev <-chan []*discovery.ServiceMeta) {
	i.logger.Log("file", "instancer.go",
		"function", "loop",
		"action", "loop start")

	tick := time.Tick(10 * time.Second)

	for {
		select {
		case msg := <-ev:
			i.logger.Log("file", "instancer.go",
				"function", "loop",
				"action", "update")
			instances := make([]string, 0, len(msg))
			for _, v := range msg {
				instances = append(instances, v.String())
			}
			i.cache.Update(sd.Event{
				Instances: instances,
				Err:       nil,
			})

		case <-i.quitC:
			i.logger.Log("file", "instancer.go",
				"function", "loop",
				"action", "quit")
			return

		case <-tick:
			i.logger.Log("file", "instancer.go",
				"function", "loop",
				"action", "loop wait is running")
			continue
		}
	}
}

func NewInstancer(cli discovery.Client, meta *discovery.ServiceMeta, logger log.Logger) (sd.Instancer, error) {
	inst := &Instancer{
		cache:  instance.NewCache(),
		cli:    cli,
		logger: logger,
		quitC:  make(chan struct{}),
	}

	instList, err := cli.GetServices(meta)
	if err != nil {
		return nil, err
	}

	logger.Log("file", "instancer.go",
		"function", "NewInstancer",
		"action", "GetEntity",
		"instance", len(instList),
		"error", err)

	instances := make([]string, 0, len(instList))
	for _, v := range instList {
		instances = append(instances, v.String())
	}

	inst.cache.Update(sd.Event{
		Instances: instances,
		Err:       err,
	})

	evc := cli.Watch(meta)
	go inst.loop(evc)

	return inst, nil
}