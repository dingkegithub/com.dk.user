package transports

import (
	"context"
	"encoding/json"
	"github.com/dingkegithub/com.dk.user/logic/service"
	"io/ioutil"
	"net/http"
)

/**
 * http 错误响应
 *
 * @param ctx 请求上下文
 * @param err 内部错误
 * @param w 携带http响应
 */
func EncodeError(_ context.Context, err error, w http.ResponseWriter)  {
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.WriteHeader(http.StatusInternalServerError)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"err": 50000,
		"msg": err.Error(),
	})
}

/**
 * http 错误响应
 *
 * @param ctx 请求上下文
 * @param req http请求
 * @return http请求转换为站点输入参数
 */
func decodeRegisterRequest(ctx context.Context, req *http.Request) (interface{}, error) {
	defer func() {
		_ = req.Body.Close()
	}()

	data, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}

	registerReq := &service.RegisterUsrRequest{}
	err = json.Unmarshal(data, &registerReq)
	if err != nil {
		return nil, err
	}

	return registerReq, nil
}

/**
 * 站点响应写入http响应
 *
 * @param ctx 请求上下文
 * @param rw http响应
 * @param resp 站点响应
 */
func encodeRegisterResponse(ctx context.Context, rw http.ResponseWriter, resp interface{}) error {
	rw.Header().Set("Content-Type", "application/json;charset=utf-8")
	return json.NewEncoder(rw).Encode(resp)
}
