package main

import (
	"github.com/go-playground/validator"
)

func initValidator() *validator.Validate {
	// translator := en.New()
	// uni := ut.New(translator, translator)

	// // this is usually known or extracted from http 'Accept-Language' header
	// // also see uni.FindTranslator(...)
	// trans, found := uni.GetTranslator("en")
	// if !found {
	// 	log.Fatal("translator not found")
	// }

	v := validator.New()

	// if err := en_translations.RegisterDefaultTranslations(v, trans); err != nil {
	// 	log.Fatal(err)
	// }

	// _ = v.RegisterTranslation("required", trans, func(ut ut.Translator) error {
	// 	return ut.Add("required", "{0} is a required field", true) // see universal-translator for details
	// }, func(ut ut.Translator, fe validator.FieldError) string {
	// 	t, _ := ut.T("required", fe.Field())
	// 	return t
	// })

	return v
}
