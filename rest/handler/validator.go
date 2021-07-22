package handler

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
)

var (
	uni      *ut.UniversalTranslator
	validate *validator.Validate
	trans    ut.Translator
)

func init() {
	translator := zh.New()
	uni = ut.New(translator, translator)
	trans, _ = uni.GetTranslator("zh")
	validate = binding.Validator.Engine().(*validator.Validate)
	_ = zh_translations.RegisterDefaultTranslations(validate, trans)
}

//Translate Translate
func Translate(err error) string {
	var result string
	if fieldErr, ok := err.(validator.ValidationErrors); ok {
		for _, err := range fieldErr {
			result += err.Translate(trans) + ";"
			break
		}
	} else {
		result = err.Error()
	}
	return result
}
