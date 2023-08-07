package web

import (
	"encoding/json"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"message"`
}

const (
	ERROR     = 500
	SUCCESS   = 200
	JWTERROR  = 401
	AUTHERROR = 403
)

// FmtLongInt 超过16位的数字转为字符串
func FmtLongInt(data interface{}) interface{} {
	j, _ := json.Marshal(data)
	reg := regexp.MustCompile(`:(\d{16,20})`)
	l := len(reg.FindAllString(string(j), -1)) //正则匹配16-20位的数字，如果找到了就开始正则替换并解析

	if l != 0 {
		var mapResult interface{}
		str := reg.ReplaceAllString(string(j), `:"${1}"`)
		_ = json.Unmarshal([]byte(str), &mapResult)
		data = &mapResult
	}

	return data
}

func Result(code int, data interface{}, msg string, c *gin.Context) {
	data = FmtLongInt(data)
	// 开始时间
	c.JSON(http.StatusOK, Response{
		code,
		data,
		msg,
	})
}

func Ok(c *gin.Context) {
	Result(SUCCESS, map[string]interface{}{}, "success", c)
}

func OkWithMessage(message string, c *gin.Context) {
	Result(SUCCESS, map[string]interface{}{}, message, c)
}

func OkWithData(data interface{}, c *gin.Context) {
	Result(SUCCESS, data, "success", c)
}

func OkWithString(data string, c *gin.Context) {
	c.JSON(http.StatusOK, Response{
		SUCCESS,
		data,
		"success",
	})
}

func OkWithDetailed(data interface{}, message string, c *gin.Context) {
	Result(SUCCESS, data, message, c)
}

func Fail(c *gin.Context) {
	Result(ERROR, map[string]interface{}{}, "error", c)
}

func FailWithMessage(message string, c *gin.Context) {
	Result(ERROR, map[string]interface{}{}, message, c)
}

func FailWithMessageV(err interface{}, c *gin.Context) {
	Result(ERROR, err, "参数校验错误", c)
}

func FailWithDetailed(data interface{}, message string, c *gin.Context) {
	Result(ERROR, data, message, c)
}

func JWTFailWithDetailed(data interface{}, message string, c *gin.Context) {
	c.JSON(http.StatusUnauthorized, Response{
		JWTERROR,
		data,
		message,
	})
}

func OkWithDataLogin(data interface{}, message string, c *gin.Context) {
	c.JSON(http.StatusOK, Response{
		JWTERROR,
		data,
		message,
	})
}

func OkStatusData(code int, data interface{}, message string, c *gin.Context) {
	c.JSON(http.StatusOK, Response{
		code,
		data,
		message,
	})
}
