package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

type validateHandlerFunc func(key string, value string, fieldName string, field interface{}) error

type stringValidator string

var (
	length       stringValidator = "len"
	regexpString stringValidator = "regexp"
	subsetString stringValidator = "in"
)

type numberValidator string

var (
	minimum      numberValidator = "min"
	maximum      numberValidator = "max"
	subsetNumber numberValidator = "in"
)

var (
	ErrInterfaceType       = errors.New("interface is not struct")
	ErrInterfaceConversion = errors.New("invalid type for conversion")

	ErrInvalidEmptyTag = errors.New("invalid tag: tag can't be empty")

	ErrStringLength = errors.New("length of string not equals number in tag")
	ErrStringRegexp = errors.New("string not equals regexp in tag")
	ErrStringIn     = errors.New("string is not contains in subset of tag")

	ErrNumberMax = errors.New("number is bigger than max")
	ErrNumberMin = errors.New("number is less than min")
	ErrNumberIn  = errors.New("number is not contains in subset of tag")
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	errorString := strings.Builder{}
	for _, err := range v {
		errorString.WriteString(fmt.Sprintf("Field: %s - Error: %v\n", err.Field, err.Err))
	}
	return errorString.String()
}

func Validate(v interface{}) error {
	var (
		validateTagName  = "validate"
		validationErrors ValidationErrors
		wg               sync.WaitGroup
		lock             sync.Mutex
	)
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Struct {
		validationErrors = append(validationErrors, ValidationError{
			Field: val.Type().Name(),
			Err:   ErrInterfaceType,
		})
		return validationErrors
	}
	for i := 0; i < val.NumField(); i++ {
		wg.Add(1)
		i := i
		go func() {
			defer wg.Done()
			field := val.Type().Field(i)
			tag, ok := field.Tag.Lookup(validateTagName)
			if !ok {
				return
			}
			switch field.Type.Kind() {
			case reflect.String:
				fieldVal := val.Field(i).String()
				if err := validateHandler(fieldVal, tag, field.Name, validateString); err != nil {
					lock.Lock()
					validationErrors = errHandler(err, validationErrors, field.Name)
					lock.Unlock()
				}
			case reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8, reflect.Int:
				fieldVal := val.Field(i).Int()
				if err := validateHandler(fieldVal, tag, field.Name, validateNumber); err != nil {
					lock.Lock()
					validationErrors = errHandler(err, validationErrors, field.Name)
					lock.Unlock()
				}
			case reflect.Slice:
				fieldVal := val.Field(i).Interface()
				if err := validateHandler(fieldVal, tag, field.Name, validateSlice); err != nil {
					lock.Lock()
					validationErrors = errHandler(err, validationErrors, field.Name)
					lock.Unlock()
				}
			case reflect.Struct:
				fieldVal := val.Field(i).Interface()
				err := Validate(fieldVal)
				lock.Lock()
				validationErrors = errHandler(err, validationErrors, field.Name)
				lock.Unlock()
			default:
				return
			}
		}()
	}
	wg.Wait()
	return validationErrors
}

func errHandler(err error, validationErrors ValidationErrors, fieldName string) ValidationErrors {
	var valErr ValidationErrors
	if errors.As(err, &valErr) {
		validationErrors = append(validationErrors, valErr...)
	} else {
		validationErrors = append(validationErrors, ValidationError{
			Field: fieldName,
			Err:   err,
		})
	}
	return validationErrors
}

func validateHandler(field interface{}, fullCondition, fieldName string, validateHandler validateHandlerFunc) error {
	var (
		conditionsSeparator = "|"
		valErrors           = make(ValidationErrors, 0)
	)
	if fullCondition == "" || field == nil {
		return nil
	}
	conditions := strings.Split(fullCondition, conditionsSeparator)
	for _, cond := range conditions {
		key, value, err := getValidationPair(cond)
		if err != nil {
			valErrors = append(valErrors, ValidationError{
				Field: fieldName,
				Err:   err,
			})
			continue
		}
		if err = validateHandler(key, value, fieldName, field); err != nil {
			valErrors = errHandler(err, valErrors, fieldName)
		}
	}
	return valErrors
}

