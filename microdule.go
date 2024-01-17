package microdule

type Service interface {
	Name() string
	Init(...Option)
	Options() *Options
	Close()
	Start() error
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

	s.opts.Teamwork.Close()

	if s.opts.Rpc != nil {
		s.opts.Rpc.Close()
	}

	if s.opts.Gorm != nil {
		s.opts.Gorm.Close()
	}

	if s.opts.Redis != nil {
		s.opts.Redis.Close()
	}

	if s.opts.Etcd != nil {
		s.opts.Etcd.Close()
	}
}

func (s *service) Start() error {
	if s.opts.Http != nil {
		s.opts.Teamwork.Reginster("http", func() error {
			return s.opts.Http.Run()
		})
	}

	if s.opts.Rpc != nil {
		s.opts.Teamwork.Reginster("rpc ", func() error {
			return s.opts.Rpc.Run()
		})
	}

	return s.opts.Teamwork.Start()
}
