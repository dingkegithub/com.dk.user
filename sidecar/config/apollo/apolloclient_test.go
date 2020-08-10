package apollo

import (
	"github.com/go-kit/kit/log"
	"os"
	"testing"
)

func TestNewSetting(t *testing.T) {
	firstSetting, err := NewApolloCfgCenterClient("127.0.0.1:8080",
		"UserDasServiceId", log.NewLogfmtLogger(os.Stderr))

	secondSetting, err := NewApolloCfgCenterClient("127.0.0.1:8080",
		"UserDasServiceId", log.NewLogfmtLogger(os.Stderr), WithCluster("BiGeYun"))
	if err != nil {
		t.Log("new setting failed", err)
		t.FailNow()
	}

	// application
	cfg := firstSetting.GetNamespace("application")
	t.Log( "namespace:", "application", "cfg", cfg)

	cfg = secondSetting.GetNamespace("application")
	t.Log("namespace:", "application", "cfg", cfg)

	// application
	v := firstSetting.Get("db", "application")
	t.Log("namespace:", "application", "db", v)

	v = secondSetting.Get("db", "application")
	t.Log("namespace:", "application", "db", v)

	// registerTable
	v = firstSetting.Get("registerTable", "application")
	t.Log("namespace:", "application", "registerTable", v)

	v = secondSetting.Get("registerTable", "application")
	t.Log("namespace:", "application", "registerTable", v)

	// namespace
	cfg = firstSetting.GetNamespace("Tech-2.basic-middleware")
	t.Log("namespace:", "Tech-2.basic-middleware", "cfg", cfg)

	cfg = secondSetting.GetNamespace("Tech-2.basic-middleware")
	t.Log( "namespace:", "Tech-2.basic-middleware", "cfg", cfg)

	// register centre
	v = firstSetting.Get("register-centre", "Tech-2.basic-middleware")
	t.Log("namespace:", "Tech-2.basic-middleware", "register-centre", v)

	v = secondSetting.Get("register-centre", "Tech-2.basic-middleware")
	t.Log("namespace:", "Tech-2.basic-middleware", "register-centre", v)

}