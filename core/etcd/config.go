package etcd

import (
	"encoding/json"
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

// DeConfig 解析配置json
func DeConfig(conf string) (c Config, err error) {

	err = json.Unmarshal([]byte(conf), &c)

	if err != nil {
		return Config{}, err
	}

	if c.Addr == "" {
		return Config{}, ErrEmptyAddr
	}

	return
}
