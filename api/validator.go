package api

import (
	"github.com/forabbie/vank-app/util"
	"github.com/go-playground/validator/v10"
)

var validCurrency validator.Func = func(fl validator.FieldLevel) bool {
	if currency, ok := fl.Field().Interface().(string); ok {
		if currency == "USD" || currency == "EUR" {
			return util.IsSupportedCurrency(currency)
		}
	}
	return false
}
