package microdule

import (
	"sync"
)

type service struct {
	opts Options
	once sync.Once
	err  error
}

func newService(opts ...Option) *service {
	return &service{
		opts: newOptions(opts...),
	}
}
