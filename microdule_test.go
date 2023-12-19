package microdule

import (
	"log"
	"sync"
	"testing"
	"time"
)

var global *Options

func TestNewService(t *testing.T) {
	//初始化服务 config初始化键值为0
	s := NewService(
		Config("core/config.yml"),
		Name("test"),
	)

	global = s.Options()

	// fmt.Println(global.Config.Vp.Get("QiNu.Key"))
	////NewZapWriter 对log.New函数的再次封装，从而实现是否通过zap打印日志
	//gormConf := gorm.GetGormConfigStruct()
	//gorm.LogGorm(
	//	global.Config.Data.DB.LogMode,
	//	gormConf,
	//	gorm.SetGormLogZap(zap.NewZapWriter(global.Log.Client())),
	//)
	//
	//// 设置etcd 日志
	//global.Config.Data.Etcd.Log = global.Log.Client()
	//
	////初始化组件
	//s.Init(
	////Redis(nil),
	////Gorm(global.Config.ConfigToGormMysql(gorm.SetGormConfig(gormConf))),
	////Etcd(global.Config.Data.Etcd),
	//)
	//
	////关闭链接
	//defer s.Close()
	//
	////开启rest
	//rest := s.Http().Client()
	//
	//a := rest.Route.Group("")
	//{
	//	a.GET("/test", func(context *gin.Context) {
	//		fmt.Println("test")
	//		global.Log.Client().Info("test")
	//	})
	//	a.GET("/err", func(c *gin.Context) {
	//		panic("test")
	//	})
	//}

	//ip, _ := utils.ExternalIP()
	//global.Config.Data.Rpc.IP = ip
	//rpc := s.Rpc().Client()
	//register, err := s.Rpc().Client().Register(global.Etcd.Clients())
	//if err != nil {
	//	os.Exit(0)
	//}
	//go register.ListenLeaseRespChan()
	//
	//rpc.Run()
	//rest.Run()
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
