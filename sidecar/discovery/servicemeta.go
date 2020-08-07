package discovery

import (
	"encoding/json"
)

type InstanceMeta struct {

}

type ServiceMeta struct {
	Ver     string                 `json:"ver"`      // service instance version
	Group   string                 `json:"group"`    // service instance group
	Cluster string                 `json:"cluster"`  // cluster that service instance had been released
	Idc     string                 `json:"idc"`      // IDC, maybe service instance is distributed on multiple idc
	Weight  float64                `json:"weight"`   // service instance weight
	Tag     string                 `json:"tag"`      // service instance specify tag
	Ip      string                 `json:"ip"`       // service instance ip address
	Port    uint16                 `json:"port"`     // service instance port
	SvcName string                 `json:"svc_name"` // service name
	SvcId   string                 `json:"svc_id"`   // service id
	Check   bool                   `json:"check"`    // whether need healthy check
	Healthy string                 `json:"healthy"`  // healthy check url
	Meta    map[string]interface{} `json:"meta"`     // other meta info
}

func (sm ServiceMeta) String() string {
	b, _ := json.Marshal(sm)
	return string(b)
}

func (sm ServiceMeta) MetaString() string {
	b, _ := json.Marshal(sm.Meta)
	return string(b)
}