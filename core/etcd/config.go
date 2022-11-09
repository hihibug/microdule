package etcd

import (
	"errors"
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
}

func (c *Config) Validate() error {
	if c.Addr == "" {
		return ErrEmptyAddr
	}

	if c.TimeOut == 0 {
		c.TimeOut = 5
	}

	return nil
}
