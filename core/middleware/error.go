package middleware

import (
	"log"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

func GinErrorHttp(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			//打印错误堆栈信息
			debug.PrintStack()
			//封装通用json返回
			log.Println("recover err ", r)
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "服务器异常",
				"code":    http.StatusInternalServerError,
			})
		}
	}()

	c.Next()
}
