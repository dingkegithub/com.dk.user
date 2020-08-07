package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/dingkegithub/com.dk.user/das/common/localcfg"
	"github.com/dingkegithub/com.dk.user/das/endpoints"
	"github.com/dingkegithub/com.dk.user/das/model"
	"github.com/dingkegithub/com.dk.user/das/proto/userpb"
	"github.com/dingkegithub/com.dk.user/das/service"
	"github.com/dingkegithub/com.dk.user/das/service/impl"
	"github.com/dingkegithub/com.dk.user/das/transport"
	"github.com/dingkegithub/com.dk.user/sidecar/discovery"
	"github.com/dingkegithub/com.dk.user/sidecar/discovery/kitnacos"
	nacoshttp "github.com/dingkegithub/com.dk.user/sidecar/discovery/kitnacos/http"
	"github.com/dingkegithub/com.dk.user/utils/logging"
	"github.com/dingkegithub/com.dk.user/utils/netutils"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/zap"
	"github.com/go-kit/kit/sd"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"syscall"
)

var (
	ServiceName = "userpb.UserDasService"
)


func main() {
	ctx := context.Background()
	errChan := make(chan error)

	host := flag.String("host", "127.0.0.1", "-h x.x.x.x --host=x.x.x.x")
	port := flag.Uint("port", 8080, "-p x or --port=x")
	cfgFile := flag.String("config", "", "-c xxx.json --config=xxx.json")
	flag.Parse()

	if *cfgFile == "" {
		panic("need option -c xxx.json or --config=xxx.json")
	}

	_, err := localcfg.NewCfgLoader(*cfgFile)
	if err != nil {
		panic(fmt.Sprintf("config file error %s", err.Error()))
	}

	cfgH := localcfg.GetCfg()

	var logger log.Logger
	{
		logCfg := cfgH.GetLogCfg()
		zapLogger := logging.LogInit(logCfg.FileName, logCfg.MaxSize, logCfg.MaxBackups, logCfg.MaxAge, logCfg.Level)
		logger = zap.NewZapSugarLogger(zapLogger, zapcore.Level(logCfg.Level-1))
	}

	logger.Log("file", "main.go", "function", "main", "service", ServiceName, "status", "init model")
	model.Init("root", "123456", "user")

	logger.Log("file", "main.go", "function", "main", "service", ServiceName, "status", "init service")
	var svc service.UserSvc
	{
		svc = impl.NewUserSvc()
	}

	logger.Log("file", "main.go", "function", "main", "service", ServiceName, "status", "register service")
	ends := endpoints.NewUsrEndpoints(svc)
	handler := transport.NewRpcUsrSvc(ctx, ends)

	var discoverCli sd.Registrar
	{
		svcId := fmt.Sprintf("%s-%d", "UserDas", rand.Int())
		svcMeta := &discovery.ServiceMeta{
			Ip:      *host,
			Port:    uint16(*port),
			SvcName: ServiceName,
			SvcId:   svcId,
			Weight:  0,
			Group:   "default",
			Cluster: "default",
			Check:   true,
			Healthy: "",
			Meta: map[string]interface{}{
				"idc": "ChongQing",
				"need_login": false,
			},
		}

		clusterNodeManager, err := netutils.NewClusterNodeManager(5, "127.0.0.1:8848")
		if err != nil {
			logger.Log("file", "main.go",
				"function", "main",
				"action", "init cluster manager",
				"error", err)
			panic(err)
		}

		nacosCli, err := nacoshttp.NewDefaultClient(
			"/Users/dk/github/nacos/nacos/mydata", logger, clusterNodeManager)
		if err != nil {
			logger.Log("file", "main.go", "function", "main", "register_cli", err, "status", "panic exit")
			panic(err)
		}

		discoverCli, err = kitnacos.NewRegistrar(nacosCli, svcMeta, logger)
		if err != nil {
			logger.Log("file", "main.go", "function", "main", "register", err, "status", "panic exit")
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

		logger.Log("file", "main.go", "function", "main", "service", "UserDasService", "listen", lisAddr)
		err := rpcServer.Serve(ls)
		if err != nil {
			logger.Log("file", "main.go", "function", "main", "service", "UserDasService", "status", "exit", "error", err)
			errChan <- fmt.Errorf("%s", err.Error())
		}
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()

	<-errChan
	logger.Log("file", "main.go", "function", "main", "service", "UserDasService", "listen", lisAddr)
}
