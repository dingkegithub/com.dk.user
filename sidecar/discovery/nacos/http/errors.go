package http

import "fmt"

var (
	ErrCacheFileExist      = fmt.Errorf("not found cache file")
	ErrNotFoundHealthyNode = fmt.Errorf("could not found healthy node, check your cluster")
)
