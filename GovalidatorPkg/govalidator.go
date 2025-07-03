package govalidatorpkg

import (
	"reflect"
	"strings"
	"tasks/helpers"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/schema"
)

var validate = validator.New()
var Decoder = schema.NewDecoder()

// func CleanAndValidateStruct(task any) error {
// 	return validate.Struct(task)
// }

//Example
// type Input struct {
// 	Name  string `valid:"required,trim"`
// 	Email string `valid:"required,email,trim,lower"`
// }

func CleanAndValidateStruct(pDebug *helpers.HelperStruct, pInputStruct any) error {
	pDebug.Log(helpers.Statement, "CleanAndValidateStruct(+)")

	var lErr error

	val := reflect.ValueOf(pInputStruct)
	if val.Kind() == reflect.Ptr {
		val = val.Elem() // Dereference to access fields
	}
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		fieldVal := val.Field(i)
		fieldType := typ.Field(i)
		tag := fieldType.Tag.Get("valid")

		// Apply trimming
		if strings.Contains(tag, "trim") && fieldVal.Kind() == reflect.String {
			// pDebug.Log(helpers.Statement, "Trimming field: ", fieldType.Name, "Value lenght : ", len(fieldVal.String()))
			fieldVal.SetString(strings.TrimSpace(fieldVal.String()))
			// pDebug.Log(helpers.Statement, "Trimmed value: ", fieldVal.String(), "Value lenght : ", len(fieldVal.String()))
		}

		// Apply lowercasing
		if strings.Contains(tag, "lower") && fieldVal.Kind() == reflect.String {
			fieldVal.SetString(strings.ToLower(fieldVal.String()))
		}
	}
	lErr = validate.Struct(pInputStruct)
	if lErr != nil {
		pDebug.Log(helpers.Elog, "lErr : ", lErr)
		return helpers.ErrReturn(lErr)
	}

	pDebug.Log(helpers.Statement, "CleanAndValidateStruct(-)")
	return nil
}
