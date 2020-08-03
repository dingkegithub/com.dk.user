package netutils

import (
	"testing"
)

func TestNewClusterNodeManager(t *testing.T) {
	manager, err := NewClusterNodeManager(3, "127.0.0.1:8848", "127.0.0.1:9988")
	if err != nil {
		t.Error("error: ", err)
		t.FailNow()
	}

	for i:=0; i<20; i++ {
		t.Log("signal thread run")
		randomNode, err := manager.Random()
		if err != nil {
			t.Error("error: ", err)
			t.FailNow()
		}

		t.Log("random: ", randomNode, len(manager.allNodes), len(manager.availableNodes), len(manager.brokenNodes))
	}
}
