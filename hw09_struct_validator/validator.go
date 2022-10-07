package hw09structvalidator

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type (
	// Все валидаторы объедены общим интерфейсом.
	FieldValidator interface {
		ValidateField(field reflect.StructField, value reflect.Value) ValidationError
	}

	// Валидатор проверяет вхождение значения в множество.
	InStringFieldValidator struct {
		in map[string]struct{}
	}

	// Валидатор проверяет длину строки.
	LengthStringFieldValidator struct {
		length int
	}

	// Валидатор проверяет соответствие регулярному выражению.
	RegexpStringFieldValidator struct {
		regexp *regexp.Regexp
	}

	// Валидатор проверяет вхождение значения в множество.
	InIntFieldValidator struct {
		in map[int]struct{}
	}

	// Валидатор проверяет что значение больше заданного в описании тега.
	MinIntFieldValidator struct {
		min int
	}

	// Валидатор проверяет что значение меньше заданного в описании тега.
	MaxIntFieldValidator struct {
		max int
	}

	// Ошибка валидации возникающая если проверка не прошла.
	ValidationError struct {
		Field string
		Err   error
	}

	// Ошибка валидации возникающая если при создании какого-то поля произошла ошибка.
	CreationError struct {
		Field string
		Err   error
	}

	// Ошибки валидации объединенные в массив.
	ValidationErrors []ValidationError
)

var (
	ErrIsNotStruct           = errors.New("is not struct")
	ErrHasUnknownValidator   = errors.New("has unknown validator")
	ErrHasInvalidValidator   = errors.New("has invalid validator")
	ErrWhanCreatingValidator = errors.New("when creating validator")

	ErrValidateLessValue       = errors.New("value is less than expected")
	ErrValidateOutOfScopeValue = errors.New("value out of scope")
	ErrValidateGreaterValue    = errors.New("value is greater than expected")
	ErrValidateStringLength    = errors.New("string length does not match expected")
	ErrValidateRegularMatch    = errors.New("value does not match regular expression")

	NilValidationError        ValidationError
	mappingCreateFuncWithKind = map[reflect.Kind]map[string]func(cond string) (FieldValidator, error){
		reflect.String: {
			"in":     NewInStringValidator,
			"len":    NewLengthStringValidator,
			"regexp": NewRegexpStringFieldValidator,
		},
		reflect.Int: {
			"in":  NewInIntValidator,
			"min": NewMinIntValidator,
			"max": NewMaxIntValidator,
		},
	}
)

func NewMinIntValidator(cond string) (FieldValidator, error) {
	var validator FieldValidator

	number, err := strconv.Atoi(cond)
	if err != nil {
		return validator, err
	}

	validator = &MinIntFieldValidator{
		min: number,
	}

	return validator, nil
}

func (fv *MinIntFieldValidator) ValidateField(field reflect.StructField, value reflect.Value) ValidationError {
	kind := field.Type.Kind()
	if kind == reflect.Slice || kind == reflect.Array {
		for _, v := range value.Interface().([]int64) {
			if v < int64(fv.min) {
				return ValidationError{
					Field: field.Name,
					Err:   ErrValidateLessValue,
				}
			}
		}
	}

	if value.Int() < int64(fv.min) {
		return ValidationError{
			Field: field.Name,
			Err:   ErrValidateLessValue,
		}
	}

	return NilValidationError
}

func NewInIntValidator(cond string) (FieldValidator, error) {
	var validator FieldValidator

	inValues := make(map[int]struct{})
	for _, value := range strings.Split(cond, ",") {
		val, err := strconv.Atoi(value)
		if err != nil {
			return validator, err
		}
		inValues[val] = struct{}{}
	}

	validator = &InIntFieldValidator{
		in: inValues,
	}

	return validator, nil
}

func (fv *InIntFieldValidator) ValidateField(field reflect.StructField, value reflect.Value) ValidationError {
	kind := field.Type.Kind()
	if kind == reflect.Slice || kind == reflect.Array {
		for _, v := range value.Interface().([]int) {
			if _, ok := fv.in[v]; !ok {
				return ValidationError{
					Field: field.Name,
					Err:   ErrValidateOutOfScopeValue,
				}
			}
		}
	}

	if _, ok := fv.in[int(value.Int())]; !ok {
		return ValidationError{
			Field: field.Name,
			Err:   ErrValidateOutOfScopeValue,
		}
	}

	return NilValidationError
}

func NewMaxIntValidator(cond string) (FieldValidator, error) {
	var validator FieldValidator

	m, err := strconv.Atoi(cond)
	if err != nil {
		return validator, err
	}

	validator = &MaxIntFieldValidator{
		max: m,
	}

	return validator, nil
}

func (fv *MaxIntFieldValidator) ValidateField(field reflect.StructField, value reflect.Value) ValidationError {
	kind := field.Type.Kind()
	if kind == reflect.Slice || kind == reflect.Array {
		for _, v := range value.Interface().([]int64) {
			if v > int64(fv.max) {
				return ValidationError{
					Field: field.Name,
					Err:   ErrValidateGreaterValue,
				}
			}
		}
	}

	if value.Int() > int64(fv.max) {
		return ValidationError{
			Field: field.Name,
			Err:   ErrValidateGreaterValue,
		}
	}

	return NilValidationError
}

