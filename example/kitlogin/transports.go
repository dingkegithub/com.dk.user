package main

import (
	"context"
	"encoding/json"
	"fmt"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"strings"
)

func NewHtppHandler(ctx context.Context, eps *Endpoints) http.Handler {
	route := mux.NewRouter()

	option := []kithttp.ServerOption {
		kithttp.ServerErrorEncoder(encodeError),
	}

	route.Path("/api/kitlogin").Handler(kithttp.NewServer(
		eps.LoginEndpoint,
		decodeLoginRequest,
		encodeJsonResponse,
		option...,
	))

	route.Path("/api/logout").Handler(kithttp.NewServer(
		eps.LogoutEndpoint,
		decodeLogoutRequest,
		encodeJsonResponse,
		option...,
	))


	return route
}

func decodeLoginRequest(cxt context.Context, req *http.Request) (interface{}, error) {

	/*
	// form-data 机械
	req.ParseMultipartForm(1024)
	return &LoginRequest{
		Name: req.PostFormValue("name"),
		Pwd:  req.PostFormValue("pwd"),
	}, nil
	 */

	defer req.Body.Close()

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, fmt.Errorf("need kitlogin info")
	}

	loginRequest := &LoginRequest{}
	err = json.Unmarshal(body, &loginRequest)
	if err != nil {
		return nil, fmt.Errorf("body info error")
	}

	return loginRequest, nil
}

func decodeLogoutRequest(ctx context.Context, req *http.Request) (interface{}, error)  {
	token := req.Header.Get("token")
	//token-bql-123
	info := strings.Split(token, "-")
	if len(info) < 3 {
		return nil, fmt.Errorf("token need")
	}

	logoutRequest := &LogoutRequest{
		Name: info[1],
	}

	return logoutRequest, nil
}

func encodeJsonResponse(ctx context.Context, wr http.ResponseWriter, resp interface{}) error {
	wr.Header().Set("Content-Type", "application/json;charset=utf-8")
	return json.NewEncoder(wr).Encode(resp)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter)  {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.Header().Set("Access-Control-Allow-Methods", "GET,POST,OPTIONS")
	w.Header().Set("Access-Control-Max-Age", "3600")
	w.WriteHeader(http.StatusInternalServerError)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

func accessControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")

		if r.Method == "OPTIONS" {
			return
		}

		h.ServeHTTP(w, r)
	})
}