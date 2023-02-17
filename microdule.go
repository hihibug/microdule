package microdule

import (
	"github.com/hihibug/microdule/rpc"
	"github.com/hihibug/microdule/web"
	"google.golang.org/grpc"
)

type Service interface {
	Name() string
	Init(...Option)
	Options() *Options
	Http() web.Web
	Rpc(...grpc.ServerOption) rpc.Rpc
	Close()
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

func (s *service) Close() {

	if s.opts.Gorm != nil {
		s.opts.Gorm.Close()
	}

	if s.opts.Redis != nil {
		s.opts.Redis.Close()
	}

	if s.opts.Etcd != nil {
		s.opts.Etcd.Close()
	}

	return
}

func (s *service) Run() error {
	//fmt.Println(s)
	return nil
}

func (s *service) Http() web.Web {
	return web.NewGin(s.opts.Config.Data.Http)
}

func (s *service) Rpc(opt ...grpc.ServerOption) rpc.Rpc {
	return rpc.NewGrpc(s.opts.Config.Data.Rpc, opt...)
}

func (s *service) Stop() {

}
