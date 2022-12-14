package rest

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hihibug/microdule/core/middleware"
	"io"
	"net/http"
	"os"
	"time"
)

type Gin struct {
	Route  *gin.Engine
	Config *Config
}

func NewGin(conf *Config) Rest {
	gin.SetMode(conf.Mode)

	if !conf.LogColType {
		// 禁用控制台颜色，将日志写入文件时不需要控制台颜色。
		gin.DisableConsoleColor()
		accessLogPath := conf.LogPath + "/access-" + time.Now().Format("2006-01-02") + ".log"
		// 记录到文件。
		f, _ := os.Create(accessLogPath)
		gin.DefaultWriter = io.MultiWriter(f)
	}

	var route = gin.Default()

	// 初始化页面
	if conf.UseHtml {
		route.Delims("{[{", "}]}")
		route.Static("/static/page", "/resource/page")
		route.LoadHTMLGlob("/resource/templates/*")
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

func (g *Gin) GetGin() *Gin {
	return g
}

func (g *Gin) Run() {
	fmt.Printf("Init Rest  Success Port: %s \n", g.Config.Addr)
	s := &http.Server{
		Addr:           ":" + g.Config.Addr,
		Handler:        g.Route,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	err := s.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
