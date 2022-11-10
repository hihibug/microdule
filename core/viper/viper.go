package viper

import (
	"github.com/fsnotify/fsnotify"
	"github.com/hihibug/microdule/core/gorm"
	"github.com/spf13/viper"
)

type ConfigViper interface {
	ConfigToGormMysql(opt gorm.OptConfig) *gorm.Config
}

type Viper struct {
	Vp   *viper.Viper
	Data *Config
	Err  error
}

func NewViper(path string) Viper {
	conf := &Config{}
	v := viper.New()
	v.SetConfigFile(path)
	err := v.ReadInConfig()
	if err != nil {
		return Viper{nil, conf, err}
	}

	v.WatchConfig()

	v.OnConfigChange(func(e fsnotify.Event) {
		_ = v.Unmarshal(conf)
	})

	err = v.Unmarshal(conf)
	if err != nil {
		return Viper{nil, conf, err}
	}

	return Viper{v, conf, nil}
}

func (v *Viper) ConfigToGormMysql(opt gorm.OptConfig) *gorm.Config {
	return &gorm.Config{
		DbType:      v.Data.DB.DbType,
		MaxIdleCons: v.Data.DB.MaxIdleCons,
		MaxOpenCons: v.Data.DB.MaxOpenCons,
		LogMode:     v.Data.DB.LogMode,
		Opt:         opt,
		Mysql: &gorm.MysqlConfig{
			Path:     v.Data.DB.Path,
			Config:   v.Data.DB.Config,
			Dbname:   v.Data.DB.Dbname,
			Username: v.Data.DB.Username,
			Password: v.Data.DB.Password,
		},
	}
}
