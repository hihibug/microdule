package viper

import (
	"github.com/hihibug/microdule/core/etcd"
	"github.com/hihibug/microdule/core/redis"
	"github.com/hihibug/microdule/core/zap"
	"github.com/hihibug/microdule/rpc"
	"github.com/hihibug/microdule/web"
)

type Config struct {
	DB    DbConfig      `json:"db" yaml:"db"`
	Etcd  *etcd.Config  `json:"etcd" yaml:"etcd"`
	Redis *redis.Config `json:"redis" yaml:"redis"`
	Log   *zap.Config   `json:"log" yaml:"log"`
	Rest  *web.Config   `json:"rest" yaml:"rest"`
	Rpc   *rpc.Config   `json:"rpc" yaml:"rpc"`
}

type DbConfig struct {
	DbType      string `json:"db-type" yaml:"dbType"`
	Path        string `json:"path" yaml:"path"`
	Config      string `json:"config" yaml:"config"`
	Dbname      string `json:"dbname" yaml:"dbName"`
	Username    string `json:"username" yaml:"username"`
	Password    string `json:"password" yaml:"password"`
	MaxIdleCons int    `json:"maxIdleCons" yaml:"maxIdleCons"`
	MaxOpenCons int    `json:"maxOpenCons" yaml:"maxOpenCons"`
	LogMode     string `json:"logMode" yaml:"logMode"`
}
