package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/dingkegithub/com.dk.user/logic/common"
	"github.com/dingkegithub/com.dk.user/logic/common/localcfg"
	usrep "github.com/dingkegithub/com.dk.user/logic/endpoints"
	"github.com/dingkegithub/com.dk.user/logic/proxymiddleware"
	"github.com/dingkegithub/com.dk.user/logic/service/impl"
	usrtr "github.com/dingkegithub/com.dk.user/logic/transports"
	"github.com/dingkegithub/com.dk.user/sidecar/discovery"
	"github.com/dingkegithub/com.dk.user/sidecar/discovery/kitnacos"
	nacoshttp "github.com/dingkegithub/com.dk.user/sidecar/discovery/nacos/http"
	"github.com/dingkegithub/com.dk.user/utils/logging"
	"github.com/dingkegithub/com.dk.user/utils/netutils"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/zap"
	"github.com/go-kit/kit/sd"
	"go.uber.org/zap/zapcore"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)


func main() {
	host := flag.String("host", "127.0.0.1", "-h x.x.x.x or --host=x.x.x.x")
	port := flag.Uint("port", 8081, "-p 8081 --port=8081")
	weight := flag.Float64("weight", 0, "-w 100 --weight=100")
	group := flag.String("group", "default", "-g usrLogic --group=usrLogic")
	clusterName := flag.String("cluster", "default", "-c usrCluster --cluster=usrCluster")
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

	clusterNodeManager, err := netutils.NewClusterNodeManager(5, logger, "127.0.0.1:8848")
	if err != nil {
		logger.Log("file", "main.go",
			"function", "main",
			"action", "init cluster manager",
			"error", err)
		panic(err)
	}

	nacosCli, err := nacoshttp.NewDefaultClient(
		"/tmp/discovery/logic", logger, clusterNodeManager)
	if err != nil {
		logger.Log("file", "main.go",
			"function", "main",
			"register_cli", err,
			"status", "panic exit")
		panic(err)
	}

	ctx := context.Background()
	usrSvc := impl.NewUserLogicService(logger)
	usrSvc = proxymiddleware.UserLogicProxy(ctx, logger, nacosCli)(usrSvc)

	endpoints := usrep.NewUserLogicEndpoints(usrSvc)
	route := usrtr.MakeUserSvcHttpHandler(ctx, endpoints, logger)
	http.Handle("/", usrtr.AccessControl(route))

	errChan := make(chan error)

	svcId := fmt.Sprintf("%s-%d", "user", rand.Int())
	svcMeta := &discovery.ServiceMeta{
		Ip:      *host,
		Port:    uint16(*port),
		SvcName: common.ServiceUsrLogicSrv,
		SvcId:   svcId,
		Weight:  *weight,
		Group:   *group,
		Cluster: *clusterName,
		Check:   true,
		Healthy: "",
	}

	logger.Log("file", "main.go",
		"function", "main",
		"service", common.ServiceUsrLogicSrv,
		"status", "registering")
	var discoverCli sd.Registrar
	{
		discoverCli, err = kitnacos.NewRegistrar(nacosCli, svcMeta, logger)
		if err != nil {
			logger.Log("file", "main.go", "function", "main", "register", err, "status", "panic exit")
			panic(err)
		}
		discoverCli.Register()
		defer discoverCli.Deregister()
	}

	logger.Log("file", "main.go",
		"function", "main",
		"service", common.ServiceUsrLogicSrv,
		"status", "service routine starting")
	go func() {
		listenAddr := fmt.Sprintf("%s:%d", *host, *port)
		logger.Log("file", "main.go",
			"function", "main",
			"service", "UserLogicService",
			"status", "start",
			"listen", listenAddr)
		errChan <- http.ListenAndServe(listenAddr, nil)
	}()

	logger.Log(
		"file", "main.go",
		"function", "main",
		"service", common.ServiceUsrLogicSrv,
		"status", "signal routine starting")
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()

	er := <-errChan
	logger.Log("file", "main.go",
		"function", "main",
		"service", common.ServiceUsrLogicSrv,
		"status", "exit",
		"error", er)
}