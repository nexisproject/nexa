// Copyright (C) micros. 2025-present.
//
// Created at 2025-01-04, by liasica

package rest

import (
	zhLocale "github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhTranslation "github.com/go-playground/validator/v10/translations/zh"
)

type Validator struct {
	validator *validator.Validate
}

func NewValidator() *Validator {
	zh := zhLocale.New()
	uni := ut.New(zh, zh)
	trans, _ := uni.GetTranslator("zh")

	validate := validator.New()

	_ = zhTranslation.RegisterDefaultTranslations(validate, trans)

	return &Validator{validator: validate}
}

func (v *Validator) Validate(i any) error {
	// if err := v.validator.Struct(i); err != nil {
	// 	return NewError(http.StatusBadRequest, err.Error())
	// }
	// return nil
	return v.validator.Struct(i)
}
