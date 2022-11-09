package gorm

import (
	"fmt"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"testing"
	"time"
)

func TestNewGorm(t *testing.T) {

	opt := GetGormConfigStruct()

	_default := logger.New(NewWriter(log.New(os.Stdout, "\r\n", log.LstdFlags)), logger.Config{
		SlowThreshold: 200 * time.Millisecond,
		LogLevel:      logger.Warn,
		Colorful:      true,
	})

	opt.Logger = _default.LogMode(logger.Info)

	DB, err := NewGorm(&Config{
		DbType:  "mysql",
		LogMode: "info",
		Mysql: &MysqlConfig{
			Dbname:   "user_content",
			Path:     "127.0.0.1:3306",
			Config:   "charset=utf8mb4&parseTime=True&loc=Local",
			Username: "root",
			Password: "root",
		},
		Opt: SetGormConfig(opt),
	})

	if err != nil {
		fmt.Println(err)
	}

	a := make([]map[string]interface{}, 0)
	err = DB.Client().Table("users").Find(&a).Error
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(a)
}
