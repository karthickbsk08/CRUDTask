package govalidatorpkg

import (
	"log"
	"reflect"
	"strings"

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

func CleanAndValidateStruct(s any) (any, error) {

	log.Println("s : ", s)
	val := reflect.ValueOf(s)
	typ := val.Type()
	log.Println("val : ", val)
	log.Println("typ : ", typ)

	for i := range val.NumField() {
		fieldVal := val.Field(i)
		fieldType := typ.Field(i)
		tag := fieldType.Tag.Get("valid")

		// Apply trimming
		if strings.Contains(tag, "trim") && fieldVal.Kind() == reflect.String {
			fieldVal.SetString(strings.TrimSpace(fieldVal.String()))
		}

		// Apply lowercasing
		if strings.Contains(tag, "lower") && fieldVal.Kind() == reflect.String {
			fieldVal.SetString(strings.ToLower(fieldVal.String()))
		}
	}
	err := validate.Struct(s)
	if err != nil {
		log.Println("err : ", err)
		return s, err
	}
	return s, nil
}
