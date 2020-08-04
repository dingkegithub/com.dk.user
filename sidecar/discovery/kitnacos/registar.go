package kitnacos

import (
	"github.com/dingkegithub/com.dk.user/sidecar/discovery"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"
)

type Registrar struct {
	cli     Client
	service *discovery.ServiceMeta
	logger  log.Logger
}

func (r *Registrar) Register() {
	err := r.cli.Register(r.service)
	if err != nil {
		r.logger.Log("file", "registar.go", "function", "Register", "service", r.service.SvcName, "error", err)
	} else {
		r.logger.Log("file", "registar.go", "function", "Register", "service", r.service.SvcName, "register", "ok")
	}
}

func (r *Registrar) Deregister() {
	err := r.cli.Deregister(r.service)
	if err != nil {
		r.logger.Log("file", "registar.go", "function", "Deregister", "service", r.service.SvcName, "error", err)
	} else {
		r.logger.Log("file", "registar.go", "function", "Deregister", "service", r.service.SvcName, "register", "ok")
	}
}

func NewRegistrar(cli Client, svc *discovery.ServiceMeta, logger log.Logger) (sd.Registrar, error) {
	if cli == nil {
		return nil, ErrorParam
	}

	return &Registrar {
		cli:     cli,
		service: svc,
		logger:  logger,
	}, nil
}
