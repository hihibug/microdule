package microdule

import (
	"context"
	"github.com/hihibug/microdule/core/etcd"
	"github.com/hihibug/microdule/core/gorm"
	"github.com/hihibug/microdule/core/redis"
	"github.com/hihibug/microdule/core/viper"
	"github.com/hihibug/microdule/core/zap"
	"github.com/hihibug/microdule/rest"
	"github.com/hihibug/microdule/rpc"
	"reflect"
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

		Rest rest.Rest
		Rpc  rpc.Rpc

		Context context.Context
	}
)

func newOptions(opts ...Option) Options {
	opt := Options{
		Context: context.Background(),
	}

	for k, o := range opts {
		if k > 0 && reflect.DeepEqual(opt.Config, viper.Viper{}) {
			opt.Config = viper.NewViper("config.yml")
		}

		if opt.Config.Err != nil {
			panic(opt.Config.Err)
		}
		o(&opt)
	}

	if opt.Log == nil {
		opt.Log = zap.NewZap(opt.Config.Data.Log)
	}

	return opt
}

func Name(n string) Option {
	return func(options *Options) {
		options.Name = n
	}
}

func Config(path string) Option {
	return func(options *Options) {
		options.Config = viper.NewViper(path)
	}
}

func Gorm(dbConf *gorm.Config) Option {
	return func(options *Options) {
		if dbConf == nil {
			dbConf = options.Config.ConfigToGormMysql(gorm.SetGormConfig(gorm.GetGormConfigStruct()))
		}
		db, err := gorm.NewGorm(dbConf)
		if err != nil {
			panic("mysql error " + err.Error())
		}
		options.Gorm = db
	}
}

func Etcd(e *etcd.Config) Option {
	return func(options *Options) {
		if e == nil {
			e = options.Config.Data.Etcd
		}
		etd, err := etcd.NewEtcd(e)
		if err != nil {
			panic("etcd error " + err.Error())
		}
		options.Etcd = etd
	}
}

func Redis(r *redis.Config) Option {
	return func(options *Options) {
		if r == nil {
			r = options.Config.Data.Redis
		}
		rds, err := redis.NewRedis(r)
		if err != nil {
			panic("redis error " + err.Error())
		}
		options.Redis = rds
	}
}
