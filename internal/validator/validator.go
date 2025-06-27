package validator

import (
	"slices"
	"strings"
	"unicode/utf8"
)

type Validator struct {
	FieldErrors map[string]string
}

func (v *Validator) IsValid() bool {
	return len(v.FieldErrors) == 0
}

func (v *Validator) AddFieldError(key, message string) {
	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	}

	if _, exists := v.FieldErrors[key]; !exists {
		v.FieldErrors[key] = message
	}
}

func (v *Validator) CheckField(key, message string, ok bool) {
	if !ok {
		v.AddFieldError(key, message)
	}
}

func NotBlank(s string) bool {
	return strings.TrimSpace(s) != ""
}

func MaxChars(s string, n int) bool {
	return utf8.RuneCountInString(s) <= n
}

func PermittedValues[T comparable](value T, permittedValues ...T) bool {
	return slices.Contains(permittedValues, value)
}
