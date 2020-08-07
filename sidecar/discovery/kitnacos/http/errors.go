package http

import "fmt"

var (
	ErrNotFoundHealthyNode = fmt.Errorf("could not found healthy node, check your cluster")
	ErrCache = fmt.Errorf("could not cache")
    ErrCacheFileExist = fmt.Errorf("not found cache file")
)
