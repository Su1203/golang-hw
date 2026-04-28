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
	if len(v) == 0 {
		return "no validation errors"
	}
	var sb strings.Builder
	sb.WriteString("validation errors: ")
	for i, err := range v {
		if i > 0 {
			sb.WriteString("; ")
		}
		sb.WriteString(fmt.Sprintf("%s: %v", err.Field, err.Err))
	}
	return sb.String()
}

var (
	ErrInvalidValidateTag = errors.New("invalid validate tag")
	ErrUnsupportedType    = errors.New("unsupported type")
	ErrNotStruct          = errors.New("input is not a struct")
)

func Validate(v interface{}) error {
	if v == nil {
		return ErrNotStruct
	}

	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return ErrNotStruct
	}

	var validationErrors ValidationErrors
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		fieldValue := val.Field(i)

		if !field.IsExported() {
			continue
		}

		validateTag := field.Tag.Get("validate")
		if validateTag == "" {
			continue
		}

		errs, err := validateField(field.Name, fieldValue, validateTag)
		if err != nil {
			return fmt.Errorf("field %s: %w", field.Name, err)
		}
		validationErrors = append(validationErrors, errs...)
	}

	if len(validationErrors) > 0 {
		return validationErrors
	}
	return nil
}

func validateField(fieldName string, fieldValue reflect.Value, validateTag string) ([]ValidationError, error) {
	var errors []ValidationError

	rules := strings.Split(validateTag, "|")

	switch fieldValue.Kind() {
	case reflect.String:
		for _, rule := range rules {
			if err := validateString(fieldValue.String(), rule); err != nil {
				errors = append(errors, ValidationError{Field: fieldName, Err: err})
			}
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		for _, rule := range rules {
			if err := validateInt(fieldValue.Int(), rule); err != nil {
				errors = append(errors, ValidationError{Field: fieldName, Err: err})
			}
		}
	case reflect.Slice:
		for i := 0; i < fieldValue.Len(); i++ {
			elem := fieldValue.Index(i)
			for _, rule := range rules {
				var err error
				switch elem.Kind() {
				case reflect.String:
					err = validateString(elem.String(), rule)
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					err = validateInt(elem.Int(), rule)
				default:
					return nil, fmt.Errorf("%w: %v", ErrUnsupportedType, elem.Kind())
				}
				if err != nil {
					errors = append(errors, ValidationError{Field: fieldName, Err: err})
				}
			}
		}
	default:
	}

	return errors, nil
}

func validateString(value, rule string) error {
	parts := strings.SplitN(rule, ":", 2)
	if len(parts) != 2 {
		return fmt.Errorf("%w: %s", ErrInvalidValidateTag, rule)
	}

	validatorName := parts[0]
	validatorValue := parts[1]

	switch validatorName {
	case "len":
		expectedLen, err := strconv.Atoi(validatorValue)
		if err != nil {
			return fmt.Errorf("%w: invalid len value: %s", ErrInvalidValidateTag, validatorValue)
		}
		if len(value) != expectedLen {
			return fmt.Errorf("length must be %d", expectedLen)
		}
	case "regexp":
		re, err := regexp.Compile(validatorValue)
		if err != nil {
			return fmt.Errorf("%w: invalid regexp: %s", ErrInvalidValidateTag, validatorValue)
		}
		if !re.MatchString(value) {
			return fmt.Errorf("does not match regexp %s", validatorValue)
		}
	case "in":
		allowedValues := strings.Split(validatorValue, ",")
		found := false
		for _, allowed := range allowedValues {
			if value == allowed {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("must be one of [%s]", validatorValue)
		}
	default:
		return fmt.Errorf("%w: unknown validator %s", ErrInvalidValidateTag, validatorName)
	}

	return nil
}

func validateInt(value int64, rule string) error {
	parts := strings.SplitN(rule, ":", 2)
	if len(parts) != 2 {
		return fmt.Errorf("%w: %s", ErrInvalidValidateTag, rule)
	}

	validatorName := parts[0]
	validatorValue := parts[1]

	switch validatorName {
	case "min":
		minValue, err := strconv.ParseInt(validatorValue, 10, 64)
		if err != nil {
			return fmt.Errorf("%w: invalid min value: %s", ErrInvalidValidateTag, validatorValue)
		}
		if value < minValue {
			return fmt.Errorf("must be at least %d", minValue)
		}
	case "max":
		maxValue, err := strconv.ParseInt(validatorValue, 10, 64)
		if err != nil {
			return fmt.Errorf("%w: invalid max value: %s", ErrInvalidValidateTag, validatorValue)
		}
		if value > maxValue {
			return fmt.Errorf("must be at most %d", maxValue)
		}
	case "in":
		allowedValues := strings.Split(validatorValue, ",")
		found := false
		for _, allowed := range allowedValues {
			allowedInt, err := strconv.ParseInt(allowed, 10, 64)
			if err != nil {
				return fmt.Errorf("%w: invalid in value: %s", ErrInvalidValidateTag, allowed)
			}
			if value == allowedInt {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("must be one of [%s]", validatorValue)
		}
	default:
		return fmt.Errorf("%w: unknown validator %s", ErrInvalidValidateTag, validatorName)
	}

	return nil
}
