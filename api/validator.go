package api

import (
	"github.com/Iowel/course-simple-bank/util"
	"github.com/go-playground/validator/v10"
)

// Валидационная функция для поля с валютой, которая проверяет, поддерживается ли указанная валюта в системе
var validCurrency validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if currency, ok := fieldLevel.Field().Interface().(string); ok {
		return util.IsSupportedCurrency(currency)
	}
	return false
}
