package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	var errStrings []string
	for _, err := range v {
		errStrings = append(errStrings, fmt.Sprintf("%s: %v", err.Field, err.Err))
	}
	return strings.Join(errStrings, "; ")
}

func Validate(v interface{}) error {
	valueV := reflect.ValueOf(v)
	typeV := valueV.Type()

	if typeV.Kind() != reflect.Struct {
		return errors.New("input must be a struct")
	}

	var validationErrors ValidationErrors

	for i := 0; i < typeV.NumField(); i++ {
		field := typeV.Field(i)
		validateTag, ok := field.Tag.Lookup("validate")
		if !ok { // нет тега, пропустим это поле
			continue
		}
		if errs := validateField(field.Name, valueV.Field(i), validateTag); errs != nil {
			validationErrors = append(validationErrors, errs...)
		}
	}

	if len(validationErrors) > 0 {
		return validationErrors
	}
	return nil
}

func validateField(fieldName string, value reflect.Value, rules string) ValidationErrors {
	validators := strings.Split(rules, "|")
	var errors ValidationErrors

	switch value.Kind() {
	case reflect.Int:
		for _, validator := range validators {
			if err := validateInt(value.Int(), validator); err != nil {
				errors = append(errors, ValidationError{Field: fieldName, Err: err})
			}
		}
	case reflect.String:
		for _, validator := range validators {
			if err := validateString(value.String(), validator); err != nil {
				errors = append(errors, ValidationError{Field: fieldName, Err: err})
			}
		}
	case reflect.Slice:
		for i := 0; i < value.Len(); i++ {
			elem := value.Index(i)
			if errs := validateField(fmt.Sprintf("%s[%d]", fieldName, i), elem, rules); errs != nil {
				errors = append(errors, errs...)
			}
		}
	}

	if len(errors) > 0 {
		return errors
	}
	return nil
}

func validateInt(value int64, validator string) error {
	parts := strings.Split(validator, ":")
	if len(parts) != 2 {
		return fmt.Errorf("invalid validator format: %s", validator)
	}

	switch parts[0] {
	case "min":
		min, err := strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid min value: %s", parts[1])
		}
		if value < min {
			return fmt.Errorf("value %d is less than min %d", value, min)
		}
	case "max":
		max, err := strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid max value: %s", parts[1])
		}
		if value > max {
			return fmt.Errorf("value %d is greater than max %d", value, max)
		}
	case "in":
		values := strings.Split(parts[1], ",")
		valid := false
		for _, v := range values {
			num, err := strconv.ParseInt(v, 10, 64)
			if err != nil {
				return fmt.Errorf("invalid in value: %s", v)
			}
			if value == num {
				valid = true
				break
			}
		}
		if !valid {
			return fmt.Errorf("value %d is not in allowed set %s", value, parts[1])
		}
	default:
		return fmt.Errorf("unknown validator: %s", parts[0])
	}
	return nil
}

func validateString(value string, validator string) error {
	parts := strings.Split(validator, ":")
	if len(parts) != 2 {
		return fmt.Errorf("invalid validator format: %s", validator)
	}

	switch parts[0] {
	case "len":
		length, err := strconv.Atoi(parts[1])
		if err != nil {
			return fmt.Errorf("invalid length value: %s", parts[1])
		}
		if len(value) != length {
			return fmt.Errorf("string length %d does not match required length %d", len(value), length)
		}
	case "regexp":
		re, err := regexp.Compile(parts[1])
		if err != nil {
			return fmt.Errorf("invalid regexp pattern: %s", parts[1])
		}
		if !re.MatchString(value) {
			return fmt.Errorf("string does not match pattern %s", parts[1])
		}
	case "in":
		values := strings.Split(parts[1], ",")
		valid := false
		for _, v := range values {
			if value == v {
				valid = true
				break
			}
		}
		if !valid {
			return fmt.Errorf("value %s is not in allowed set %s", value, parts[1])
		}
	default:
		return fmt.Errorf("unknown validator: %s", parts[0])
	}
	return nil
}
