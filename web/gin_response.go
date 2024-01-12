package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type (
	GinResponse struct {
		Context *gin.Context
	}
)

func NewGinResponse(c *gin.Context) *GinResponse {
	return &GinResponse{Context: c}
}

func GinResult(code int, data interface{}, msg string, c *gin.Context) {
	data = FmtLongInt(data)
	// 开始时间
	c.JSON(http.StatusOK, ResponseData{
		code,
		data,
		msg,
	})
}

func (r *GinResponse) Ok() {
	GinResult(SUCCESS, map[string]interface{}{}, "success", r.Context)
}

func (r *GinResponse) OkWithMessage(message string) {
	GinResult(SUCCESS, map[string]interface{}{}, message, r.Context)
}

func (r *GinResponse) OkWithData(data interface{}) {
	GinResult(SUCCESS, data, "success", r.Context)
}

func (r *GinResponse) OkWithString(data string) {
	r.Context.JSON(http.StatusOK, ResponseData{
		SUCCESS,
		data,
		"success",
	})
}

func (r *GinResponse) OkWithDetailed(data interface{}, message string) {
	GinResult(SUCCESS, data, message, r.Context)
}

func (r *GinResponse) Fail() {
	GinResult(ERROR, map[string]interface{}{}, "error", r.Context)
}

func (r *GinResponse) FailWithMessage(message string) {
	GinResult(ERROR, map[string]interface{}{}, message, r.Context)
}

func (r *GinResponse) FailWithDataMessages(err interface{}) {
	GinResult(ERROR, err, "参数校验错误", r.Context)
}

func (r *GinResponse) FailWithDetailed(data interface{}, message string) {
	GinResult(ERROR, data, message, r.Context)
}

func (r *GinResponse) JWTFailWithDetailed(data interface{}, message string) {
	r.Context.JSON(http.StatusUnauthorized, ResponseData{
		JWTERROR,
		data,
		message,
	})
}

func (r *GinResponse) OkWithDataLogin(data interface{}, message string) {
	r.Context.JSON(http.StatusOK, ResponseData{
		JWTERROR,
		data,
		message,
	})
}

func (r *GinResponse) OkStatusData(code int, data interface{}, message string) {
	r.Context.JSON(http.StatusOK, ResponseData{
		code,
		data,
		message,
	})
}
