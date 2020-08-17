package apollo

import (
	"fmt"
	uuid2 "github.com/pborman/uuid"
	"time"
)

type Option func(options *apolloOptions)

type apolloOptions struct {
	Cluster string
	CacheFile string
}

func newApolloOptions() *apolloOptions {
	s := time.Now().Unix()
	uuid := uuid2.New()
	return &apolloOptions{
		Cluster:   "default",
		CacheFile: fmt.Sprintf("/tmp/.%s-%s", s, uuid),
	}
}

func (a *apolloOptions) apply(options ...Option) {
	apolloOptions := &apolloOptions{}

	for _, option := range options {
		option(apolloOptions)
	}
}

func WithCache(f string) Option {
	return func(options *apolloOptions) {
		options.CacheFile = f
	}
}

func WithCluster(c string) Option {
	return func(options *apolloOptions) {
		if c == "" {
			return
		}
		options.Cluster = c
	}
}