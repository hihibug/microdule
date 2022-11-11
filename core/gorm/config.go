package gorm

import (
	"errors"
	"github.com/hihibug/microdule/core"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"reflect"
	"time"
)

var (
	ErrEmptyPath   = errors.New("config empty path")
	ErrEmptyDbName = errors.New("config empty db_name")
)

type (
	ConfigGorm interface {
		Validate() error
	}

	MysqlConfig struct {
		Path     string `json:"path" yaml:"path"`
		Config   string `json:"config" yaml:"config"`
		Dbname   string `json:"dbname" yaml:"db-name"`
		Username string `json:"username" yaml:"username"`
		Password string `json:"password" yaml:"password"`
	}

	Config struct {
		DbType      string `json:"dbType" yaml:"db-type"`
		MaxIdleCons int    `json:"maxIdleCons" yaml:"max-idle-cons"`
		MaxOpenCons int    `json:"maxOpenCons" yaml:"max-open-cons"`
		LogMode     string `json:"logMode" yaml:"log-mode"`
		Opt         func(mode string) *gorm.Config
		Mysql       *MysqlConfig
	}

	Writer struct {
		logger.Writer
	}
)

func DefaultOpt(mode string) *gorm.Config {
	config := &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	}
	//NewWriter 对log.New函数的再次封装，从而实现是否通过zap打印日志
	_default := logger.New(NewWriter(log.New(os.Stdout, "\r\n", log.LstdFlags)), logger.Config{
		SlowThreshold: 200 * time.Millisecond,
		LogLevel:      logger.Warn,
		Colorful:      true,
	})
	LogGorm(mode, config, _default)

	return config
}

func LogGorm(mode string, config *gorm.Config, _default logger.Interface) {
	//设置logger的日志输出等级
	switch mode {
	case "silent", "Silent":
		config.Logger = _default.LogMode(logger.Silent)
	case "error", "Error":
		config.Logger = _default.LogMode(logger.Error)
	case "warn", "Warn":
		config.Logger = _default.LogMode(logger.Warn)
	case "info", "Info":
		config.Logger = _default.LogMode(logger.Info)
	default:
		config.Logger = _default.LogMode(logger.Info)
	}
}

func NewWriter(w logger.Writer) *Writer {
	return &Writer{Writer: w}
}

func (w *Writer) Printf(format string, args ...interface{}) {
	w.Writer.Printf(format, args...)
	return
}

func (c *Config) Validate() error {
	if reflect.DeepEqual(c, &Config{}) {
		return core.ErrEmptyConfig
	}

	if c.MaxIdleCons <= 0 {
		c.MaxIdleCons = 100
	}

	if c.MaxOpenCons <= 0 {
		c.MaxOpenCons = 10
	}

	if c.Opt == nil {
		opt := DefaultOpt(c.LogMode)
		c.Opt = func(mode string) *gorm.Config {
			return opt
		}
	}
	return nil
}

func (c *MysqlConfig) Validate() error {
	if c.Path == "" {
		return ErrEmptyPath
	}

	if c.Dbname == "" {
		return ErrEmptyDbName
	}
	return nil
}
