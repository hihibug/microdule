package rest_test

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hihibug/microdule/rest"
	"testing"
)

func TestRest(t *testing.T) {
	r := rest.NewGin(&rest.Config{
		Mode:       "debug",
		LogColType: false,
		LogPath:    "",
		UseHtml:    false,
		Addr:       "8999",
	})

	rs := r.GetGin()

	a := rs.Route.Group("")
	{
		a.GET("/test", func(context *gin.Context) {
			fmt.Println("test")
		})
	}

	r.Run()
}
