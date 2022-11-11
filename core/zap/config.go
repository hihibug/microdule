package zap

import (
	"errors"
	"github.com/hihibug/microdule/core"
	"reflect"
)

var (
	ErrEmptyDirector = errors.New("empty log director")
)

type Config struct {
	Level         string `json:"level" yaml:"level"`
	Format        string `json:"format" yaml:"format"`
	Prefix        string `json:"prefix" yaml:"prefix"`
	Director      string `json:"director"  yaml:"director"`
	LinkName      string `json:"linkName" yaml:"linkName"`
	ShowLine      bool   `json:"showLine" yaml:"showLine"`
	EncodeLevel   string `json:"encodeLevel" yaml:"encodeLevel"`
	StacktraceKey string `json:"stacktraceKey" yaml:"stacktraceKey"`
	LogInConsole  bool   `json:"logInConsole" yaml:"logInConsole"`
}

func (c *Config) Validate() error {
	if reflect.DeepEqual(c, &Config{}) {
		return core.ErrEmptyConfig
	}

	if c.Director == "" {
		return ErrEmptyDirector
	}

	return nil
}
