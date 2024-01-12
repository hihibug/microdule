package web

import "github.com/gin-gonic/gin"

type Web interface {
	Client() any
	Run() error
}

type Gin struct {
	Route  *gin.Engine
	Config *Config
}
