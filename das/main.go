package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/dingkegithub/com.dk.user/das/endpoints"
	"github.com/dingkegithub/com.dk.user/das/model"
	"github.com/dingkegithub/com.dk.user/das/proto/userpb"
	"github.com/dingkegithub/com.dk.user/das/service"
	"github.com/dingkegithub/com.dk.user/das/service/impl"
	"github.com/dingkegithub/com.dk.user/das/transport"
	"github.com/dingkegithub/com.dk.user/sidecar/discovery"
	"github.com/dingkegithub/com.dk.user/sidecar/discovery/kitnacos"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"
	"google.golang.org/grpc"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"syscall"
)

var (
	ServiceName = "UserDasService"
)


func main() {
	ctx := context.Background()
	errChan := make(chan error)

	flag.Parse()
	host := flag.String("host", "127.0.0.1", "-h x.x.x.x --host=x.x.x.x")
	port := flag.Uint("port", 8080, "-p x or --port=x")

	var logging log.Logger
	{
		logging = log.NewLogfmtLogger(os.Stderr)
	}

	logging.Log("service", ServiceName, "status", "init model")
	model.Init("root", "123456", "user")

	logging.Log("service", ServiceName, "status", "init service")
	var svc service.UserSvc
	{
		svc = impl.NewUserSvc()
		//svc = middleware.LoggingMiddleware(logging)(svc)
	}

	logging.Log("service", ServiceName, "status", "register service")
	ends := endpoints.NewUsrEndpoints(svc)
	handler := transport.NewRpcUsrSvc(ctx, ends)

	var discoverCli sd.Registrar
	{
		svcId := fmt.Sprintf("%s-%d", "UserDas", rand.Int())
		svcMeta := &discovery.ServiceMeta{
			Ip:          *host,
			Port:        uint16(*port),
			SvcName:     ServiceName,
			SvcId:       svcId,
			Weight:      0,
			Group:       "default",
			ClusterName: "default",
			Check:       true,
			Healthy:     "",
		}

		//nacosCli, err := kitnacos.NewDefaultClient("public", "/Users/dk/github/nacos/nacos/mydata", "127.0.0.1", 8848, logging)
		nacosCli, err := kitnacos.NewDefaultClient(2, "/Users/dk/github/nacos/nacos/mydata", logging, "127.0.0.1:8848")
		if err != nil {
			logging.Log("register_cli", err, "status", "panic exit")
			panic(err)
		}

		discoverCli, err = kitnacos.NewRegistrar(nacosCli, svcMeta, logging)
		if err != nil {
			logging.Log("register", err, "status", "panic exit")
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

		logging.Log("service", "UserDasService", "listen", lisAddr)
		err := rpcServer.Serve(ls)
		if err != nil {
			logging.Log("service", "UserDasService", "status", "exit", "error", err)
			errChan <- fmt.Errorf("%s", err.Error())
		}
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()

	<-errChan
	fmt.Println("program exit")
	logging.Log("service", "UserDasService", "listen", lisAddr)
}
