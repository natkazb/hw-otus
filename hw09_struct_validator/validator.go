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

var ErrNotStruct = errors.New("input must be a struct")
var ErrInvalidValidator = errors.New("invalid validator format")
var ErrUnknownValidator = errors.New("unknown validator")
var ErrINTMinInvalid = errors.New("invalid min value")
var ErrINTMinLess = errors.New("value is less than min")
var ErrINTMaxInvalid = errors.New("invalid max value")
var ErrINTMaxGreater = errors.New("value is greater than min")
var ErrINInvalid = errors.New("invalid in value")
var ErrINNotInSet = errors.New("value is not in allowed set")

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
		return ErrNotStruct
	}

	var validationErrors ValidationErrors

	for i := 0; i < typeV.NumField(); i++ {
		field := typeV.Field(i)
		validateTag, ok := field.Tag.Lookup("validate")
		if !ok { // нет тега, пропускаем
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
		return fmt.Errorf("%w: %s", ErrInvalidValidator, validator)
	}

	switch parts[0] {
	case "min":
		min, err := strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			return fmt.Errorf("%w: %s", ErrINTMinInvalid, parts[1])
		}
		if value < min {
			return fmt.Errorf("%w: %d is less than min %d", ErrINTMinLess, value, min)
		}
	case "max":
		max, err := strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			return fmt.Errorf("%w: %s", ErrINTMaxInvalid, parts[1])
		}
		if value > max {
			return fmt.Errorf("%w: %d is greater than max %d", ErrINTMaxGreater, value, max)
		}
	case "in":
		values := strings.Split(parts[1], ",")
		valid := false
		for _, v := range values {
			num, err := strconv.ParseInt(v, 10, 64)
			if err != nil {
				return fmt.Errorf("%w: %s", ErrINInvalid, v)
			}
			if value == num {
				valid = true
				break
			}
		}
		if !valid {
			return fmt.Errorf("%w: %d is not in allowed set %s", ErrINTMaxGreater, value, parts[1])
		}
	default:
		return fmt.Errorf("%w: %s", ErrUnknownValidator, parts[0])
	}
	return nil
}

func validateString(value string, validator string) error {
	parts := strings.Split(validator, ":")
	if len(parts) != 2 {
		return fmt.Errorf("%w: %s", ErrInvalidValidator, validator)
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
		return fmt.Errorf("%w: %s", ErrUnknownValidator, parts[0])
	}
	return nil
}
