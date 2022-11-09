package gorm

import (
	"fmt"
	"testing"
)

func TestNewGorm(t *testing.T) {
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
