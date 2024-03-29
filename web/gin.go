package web

import (
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hihibug/microdule/core/middleware"
	"github.com/hihibug/microdule/core/utils"
)

func NewGin(conf *Config) Web {
	gin.SetMode(conf.Mode)

	if !conf.LogColType {
		// 禁用控制台颜色，将日志写入文件时不需要控制台颜色。
		gin.DisableConsoleColor()
		defPath, _ := os.Getwd()
		path := defPath + "/" + conf.LogPath
		if ok, _ := utils.PathExists(path); !ok { // 判断是否有Director文件夹
			_ = os.Mkdir(path, os.ModePerm)
		}
		accessLogPath := path + "/access-" + time.Now().Format("2006-01-02") + ".log"
		// 记录到文件。
		f, _ := os.Create(accessLogPath)
		gin.DefaultWriter = io.MultiWriter(f)
	}

	var route = gin.Default()

	// 初始化页面
	if conf.UseHtml {
		defPath, _ := os.Getwd()
		route.Delims(conf.DelimsStr, conf.DelimsEnd)
		route.Static(defPath+"/"+conf.StaticPath, defPath+"/"+conf.TmpPath)
		route.LoadHTMLGlob(defPath + "/" + conf.TmpPath + "/*")
	}

	//注册GinCors
	route.Use(middleware.GinCors(), middleware.GinErrorHttp)
	route.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"message": "page not found", "code": http.StatusNotFound})
	})

	return &Gin{
		route,
		conf,
	}
}

func (g *Gin) Client() any {
	return g
}

func (g *Gin) Run() error {
	log.Printf("rest  port: %s \n", g.Config.Addr)
	s := &http.Server{
		Addr:           ":" + g.Config.Addr,
		Handler:        g.Route,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	err := s.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}
