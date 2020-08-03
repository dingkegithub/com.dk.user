package main

import (
	"context"
	"fmt"
	"net/http"
)

func main() {
	ep := NewEps(NewDefaultService())
	tr := NewHtppHandler(context.Background(), ep)

	http.Handle("/", accessControl(tr))

	errCh := make(chan error)
	go func() {
		fmt.Println("server listen on 8080")
		errCh <- http.ListenAndServe("localhost:8080", nil)
	}()

	e := <-errCh
	fmt.Println("server exit: ", e)
}
