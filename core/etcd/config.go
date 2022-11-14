package etcd

import (
	"errors"
	"github.com/hihibug/microdule/core"
	"go.uber.org/zap"
	"reflect"
)

var (
	// ErrEmptyAddr etcd 地址为空输出错误
	ErrEmptyAddr = errors.New("empty etcd addr")
)

// Config Etcd 配置文件
type Config struct {
	Addr     string `json:"addr" yaml:"addr"`
	Password string `json:"password" yaml:"password"`
	TimeOut  int    `json:"time-out" yaml:"time-out"`
	Log      *zap.Logger
}

func (c *Config) Validate() error {

	if reflect.DeepEqual(c, &Config{}) {
		return core.ErrEmptyConfig
	}

	if c.Addr == "" {
		return ErrEmptyAddr
	}

	if c.TimeOut == 0 {
		c.TimeOut = 5
	}

	return nil
}
