package main

import (
	"gopkg.in/go-playground/validator.v9"
	"unicode"
)

func ExcludeControl(fl validator.FieldLevel) bool {
	for _, c := range fl.Field().String() {
		if unicode.IsControl(c) {
			return false
		}
	}
	return true
}

func Validate(i interface{}) (field, tag string) {
	v := validator.New()
	v.RegisterValidation("excludes_control", ExcludeControl)
	err := v.Struct(i)

	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			field = err.Field()
			tag = err.Tag()
			break
		}
	}
	return
}
