package http

import (
	"github.com/dingkegithub/com.dk.user/sidecar/discovery"
	"github.com/go-kit/kit/log"
	"os"
	"testing"
)

func TestNewLocalCache(t *testing.T) {
	logger := log.NewLogfmtLogger(os.Stderr)
	cache, err := NewLocalCache("/tmp/discover", logger)
	if err != nil {
		t.Fatalf("new cache failed %s", err.Error())
	}

	instanceList := []*discovery.ServiceMeta {
		&discovery.ServiceMeta{
			Ip:          "a.com",
			Port:        80,
			SvcName: "Svc.1",
			Cluster: "",
			Meta: map[string]interface{}{"a": 1, "c": "dd"},
			Weight:      0,
		},
		&discovery.ServiceMeta{
			Ip:          "b.com",
			Port:        81,
			SvcName: "Svc.1",
			Cluster: "",
			Meta:    nil,
			Weight:      0,
		},
	}

	instanceList2 := []*discovery.ServiceMeta {
		&discovery.ServiceMeta{
			Ip:          "c.com",
			Port:        8080,
			SvcName: "Svc.2",
			Cluster: "",
			Meta:    nil,
			Weight:      0,
		},
	}

	ci1 := &CacheInstance{
		Mils:      0,
		Instances: instanceList,
	}

	ci2 := &CacheInstance{
		Mils:      0,
		Instances: instanceList2,
	}

	err = cache.Store("svc.1", ci1)
	if err != nil {
		t.Fatalf("store err %s", err.Error())
	}

	err = cache.Store("svc.2", ci2)
	if err != nil {
		t.Fatalf("store svc.2 err %s", err.Error())
	}

	res := cache.Instance("svc.2")
	if res == nil {
		t.Fatalf("could not load svc.2")
	}

	if res.Instances[0].Ip != instanceList2[0].Ip && res.Instances[0].Port != instanceList2[0].Port {
		t.Failed()
	}
}
