package web_test

import (
	"github.com/hihibug/microdule/web"
	"testing"
)

func TestRest(t *testing.T) {
	_ = web.NewGin(&web.Config{
		Mode:       "debug",
		LogColType: false,
		LogPath:    "",
		UseHtml:    false,
		Addr:       "8999",
	})

	//rs := r.GetGin()
	//
	//a := rs.Route.Group("")
	//{
	//	a.GET("/test", func(context *gin.Context) {
	//		fmt.Println("test")
	//	})
	//}
	//
	//r.Run()
}