func NewLengthStringValidator(cond string) (FieldValidator, error) {
	var validator FieldValidator

	length, err := strconv.Atoi(cond)
	if err != nil {
		return validator, err
	}

	validator = &LengthStringFieldValidator{
		length: length,
	}

	return validator, nil
}

func (fv *LengthStringFieldValidator) ValidateField(field reflect.StructField, value reflect.Value) ValidationError {
	kind := field.Type.Kind()
	if kind == reflect.Slice || kind == reflect.Array {
		for _, v := range value.Interface().([]string) {
			if len(v) != fv.length {
				return ValidationError{
					Field: field.Name,
					Err:   ErrValidateStringLength,
				}
			}
		}
	}

	if len(value.String()) != fv.length {
		return ValidationError{
			Field: field.Name,
			Err:   ErrValidateStringLength,
		}
	}

	return NilValidationError
}

func NewInStringValidator(values string) (FieldValidator, error) {
	var validator FieldValidator

	inValues := make(map[string]struct{})
	for _, val := range strings.Split(values, ",") {
		inValues[val] = struct{}{}
	}

	validator = &InStringFieldValidator{
		in: inValues,
	}

	return validator, nil
}

func (fv *InStringFieldValidator) ValidateField(field reflect.StructField, value reflect.Value) ValidationError {
	kind := field.Type.Kind()
	if kind == reflect.Slice || kind == reflect.Array {
		for _, v := range value.Interface().([]string) {
			if _, ok := fv.in[v]; !ok {
				return ValidationError{
					Field: field.Name,
					Err:   ErrValidateOutOfScopeValue,
				}
			}
		}
	}

	if _, ok := fv.in[value.String()]; !ok {
		return ValidationError{
			Field: field.Name,
			Err:   ErrValidateOutOfScopeValue,
		}
	}

	return NilValidationError
}

func NewRegexpStringFieldValidator(exp string) (FieldValidator, error) {
	var validator FieldValidator

	reg, err := regexp.Compile(exp)
	if err != nil {
		return validator, err
	}

	validator = &RegexpStringFieldValidator{
		regexp: reg,
	}

	return validator, nil
}

func (fv *RegexpStringFieldValidator) ValidateField(field reflect.StructField, value reflect.Value) ValidationError {
	kind := field.Type.Kind()
	if kind == reflect.Slice || kind == reflect.Array {
		for _, v := range value.Interface().([]string) {
			if !fv.regexp.MatchString(v) {
				return ValidationError{
					Field: field.Name,
					Err:   ErrValidateRegularMatch,
				}
			}
		}
	}

	if !fv.regexp.MatchString(value.String()) {
		return ValidationError{
			Field: field.Name,
			Err:   ErrValidateRegularMatch,
		}
	}

	return NilValidationError
}

func NewFieldValidators(f reflect.StructField, k reflect.Kind) ([]FieldValidator, error) {
	var validators []FieldValidator

	tags, ok := f.Tag.Lookup("validate")
	if !ok {
		return validators, nil
	}

	createFuncMapping, ok := mappingCreateFuncWithKind[k]
	if !ok {
		return validators, ErrHasUnknownValidator
	}
	for _, validationRules := range strings.Split(tags, "|") {
		if ruleCondition := strings.Split(validationRules, ":"); len(ruleCondition) == 2 {
			createFunc, ok := createFuncMapping[ruleCondition[0]]
			if !ok {
				return validators, ErrHasUnknownValidator
			}
			validator, err := createFunc(ruleCondition[1])
			if err != nil {
				return validators, CreationError{
					Field: f.Name,
					Err:   err,
				}
			}
			validators = append(validators, validator)
		} else {
			return validators, ErrHasInvalidValidator
		}
	}

	return validators, nil
}

func (e CreationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Err)
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Err)
}

func (v ValidationErrors) Error() string {
	var buffer bytes.Buffer

	if len(v) > 0 {
		buffer.WriteString("Errors:\n")
	}
	for _, err := range v {
		buffer.WriteString(fmt.Sprintf("- %s\n", err.Error()))
	}

	return buffer.String()
}

func Validate(v interface{}) error {
	// Переданное тип переданного значения не является валидной структурой или указателем на нее
	r := reflect.Indirect(reflect.ValueOf(v))
	if r.Type().Kind() != reflect.Struct {
		return ErrIsNotStruct
	}

	validators := make(map[string][]FieldValidator, 0)
	for i := 0; i < r.NumField(); i++ {
		structField := r.Type().Field(i)
		fieldValue := r.Field(i)
		kind := structField.Type.Kind()
		if kind == reflect.Slice || kind == reflect.Array {
			if fieldValue.Len() == 0 {
				continue
			}
			kind = fieldValue.Index(0).Type().Kind()
		}
		fieldValidators, err := NewFieldValidators(structField, kind)
		if err == nil && len(fieldValidators) > 0 {
			validators[structField.Name] = fieldValidators
		} else if err != nil {
			return err
		}
	}

	validationErrors := make(ValidationErrors, 0)
	for name, fieldValidators := range validators {
		field, _ := r.Type().FieldByName(name)
		value := r.FieldByName(name)
		for _, v := range fieldValidators {
			if err := v.ValidateField(field, value); err != NilValidationError {
				validationErrors = append(validationErrors, err)
			}
		}
	}

	if len(validationErrors) > 0 {
		return validationErrors
	}

	return nil
}
