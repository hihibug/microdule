package microdule

import "github.com/hihibug/microdule/rest"

type Service interface {
	Name() string
	Init(...Option)
	Options() *Options
	Rest(string) rest.Rest
	Rpc() rest.Rest
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

func (s *service) Rest(t string) rest.Rest {
	switch t {
	case "gin":
		return rest.NewGin(s.opts.Config.Data.Rest)
	default:
		return rest.NewGin(s.opts.Config.Data.Rest)
	}
}

func (s *service) Rpc() rest.Rest {
	//TODO implement me
	panic("implement me")
}

func (s *service) Stop() {

}
