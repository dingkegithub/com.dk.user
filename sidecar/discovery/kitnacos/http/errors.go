package http

import "fmt"

var (
	ErrNotFoundHealthyNode = fmt.Errorf("could not found healthy node, check your cluster")
)
