package apollo

type Option func(options *apolloOptions)

type apolloOptions struct {
	Cluster string
	CacheFile string
}

func newApolloOptions() *apolloOptions {
	return &apolloOptions{
		Cluster:   "default",
		CacheFile: "/tmp/.cfg",
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
		options.Cluster = c
	}
}