package discovery

import "fmt"

var (
	ErrorParam = fmt.Errorf("param error")
	ErrorNotFoundAnyService = fmt.Errorf("not found any available service instance")
)
