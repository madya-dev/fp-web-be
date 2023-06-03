package validator

import (
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
)

func RequestBodyValidator(input interface{}) []string {
	english := en.New()
	uni := ut.New(english, english)

	trans, _ := uni.GetTranslator("en")

	validate := validator.New()
	_ = enTranslations.RegisterDefaultTranslations(validate, trans)

	var errorString []string

	err := validate.Struct(input)
	if err != nil {
		errs := err.(validator.ValidationErrors)

		for _, e := range errs {
			errorString = append(errorString, e.Translate(trans))
		}
	}

	return errorString
}
