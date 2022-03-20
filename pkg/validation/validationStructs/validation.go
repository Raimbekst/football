package validationStructs

import (
	"carWash/pkg/logger"
	"github.com/go-playground/validator/v10"
	"reflect"
	"strings"
)

func ValidateStruct(inp interface{}) (bool, map[string]string) {
	validate := validator.New()
	err := validate.Struct(inp)

	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			logger.Errorf("validationStructs.ValidateStruct: %v", err)
		}
		e := make(map[string]string)
		reflected := reflect.ValueOf(inp)

		for _, err := range err.(validator.ValidationErrors) {

			field, _ := reflected.Type().FieldByName(err.StructField())

			var name string
			if name = field.Tag.Get("json"); name == "" {
				name = strings.ToLower(err.StructField())
			}

			switch err.Tag() {
			case "required":
				e[name] = "The " + name + " is required."
				break
			case "e164":
				e[name] = "The " + name + " is invalid." + " Must be as +7707-404-21-42."
			default:
				e[name] = name + " is invalid type."
				break
			}

		}
		return false, e

	}
	return true, nil

}
