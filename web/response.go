package web

import (
	"encoding/json"
	"regexp"
)

const (
	ERROR     = 500
	SUCCESS   = 200
	JWTERROR  = 401
	AUTHERROR = 403
)

type (
	Response interface {
		Ok()
		OKWithMessage(string)
		OKWithData(interface{})
		OkWithString(string)
		OkWithDetailed(interface{}, string)
		Fail()
		FailWithMessage(string)
		FailWithDataMessages(interface{})
		FailWithDetailed(interface{}, string)
		JWTFailWithDetailed(interface{}, string)
		OkWithDataLogin(interface{}, string)
		OkStatusData(int, interface{}, string)
	}

	ResponseData struct {
		Code int         `json:"code"`
		Data interface{} `json:"data"`
		Msg  string      `json:"message"`
	}
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
