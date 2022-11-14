package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"runtime/debug"
)

func GinErrorHttp(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			//打印错误堆栈信息
			debug.PrintStack()
			//封装通用json返回
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": r.(string),
				"code":    http.StatusInternalServerError,
			})
		}
	}()

	c.Next()
}
