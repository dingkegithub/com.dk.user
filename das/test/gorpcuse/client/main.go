package main

import (
	"com.gz.heartbeat/test/gorpcuse/svc"
	"fmt"
	"log"
	"net/rpc"
)

func main() {
	client, err := rpc.DialHTTP("tcp", "127.0.0.1:1234")
	if err != nil {
		log.Fatal("dialing:", err)
	}

	req := &svc.StringRequest{"A", "B"}

	var reply string
	err = client.Call("StringService.Concat", req, &reply)
	if err != nil {
		log.Fatal("Concat error:", err)
	}
	fmt.Printf("StringService Concat: %s concat %s = %s\n", req.A, req.B, reply)

	req = &svc.StringRequest{"ACD", "BDF"}
	call := client.Go("StringService.Diff", req, &reply, nil)
	_ = <-call.Done
	fmt.Printf("StringService.Diff: %s diff %s = %s\n", req.A, req.B, reply)
}
