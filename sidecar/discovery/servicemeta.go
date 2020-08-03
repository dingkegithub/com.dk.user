package discovery

type ServiceMeta struct {
	// 服务ip
	Ip          string `json:"ip"`

	// 服务端口
	Port        uint16  `json:"port"`

	// 服务名称
	SvcName     string `json:"svc_name"`

	// 服务id
	SvcId       string `json:"svc_id"`

	// 服务权重
	Weight      float64    `json:"weight"`

	// 服务组
	Group       string `json:"group"`

	// 服务所在集群名称
	ClusterName string `json:"cluster_name"`

	// 健康检查是否允许
	Check       bool   `json:"check"`

	// 健康检查路径
	Healthy     string   `json:"healthy"`
}
