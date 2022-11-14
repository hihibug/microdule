package microdule

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"sync"
	"testing"
	"time"
)

var global *Options

func TestNewService(t *testing.T) {
	//初始化服务
	s := NewService(Name("test"))

	global = s.Options()

	//gormConf := gorm.GetGormConfigStruct()

	//NewZapWriter 对log.New函数的再次封装，从而实现是否通过zap打印日志
	//gorm.LogGorm(
	//	global.Config.Data.DB.LogMode,
	//	gormConf,
	//	gorm.SetGormLogZap(zap.NewZapWriter(global.Log.Client())),
	//)

	//获取db配置
	//dbConf := global.Config.ConfigToGormMysql(gorm.SetGormConfig(gormConf))

	global.Config.Data.Etcd.Log = global.Log.Client()

	//初始化组件
	s.Init(
		//Gorm(dbConf),
		Etcd(global.Config.Data.Etcd),
		//Redis(global.Config.Data.Redis),
	)

	//关闭链接
	defer s.Close()

	rest := s.Rest("gin")
	a := rest.GetGin().Route.Group("")
	{
		a.GET("/test", func(context *gin.Context) {
			fmt.Println("test")
		})
		a.GET("/err", func(c *gin.Context) {
			panic("test")
		})
	}

	rest.Run()
}

func GoMysql(num, cnum int) {
	var wg sync.WaitGroup
	ch := make(chan struct{}, cnum)
	for i := 0; i < num; i++ {
		ch <- struct{}{}
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			a := make([]map[string]interface{}, 0)
			err := global.Gorm.Client().Table("users").Find(&a).Error
			if err != nil {
				log.Println(err)
			}
			log.Println(a)
			time.Sleep(time.Second)
			<-ch
		}(i)
	}
	wg.Wait()
}

func GoRedis(num, cnum int) {
	var wg sync.WaitGroup
	ch := make(chan struct{}, cnum)
	for i := 0; i < num; i++ {
		ch <- struct{}{}
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			c := global.Redis.Client().Get("test-11")
			if c.Err() != nil {
				log.Println(c.Err())
			}
			log.Println(c.Val())
			time.Sleep(time.Second)
			<-ch
		}(i)
	}
	wg.Wait()
}
