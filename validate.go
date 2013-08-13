package mgorm

import (
	"reflect"
)

type IValidator interface {
	Validate() bool
}

func NewValidator(errorHandler IErrorHandler) IValidator {
	validater := new(Validator)
	validater.errorHandler = errorHandler
	return validater
}

type ValidateFn func(fieldValue reflect.Value, fieldType reflect.StructField) error

type Validator struct {
	errorHandler IErrorHandler
}

func (self *Validator) Validate() bool {
	refType := reflect.TypeOf(self.errorHandler)
	refValue := reflect.ValueOf(self.errorHandler)

	if refType.Kind() == reflect.Ptr {
		refType = refType.Elem()
		refValue = refValue.Elem()
	}

	numField := refType.NumField()

	var hasError bool
	for i := 0; i < numField; i++ {
		fieldType := refType.Field(i)
		fieldValue := refValue.Field(i)

		var err error
		if reflect.String == fieldType.Type.Kind() {
			tag := fieldType.Tag.Get("rules")
			switch tag {
			case "email":
				err = Validate(fieldValue, fieldType, EmailValidator)
			case "url":
				err = Validate(fieldValue, fieldType, UrlValidator)
			default:
				//do nothing
			}

			if nil != err {
				self.errorHandler.AddError(err.Error())
				hasError = true
			}
		}

		if reflect.Ptr != fieldType.Type.Kind() {
			continue
		}

		if v, ok := fieldValue.Interface().(IValidator); ok {
			if !v.Validate() {
				hasError = true
			}
		}

		if v, ok := fieldValue.Interface().(IErrorHandler); ok {
			if !NewValidator(v).Validate() {
				hasError = true
			}
		}
	}

	return hasError
}

func Validate(fieldValue reflect.Value, fieldType reflect.StructField, fn ValidateFn) error {
	return fn(fieldValue, fieldType)
}
