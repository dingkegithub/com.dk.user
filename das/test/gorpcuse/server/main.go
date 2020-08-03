package main

import (
	"com.dk.user/das/test/gorpcuse/svc"
	"log"
	"net"
	"net/http"
	"net/rpc"
)

func main() {
	stringSvc := new(svc.StringService)
	rpc.Register(stringSvc)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", "127.0.0.1:1234")
	if e != nil {
		log.Fatal("listen error: ", e)
	}
	http.Serve(l, nil)
}
