package http

import (
	"encoding/json"
	"testing"
)

func TestCmdCreateInstance(t *testing.T) {
	nacosUrl := "http://127.0.0.1:8848"
	instance := &RegisterInstanceRequest{
		Ip:          "dk.user.com",
		Port:        8080,
		NamespaceId: "public",
		Weight:      1,
		Enabled:     true,
		Healthy:     false,
		Metadata:    "",
		ClusterName: "",
		ServiceName: "UserDasTestService",
		GroupName:   "",
		Ephemeral:   false,
	}
	err := CmdCreateInstance(nacosUrl, instance)
	if err != nil {
		t.Errorf("error: %s", err.Error())
		t.Fail()
	}
}

func TestCmdDeleteInstance(t *testing.T) {
	nacosUrl := "http://127.0.0.1:8848"
	instance := &DeregisterInstanceRequest{
		Ip:          "dk.user.com",
		Port:        8080,
		NamespaceId: "public",
		ClusterName: "",
		ServiceName: "UserDasTestService",
		GroupName:   "",
		Ephemeral:   false,
	}

	err := CmdDeleteInstance(nacosUrl, instance)
	if err != nil {
		t.Errorf("error: %s", err.Error())
		t.Fail()
	}
}

func TestCmdUpdateInstance(t *testing.T) {
	nacosUrl := "http://127.0.0.1:8848"
	instance := &ModifyInstanceRequest{
		Ip:          "dk.user.com",
		Port:        8081,
		ServiceName: "UserDasTestService",
	}

	err := CmdUpdateInstance(nacosUrl, instance)
	if err != nil {
		t.Errorf("error: %s", err.Error())
		t.Fail()
	}
}

func TestCmdListInstance(t *testing.T) {
	nacosUrl := "http://127.0.0.1:8848"
	requestParam := &ListInstanceRequest{
		ServiceName: "UserDasTestService",
	}

	responseParam, err := CmdListInstance(nacosUrl, requestParam)
	if err != nil {
		t.Errorf("error: %s", err.Error())
		t.FailNow()
	}

	info, err := json.Marshal(responseParam)
	if err != nil {
		t.Errorf("error: %s", err.Error())
		t.Failed()
	}

	t.Log("info: ", string(info))
}

func TestCmdDetailInstance(t *testing.T) {
	nacosUrl := "http://127.0.0.1:8848"
	requestParam := &DetailInstanceRequest{
		ServiceName: "UserDasTestService",
		Ip: "dk.user.com",
		Port: 8081,
	}

	responseParam, err := CmdDetailInstance(nacosUrl, requestParam)
	if err != nil {
		t.Errorf("error: %s", err.Error())
		t.FailNow()
	}

	info, err := json.Marshal(responseParam)
	if err != nil {
		t.Errorf("error: %s", err.Error())
		t.Failed()
	}

	t.Log("info: ", string(info))
}

func TestCmdHeartbeatInstance(t *testing.T) {
	nacosUrl := "http://127.0.0.1:8848"
	requestParam := &HeartbeatRequest{
		ServiceName: "UserDasTestService",
		Beat: &Beat{
			Cluster:     "",
			Ip:          "dk.user.com",
			Metadata:    nil,
			Port:        8081,
			Scheduled:   false,
			ServiceName: "UserDasTestService",
			Weight:      0,
		},
	}

	err := CmdHeartbeatInstance(nacosUrl, requestParam)
	if err != nil {
		t.Errorf("error: %s", err.Error())
		t.FailNow()
	}
}