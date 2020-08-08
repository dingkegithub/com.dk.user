package kitnacos

import (
	"github.com/dingkegithub/com.dk.user/sidecar/discovery"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"
	"time"
)

type Registrar struct {
	cli     discovery.RegisterCenterClient
	service *discovery.ServiceMeta
	logger  log.Logger
}

func (r *Registrar) Register() {
	go func() {
		for {
			err := r.cli.Register(r.service)
			if err != nil {
				r.logger.Log("file", "registar.go",
					"func", "Register",
					"service", r.service.SvcName,
					"msg", "register failed",
					"error", err)
				time.Sleep(5 * time.Second)
			} else {
				break
			}
		}
		r.logger.Log("f", "registar.go",
			"func", "Register",
			"service", r.service.SvcName,
			"register", "success")
	}()

}

func (r *Registrar) Deregister() {
	err := r.cli.Deregister(r.service)
	if err != nil {
		r.logger.Log("file", "registar.go", "function", "Deregister", "service", r.service.SvcName, "error", err)
	} else {
		r.logger.Log("file", "registar.go", "function", "Deregister", "service", r.service.SvcName, "register", "ok")
	}
}

func NewRegistrar(cli discovery.RegisterCenterClient, svc *discovery.ServiceMeta, logger log.Logger) (sd.Registrar, error) {
	if cli == nil {
		return nil, discovery.ErrorParam
	}

	return &Registrar {
		cli:     cli,
		service: svc,
		logger:  logger,
	}, nil
}
