package microdule

import (
	"errors"
)

type Service interface {
	Name() string
	Init(...Option)
	Options() *Options
	Close() func()
	Run() error
	Stop()
}

func NewService(opt ...Option) Service {
	return newService(opt...)
}

func (s *service) Name() string {
	return s.opts.Name
}

func (s *service) Init(opts ...Option) {
	for _, o := range opts {
		o(&s.opts)
	}
}

func (s *service) Options() *Options {
	return &s.opts
}

func (s *service) Close() func() {
	return func() {
		s.opts.Gorm.Close()
		s.opts.Redis.Close()
		s.opts.Etcd.Close()
	}
}

func (s *service) Run() error {
	return errors.New("test")
}

func (s *service) Stop() {

}
