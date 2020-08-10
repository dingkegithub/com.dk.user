package config


type CfgListenerFunc func(key string, value interface{})


type CfgCenterClient interface {

	Register(key string, f CfgListenerFunc)

	Deregister(key string)

	GetNamespace(ns string) map[string]interface{}

	Get(key string, ns string) string

	Close()
}
