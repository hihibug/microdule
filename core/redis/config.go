package redis

import (
	"errors"
	"github.com/hihibug/microdule/core"
	"reflect"
)

var (
	// ErrEmptyAddr etcd 地址为空输出错误
	ErrEmptyAddr = errors.New("empty redis addr")
)

type Config struct {
	DB       int    `json:"db" yaml:"db"`
	Addr     string `json:"addr" yaml:"addr"`
	Password string `json:"password" yaml:"password"`
}

func (c *Config) Validate() error {
	if reflect.DeepEqual(c, &Config{}) {
		return core.ErrEmptyConfig
	}

	if c.Addr == "" {
		return ErrEmptyAddr
	}

	return nil
}
