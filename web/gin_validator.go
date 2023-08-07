package web

import (
	"fmt"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
)

type GinValidator struct {
	Validator *validator.Validate
	Trans     ut.Translator
}

// NewGinValidator 创建gin验证器
func NewGinValidator(locale string) (validators *GinValidator, err error) {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		validators = &GinValidator{}
		validators.Validator = v

		zhT := zh.New() // 中文翻译器
		enT := en.New() // 英文翻译器
		uni := ut.New(enT, zhT, enT)

		// locale 通常取决于 http 请求头的 'Accept-Language'
		var ok bool
		// 也可以使用 uni.FindTranslator(...) 传入多个locale进行查找
		validators.Trans, ok = uni.GetTranslator(locale)
		if !ok {
			return validators, fmt.Errorf("uni.GetTranslator(%s) failed", locale)
		}

		// 注册翻译器
		switch locale {
		case "en":
			err = enTranslations.RegisterDefaultTranslations(v, validators.Trans)
		case "zh":
			err = zhTranslations.RegisterDefaultTranslations(v, validators.Trans)
		default:
			err = enTranslations.RegisterDefaultTranslations(v, validators.Trans)
		}
	}
	return
}

// FetchGinValidatorError 获取gin验证器错误
func FetchGinValidatorError(err error, trans ut.Translator) interface{} {
	errs, ok := err.(validator.ValidationErrors)
	if ok {
		return errs.Translate(trans)
	}
	return nil
}
