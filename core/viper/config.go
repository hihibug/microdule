package viper

import "github.com/hihibug/microdule/core/etcd"

type Config struct {
	DB   DbConfig     `json:"db" yaml:"db"`
	Etcd *etcd.Config `json:"etcd" yaml:"etcd"`
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

type EtcdConfig struct {
	Addr     string `json:"addr" yaml:"addr"`
	Password string `json:"password" yaml:"password"`
	TimeOut  int    `json:"time-out" yaml:"time-out"`
}
