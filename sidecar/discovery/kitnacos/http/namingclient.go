package http

type NamingClient interface {
	// 注册一个实例
	RegisterInstance(request *RegisterInstanceRequest) error

	// 注销一个实例
	DeregisterInstance(request *DeregisterInstanceRequest) error

	// 健康实例列表
	HealthyInstances(request *ListInstanceRequest) (*ListInstanceResponse, error)

	// 实例详情
	DetailOfInstance(request *DetailInstanceRequest) (*DetailInstanceResponse, error)
}
