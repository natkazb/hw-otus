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

var (
	ErrNotStruct            = errors.New("input must be a struct")
	ErrInvalidValidator     = errors.New("invalid validator format")
	ErrUnknownValidator     = errors.New("unknown validator")
	ErrINTMinInvalid        = errors.New("invalid min value")
	ErrINTMinLess           = errors.New("value is less than min")
	ErrINTMaxInvalid        = errors.New("invalid max value")
	ErrINTMaxGreater        = errors.New("value is greater than max")
	ErrINInvalid            = errors.New("invalid in value")
	ErrINNotInSet           = errors.New("value is not in allowed set")
	ErrSTRINGLenInvalid     = errors.New("invalid length value")
	ErrSTRINGLenNotMatch    = errors.New("string length does not match required length")
	ErrSTRINGRegexpInvalid  = errors.New("invalid regexp pattern")
	ErrSTRINGRegexpNotMatch = errors.New("string does not match pattern")
	ErrSTRINGInNotInSet     = errors.New("value is not in allowed set")
)

func (v ValidationErrors) Error() string {
	errStrings := make([]string, 0, len(v))
	for _, err := range v {
		errStrings = append(errStrings, fmt.Sprintf("%s: %v", err.Field, err.Err)) // nozero
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

	//nolint:exhaustive
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
	default: // Skip validation for unsupported types
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
		minVal, err := strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			return fmt.Errorf("%w: %s", ErrINTMinInvalid, parts[1])
		}
		if value < minVal {
			return fmt.Errorf("%w: %d is less than min %d", ErrINTMinLess, value, minVal)
		}
	case "max":
		maxVal, err := strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			return fmt.Errorf("%w: %s", ErrINTMaxInvalid, parts[1])
		}
		if value > maxVal {
			return fmt.Errorf("%w: %d is greater than max %d", ErrINTMaxGreater, value, maxVal)
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
			return fmt.Errorf("%w: %d is not in allowed set %s", ErrINNotInSet, value, parts[1])
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
			return fmt.Errorf("%w: %s", ErrSTRINGLenInvalid, parts[1])
		}
		if len(value) != length {
			return fmt.Errorf("%w: length %d does not match length %d", ErrSTRINGLenNotMatch, len(value), length)
		}
	case "regexp":
		re, err := regexp.Compile(parts[1])
		if err != nil {
			return fmt.Errorf("%w: %s", ErrSTRINGRegexpInvalid, parts[1])
		}
		if !re.MatchString(value) {
			return fmt.Errorf("%w: %s", ErrSTRINGRegexpNotMatch, parts[1])
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
			return fmt.Errorf("%w: %s is not in set %s", ErrSTRINGInNotInSet, value, parts[1])
		}
	default:
		return fmt.Errorf("%w: %s", ErrUnknownValidator, parts[0])
	}
	return nil
}
