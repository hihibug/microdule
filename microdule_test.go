package microdule

import (
	"fmt"
	"github.com/hihibug/microdule/core/etcd"
	"github.com/hihibug/microdule/core/gorm"
	"testing"
)

var Ser Options

func TestNewService(t *testing.T) {

	//初始化服务
	s := NewService(Name("test"))

	Ser = s.Options()

	//获取db配置
	dbConf := Ser.Config.ConfigToGormMysql(nil)
	//初始化db
	db, err := gorm.NewGorm(dbConf)
	if err != nil {
		panic("mysql error " + err.Error())
	}
	defer db.Close()

	//初始化etcd
	etd, err := etcd.NewEtcd(Ser.Config.Data.Etcd)
	if err != nil {
		panic("etcd error " + err.Error())
	}
	defer etd.Close()

	s.Init(
		DB(db),
		ETCD(etd),
	)

	err = s.Run()
	if err != nil {
		fmt.Println(err)
	}
}
