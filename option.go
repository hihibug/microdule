package microdule

import (
	"context"
	"github.com/hihibug/microdule/core/etcd"
	"github.com/hihibug/microdule/core/gorm"
	"github.com/hihibug/microdule/core/redis"
	"github.com/hihibug/microdule/core/viper"
	"github.com/hihibug/microdule/core/zap"
)

type (
	Option func(*Options)

	Options struct {
		Name string

		Gorm  gorm.Gorm
		Etcd  etcd.Etcd
		Redis redis.Redis

		Log    zap.Log
		Config viper.Viper

		Context context.Context
	}
)

func newOptions(opts ...Option) Options {
	opt := Options{
		Config:  viper.NewViper("config.yml"),
		Context: context.Background(),
	}

	opt.Log = zap.NewZap(opt.Config.Data.Log)

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

func Gorm(dbConf *gorm.Config) Option {
	db, err := gorm.NewGorm(dbConf)
	if err != nil {
		panic("mysql error " + err.Error())
	}
	return func(options *Options) {
		options.Gorm = db
	}
}

func Etcd(e *etcd.Config) Option {
	etd, err := etcd.NewEtcd(e)
	if err != nil {
		panic("etcd error " + err.Error())
	}
	return func(options *Options) {
		options.Etcd = etd
	}
}

func Redis(r *redis.Config) Option {
	rds, err := redis.NewRedis(r)
	if err != nil {
		panic("redis error " + err.Error())
	}
	return func(options *Options) {
		options.Redis = rds
	}
}