func validateString(key, value, fieldName string, field interface{}) error {
	var errors ValidationErrors
	str, ok := field.(string)
	if !ok {
		errors = append(errors, ValidationError{
			Field: fieldName,
			Err:   ErrInterfaceConversion,
		})
		return errors
	}

	switch key {
	case string(length):
		lengthNumber, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			errors = append(errors, ValidationError{
				Field: fieldName,
				Err:   err,
			})
			return errors
		}

		if len(str) != int(lengthNumber) {
			errors = append(errors, ValidationError{
				Field: fieldName,
				Err:   ErrStringLength,
			})
		}
	case string(regexpString):
		reg, err := regexp.Compile(value)
		if err != nil {
			errors = append(errors, ValidationError{
				Field: fieldName,
				Err:   err,
			})
		}
		if !reg.MatchString(str) {
			errors = append(errors, ValidationError{
				Field: fieldName,
				Err:   ErrStringRegexp,
			})
		}
	case string(subsetString):
		if !strings.Contains(value, str) {
			errors = append(errors, ValidationError{
				Field: fieldName,
				Err:   ErrStringIn,
			})
		}
	}
	return errors
}

func validateNumber(key, value, fieldName string, field interface{}) error {
	var errors ValidationErrors
	number, ok := field.(int64)
	if !ok {
		errors = append(errors, ValidationError{
			Field: fieldName,
			Err:   ErrInterfaceConversion,
		})
		return errors
	}

	switch key {
	case string(maximum):
		max, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			errors = append(errors, ValidationError{
				Field: fieldName,
				Err:   err,
			})
			return errors
		}
		if number > max {
			errors = append(errors, ValidationError{
				Field: fieldName,
				Err:   ErrNumberMax,
			})
		}
	case string(minimum):
		min, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			errors = append(errors, ValidationError{
				Field: fieldName,
				Err:   err,
			})
			return errors
		}
		if number < min {
			errors = append(errors, ValidationError{
				Field: fieldName,
				Err:   ErrNumberMin,
			})
		}
	case string(subsetNumber):
		if !strings.Contains(value, strconv.Itoa(int(number))) {
			errors = append(errors, ValidationError{
				Field: fieldName,
				Err:   ErrNumberIn,
			})
		}
	}
	return errors
}

func validateSlice(key, value, fieldName string, field interface{}) error {
	var validateErrors ValidationErrors
	switch field := field.(type) {
	case []int64:
		for _, elem := range field {
			err := validateNumber(key, value, fieldName, elem)
			var valErr ValidationErrors
			if errors.As(err, &valErr) {
				validateErrors = append(validateErrors, valErr...)
			} else {
				valErr = append(valErr, ValidationError{
					Field: fieldName,
					Err:   err,
				})
			}
		}
	case []string:
		for _, elem := range field {
			err := validateString(key, value, fieldName, elem)
			var valErr ValidationErrors
			if errors.As(err, &valErr) {
				validateErrors = append(validateErrors, valErr...)
			} else {
				valErr = append(valErr, ValidationError{
					Field: fieldName,
					Err:   err,
				})
			}
		}
	default:
		validateErrors = append(validateErrors, ValidationError{
			Field: fieldName,
			Err:   ErrInterfaceConversion,
		})
	}
	return validateErrors
}

// getValidationPair возвращает ключ-значение тэга `validate`
//
// 	Пример: len:32
// 	key = len, value = 32
func getValidationPair(cond string) (string, string, error) {
	keyValueSeparator := ":"
	splitCond := strings.Split(cond, keyValueSeparator)
	if len(splitCond) == 1 {
		return "", "", ErrInvalidEmptyTag
	} else if len(splitCond) == 2 && splitCond[1] == "" {
		return "", "", ErrInvalidEmptyTag
	}
	key := splitCond[0]
	value := splitCond[1]
	return key, value, nil
}
