package validators

import (
	"errors"
	"fmt"
	"reflect"
)

const (
	VALIDATION__ERROR_MAX_EXCEEDED    = "field %s can not be greater than %v"
	VALIDATION__ERROR_MAXLEN_EXCEEDED = "field %s can not be longer than %v"
	VALIDATION__ERROR_MIN             = "field %s can not be lesser than %v"
	VALIDATION__ERROR_MINLEN          = "field %s can not be shorter than %v"
	VALIDATION__ERROR_NOT_IN_SLICE    = "field %s does not match expected values"
	VALIDATION__ERROR_NOT_STRING      = "invalid type for field %s, exected to be a string"
	VALIDATION__ERROR_NOT_FLOAT		  = "invalid type for fied %s, exected to be a float"
	VALIDATION__ERROR_EMPTY 		  = "invalid value for field %s, can not be empty"
)

func Max(field, max int, fieldName string) error {
	if field > max {
		msg := fmt.Sprintf(VALIDATION__ERROR_MAX_EXCEEDED, fieldName, max)
		return errors.New(msg)
	}
	return nil
}

func Min(field, min int, fieldName string) error {
	if field < min {
		msg := fmt.Sprintf(VALIDATION__ERROR_MIN, fieldName, min)
		return errors.New(msg)
	}
	return nil
}

func MaxF(field, max float64, fieldName string) error {
	if field > max {
		msg := fmt.Sprintf(VALIDATION__ERROR_MAX_EXCEEDED, fieldName, max)
		return errors.New(msg)
	}
	return nil
}

func MinF(field, min float64, fieldName string) error{
	if field < min {
		msg := fmt.Sprintf(VALIDATION__ERROR_MINLEN, fieldName, min)
		return errors.New(msg)
	}
	return nil
}

func Maxlen(field string, max int, fieldName string) error {
	if len(field) > max {
		msg := fmt.Sprintf(VALIDATION__ERROR_MAXLEN_EXCEEDED, fieldName, max)
		return errors.New(msg)
	}
	return nil
}

func Minlen(field string, min int, fieldName string) error {
	if len(field) < min {
		msg := fmt.Sprintf(VALIDATION__ERROR_MINLEN, fieldName, min)
		return errors.New(msg)
	}
	return nil
}

func In(field string, accepted []string, fieldName string) error {
	valid := false

	for _, v := range accepted {
		if field == v {
			valid = true
		}
	}

	if !valid {
		msg := fmt.Sprintf(VALIDATION__ERROR_NOT_IN_SLICE, fieldName)
		return errors.New(msg)
	}
	return nil
}

func NotEmpty(field interface{}, fieldName string) error {
	if field == nil {
		msg := fmt.Sprintf(VALIDATION__ERROR_EMPTY, fieldName)
		return errors.New(msg)
	}
	return nil
}

func IsString(field interface{}, fieldName string) error {
	if reflect.TypeOf(field) != reflect.TypeOf("") {
		msg := fmt.Sprintf(VALIDATION__ERROR_NOT_STRING, fieldName)
		return errors.New(msg)
	}
	return nil
}

func IsFloat(field interface{}, fieldName string) error {
	var float float64
	if reflect.TypeOf(field) != reflect.TypeOf(float) {
		msg := fmt.Sprintf(VALIDATION__ERROR_NOT_FLOAT, fieldName)
		return errors.New(msg)
	}
	return nil
}