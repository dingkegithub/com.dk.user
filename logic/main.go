package main

import (
	usrep "com.dk.user/logic/endpoints"
	"com.dk.user/logic/proxymiddleware"
	"com.dk.user/logic/service/impl"
	usrtr "com.dk.user/logic/transports"
	"com.dk.user/sidecar/discovery"
	"com.dk.user/sidecar/discovery/kitnacos"
	"context"
	"flag"
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

const (
	ServiceName = "UserLogicService"
	ServiceDasName = ""
)

func main() {

	var (
		server = "127.0.0.1:8848"
	)

	flag.Parse()
	host := flag.String("host", "127.0.0.1", "-h x.x.x.x or --host=x.x.x.x")
	port := flag.Uint("port", 8081, "-p 8081 --port=8081")
	weight := flag.Float64("weight", 0, "-w 100 --weight=100")
	group := flag.String("group", "default", "-g usrLogic --group=usrLogic")
	clusterName := flag.String("cluster", "default", "-c usrCluster --cluster=usrCluster")

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
	}

	ctx := context.Background()
	usrSvc := impl.NewUserLogicService(logger)
	usrSvc = proxymiddleware.UserLogicProxy(ctx, logger, server)(usrSvc)

	endpoints := usrep.NewUserLogicEndpoints(usrSvc)
	route := usrtr.MakeUserSvcHttpHandler(ctx, endpoints, logger)
	http.Handle("/", usrtr.AccessControl(route))

	errChan := make(chan error)

	svcId := fmt.Sprintf("%s-%d", "user", rand.Int())
	svcMeta := &discovery.ServiceMeta{
		Ip:          *host,
		Port:        uint16(*port),
		SvcName:     ServiceName,
		SvcId:       svcId,
		Weight:      *weight,
		Group:       *group,
		ClusterName: *clusterName,
		Check:       true,
		Healthy:     "",
	}

	//nacosCli, err := kitnacos.NewDefaultClient("public", "/Users/dk/github/nacos/nacos/mydata", cfgHost, cfgPort, logger)
	nacosCli, err := kitnacos.NewDefaultClient(2, "/Users/dk/github/nacos/nacos/mydata", logger, server)
	if err != nil {
		logger.Log("file", "main.go", "function", "main", "register_cli", err, "status", "panic exit")
		panic(err)
	}

	logger.Log("file", "main.go", "function", "main", "service", ServiceName, "status", "registering")
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

	logger.Log("file", "main.go", "function", "main", "service", ServiceName, "status", "service routine starting")
	go func() {
		listenAddr := fmt.Sprintf("%s:%d", *host, *port)
		logger.Log("file", "main.go", "function", "main", "service", "UserLogicService", "status", "start", "listen", listenAddr)
		errChan <- http.ListenAndServe(listenAddr, nil)
	}()

	logger.Log("file", "main.go", "function", "main", "service", ServiceName, "status", "signal routine starting")
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()

	er := <-errChan
	logger.Log("file", "main.go", "function", "main", "service", ServiceName, "status", "exit", "error", er)
}