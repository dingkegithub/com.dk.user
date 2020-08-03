package dknacos

import "com.dk.user/sidecar/discovery"

type Client interface {
	/**
	 * 服务注册
	 * @param: 服务注册元信息
	 */
	Register(meta *discovery.ServiceMeta) bool

	/**
	 * 服务注册
	 * @param: 服务注册元信息
	 */
	Deregister(meta *discovery.ServiceMeta) bool

	/**
	 * 服务注册
	 * @param: 服务注册元信息
	 */
	Service(meta *discovery.ServiceMeta) []interface{}
}
