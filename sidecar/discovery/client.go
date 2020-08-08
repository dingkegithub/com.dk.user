package discovery

type RegisterCenterClient interface {
	Register(svc *ServiceMeta) error

	Deregister(svc *ServiceMeta) error

	GetServices(svc *ServiceMeta) ([]*ServiceMeta, error)

	Watch(svc *ServiceMeta) <-chan []*ServiceMeta

	CancelWatch(svc *ServiceMeta)

	Close()
}
