package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/dingkegithub/com.dk.user/das/common"
	"github.com/dingkegithub/com.dk.user/das/common/localcfg"
	"github.com/dingkegithub/com.dk.user/das/endpoints"
	"github.com/dingkegithub/com.dk.user/das/model"
	"github.com/dingkegithub/com.dk.user/das/proto/userpb"
	"github.com/dingkegithub/com.dk.user/das/service"
	"github.com/dingkegithub/com.dk.user/das/service/impl"
	"github.com/dingkegithub/com.dk.user/das/setting"
	"github.com/dingkegithub/com.dk.user/das/transport"
	"github.com/dingkegithub/com.dk.user/sidecar/config/apollo"
	"github.com/dingkegithub/com.dk.user/sidecar/discovery"
	"github.com/dingkegithub/com.dk.user/sidecar/discovery/kitnacos"
	nacoshttp "github.com/dingkegithub/com.dk.user/sidecar/discovery/nacos/http"
	"github.com/dingkegithub/com.dk.user/utils/logging"
	"github.com/dingkegithub/com.dk.user/utils/netutils"
	"github.com/go-kit/kit/log/zap"
	"github.com/go-kit/kit/sd"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
)



func main() {
	ctx := context.Background()
	errChan := make(chan error)

	cfgCenterAddr := flag.String("cfg_center_addr", "127.0.0.1:8848", "--cfg_center x.x.x.x:x")
	registerCenterAddrs := flag.String("register_centers", "127.0.0.1:8848", "x.x.x.x:x,x.x.x.x:x...")
	appId := flag.String("app_id", "UserDasService", "--app_id application identify name")
	cluster := flag.String("c", "default", "--c cluster name")
	namespace := flag.String("ns", "default", "--ns namespace name")
	host := flag.String("host", "0.0.0.0", "--host x.x.x.x")
	port := flag.Uint("port", 18080, "--port x")
	cfgFile := flag.String("config", "", "--config=xxx.json")

	flag.Parse()

	if *cfgFile == "" {
		panic("need option -c xxx.json or --config=xxx.json")
	}

	_, err := localcfg.NewCfgLoader(*cfgFile)
	if err != nil {
		panic(fmt.Sprintf("config file error %s", err.Error()))
	}

	cfgH := localcfg.GetCfg()

	var logger logging.Logger
	{
		logCfg := cfgH.GetLogCfg()
		zapLogger := logging.LogInit(logCfg.FileName, logCfg.MaxSize, logCfg.MaxBackups, logCfg.MaxAge, logCfg.Level)
		logger = zap.NewZapSugarLogger(zapLogger, zapcore.Level(logCfg.Level-1))
	}

	cfgClient, err := apollo.NewApolloCfgCenterClient(
		*cfgCenterAddr,
		*appId,
		logger,
		apollo.WithCluster(*cluster),
	)
	if err != nil {
		logger.Log("file", "main.go",
			"func", "main",
			"msg", "init config center client failed",
			"addr", *cfgCenterAddr,
			"app", appId,
			"service", common.ServiceName)
		panic(err)
	}

	setting.New(cfgClient)
	model.New(logger)

	var svc service.UserSvc
	{
		svc = impl.NewUserSrv(logger)
	}

	logger.Log("file", "main.go",
		"func", "main",
		"msg", "register service",
		"service", common.ServiceName)
	ends := endpoints.NewUsrEndpoints(svc)
	handler := transport.NewRpcUsrSvc(ctx, ends)

	var discoverCli sd.Registrar
	{
		svcId := fmt.Sprintf("%s-%d", "UserDas", rand.Int())
		svcMeta := &discovery.ServiceMeta{
			Ip:      *host,
			Port:    uint16(*port),
			SvcName: common.ServiceName,
			SvcId:   svcId,
			Weight:  0,
			Group:   "default",
			Cluster: "default",
			Check:   true,
			Healthy: "",
			Meta: map[string]interface{}{
				"idc": "ChongQing",
				"namespace": namespace,
				"need_login": false,
			},
		}

		addrs := strings.Split(*registerCenterAddrs, ",")
		registerClusterNodeManager, err := netutils.NewClusterNodeManager(5, logger, addrs...)
		if err != nil {
			logger.Log("file", "main.go",
				"func", "main",
				"msg", "init register center cluster manager failed",
				"error", err)
			panic(err)
		}

		cacheDir := fmt.Sprintf("/tmp/discovery/.%s", common.ServiceName)
		nacosCli, err := nacoshttp.NewDefaultClient(
			cacheDir, logger, registerClusterNodeManager)
		if err != nil {
			logger.Log("file", "main.go",
				"func", "main",
				"msg", "create discover client(nacos) failed",
				"error", err)
			panic(err)
		}

		discoverCli, err = kitnacos.NewRegistrar(nacosCli, svcMeta, logger)
		if err != nil {
			logger.Log("file", "main.go",
				"func", "main",
				"msg", "init register failed",
				"error", err)
			panic(err)
		}

		discoverCli.Register()
		defer discoverCli.Deregister()
	}

	lisAddr := fmt.Sprintf("%s:%d", *host, *port)
	go func() {
		ls, _ := net.Listen("tcp", lisAddr)
		rpcServer := grpc.NewServer()
		userpb.RegisterUserDasServiceServer(rpcServer, handler)

		logger.Log("file", "main.go",
			"func", "main",
			"msg", "service is start",
			"service", common.ServiceName,
			"listen", lisAddr)
		err := rpcServer.Serve(ls)
		if err != nil {
			logger.Log("file", "main.go",
				"func", "main",
				"msg", "start rpc server failed",
				"service", common.ServiceName,
				"error", err)
			errChan <- fmt.Errorf("%s", err.Error())
		}
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()

	<-errChan
	logger.Log("file", "main.go",
		"func", "main",
		"msg", "service receive signal error and service exit",
		"service", common.ServiceName,
		"listen", lisAddr)
}
