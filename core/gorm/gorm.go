package gorm

import (
	"fmt"
	"github.com/hihibug/microdule/core/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

type (
	DB struct {
		Cli *gorm.DB
	}

	Gorm interface {
		Client() *gorm.DB
		Close() error
	}

	OptConfig func(mode string) *gorm.Config
)

func (db *DB) Client() *gorm.DB {
	return db.Cli
}

func (db *DB) Close() error {
	d, _ := db.Cli.DB()
	return d.Close()
}

// NewGorm 新建gorm
func NewGorm(conf *Config) (Gorm, error) {
	err := conf.Validate()
	if err != nil {
		return nil, err
	}

	var dcr gorm.Dialector

	switch conf.DbType {
	case "mysql":
		dcr = NewMysql(conf.Mysql)
		err := conf.Mysql.Validate()
		if err != nil {
			return nil, err
		}
	default:
		dcr = NewMysql(conf.Mysql)
		err := conf.Mysql.Validate()
		if err != nil {
			return nil, err
		}
	}

	db, err := gorm.Open(dcr, conf.Opt(conf.LogMode))

	if err != nil {
		return nil, err
	} else {
		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(conf.MaxIdleCons)
		sqlDB.SetMaxOpenConns(conf.MaxOpenCons)
		fmt.Printf("Init Gorm  Success \n")
		return &DB{
			db,
		}, nil
	}
}

// NewMysql 初始化Mysql数据库
func NewMysql(m *MysqlConfig) gorm.Dialector {
	dsn := m.Username + ":" + m.Password + "@tcp(" + m.Path + ")/" + m.Dbname + "?" + m.Config
	return mysql.New(mysql.Config{
		DSN:                       dsn,   // DSN data source name
		DefaultStringSize:         191,   // string 类型字段的默认长度
		DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据版本自动配置
	})
}

func GetGormConfigStruct() *gorm.Config {
	return &gorm.Config{}
}

func SetGormLogZap(z *zap.Zap) logger.Interface {
	return logger.New(z, logger.Config{
		SlowThreshold: 200 * time.Millisecond,
		LogLevel:      logger.Warn,
		Colorful:      false,
	})
}

func SetGormConfig(c *gorm.Config) OptConfig {
	return func(mode string) *gorm.Config {
		return c
	}
}
