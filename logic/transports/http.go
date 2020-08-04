package transports

import (
	"context"
	usrep "github.com/dingkegithub/com.dk.user/logic/endpoints"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"net/http"
)


/*
 * 路由设置，http请求绑定到站点
 *
 * @param ctx 全局上下文
 * @param endpoints 站点集合
 * @param logger 统一logger接口
 */
func MakeUserSvcHttpHandler(ctx context.Context, endpoints *usrep.UserEndpoints, logger log.Logger) http.Handler { // 初始化路由
	route := mux.NewRouter()

	// 服务器内部错误时的处理handler和响应
	options := []kithttp.ServerOption{
		kithttp.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		kithttp.ServerErrorEncoder(EncodeError),
	}

	// 路由绑定站点: 用户注册
	route.Methods("POST").Path("/api/register").Handler(kithttp.NewServer(
		endpoints.RegisterEndpoint,
		decodeRegisterRequest,
		encodeRegisterResponse,
		options...,
	))

	return route
}
