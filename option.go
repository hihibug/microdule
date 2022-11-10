package microdule

import (
	"context"
	"github.com/hihibug/microdule/core/etcd"
	"github.com/hihibug/microdule/core/gorm"
	"github.com/hihibug/microdule/core/viper"
)

type (
	Option func(*Options)

	Options struct {
		DB      gorm.Gorm
		Etcd    etcd.Etcd
		Name    string
		Config  viper.Viper
		Context context.Context
	}
)

func newOptions(opts ...Option) Options {
	opt := Options{
		Config:  viper.NewViper("config.yml"),
		Context: context.Background(),
	}

	if opt.Config.Err != nil {
		panic(opt.Config.Err)
	}

	for _, o := range opts {
		o(&opt)
	}

	return opt
}

func Name(n string) Option {
	return func(options *Options) {
		options.Name = n
	}
}

func DB(db gorm.Gorm) Option {
	return func(options *Options) {
		options.DB = db
	}
}

func ETCD(e etcd.Etcd) Option {
	return func(options *Options) {
		options.Etcd = e
	}
}
