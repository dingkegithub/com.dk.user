package http

import "testing"

func TestNewLocalCache(t *testing.T) {
	cache, err := NewLocalCache("/tmp/discover")
	if err != nil {
		t.Fatalf("new cache failed %s", err.Error())
	}

	instanceList := []*Instance {
		&Instance{
			Ip:          "a.com",
			Port:        80,
			ServiceName: "Svc.1",
			ClusterName: "",
			Enable:      false,
			InstanceId:  "",
			Metadata:    nil,
			Weight:      0,
		},
		&Instance{
			Ip:          "b.com",
			Port:        81,
			ServiceName: "Svc.1",
			ClusterName: "",
			Enable:      false,
			InstanceId:  "",
			Metadata:    nil,
			Weight:      0,
		},
	}

	instanceList2 := []*Instance {
		&Instance{
			Ip:          "c.com",
			Port:        8080,
			ServiceName: "Svc.2",
			ClusterName: "",
			Enable:      false,
			InstanceId:  "",
			Metadata:    nil,
			Weight:      0,
		},
	}

	err = cache.Store("svc.1", instanceList)
	if err != nil {
		t.Fatalf("store err %s", err.Error())
	}

	err = cache.Store("svc.2", instanceList2)
	if err != nil {
		t.Fatalf("store svc.2 err %s", err.Error())
	}

	res := cache.Instance("svc.2")
	if res == nil {
		t.Fatalf("could not load svc.2")
	}

	if res[0].Ip != instanceList2[0].Ip && res[0].Port != instanceList2[0].Port {
		t.Failed()
	}
}
