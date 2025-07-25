package validator

import (
	"net/mail"
	"slices"
	"strings"
	"unicode/utf8"
)

type Validator struct {
	NonFieldErrors []string
	FieldErrors    map[string]string
}

func (v *Validator) IsValid() bool {
	return len(v.FieldErrors) == 0 && len(v.NonFieldErrors) == 0
}

func (v *Validator) AddFieldError(key, message string) {
	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	}

	if _, exists := v.FieldErrors[key]; !exists {
		v.FieldErrors[key] = message
	}
}

func (v *Validator) AddNonFieldError(message string) {
	v.NonFieldErrors = append(v.NonFieldErrors, message)
}

func (v *Validator) CheckField(key, message string, ok bool) {
	if !ok {
		v.AddFieldError(key, message)
	}
}

func NotBlank(s string) bool {
	return strings.TrimSpace(s) != ""
}

func MinChars(s string, n int) bool {
	return utf8.RuneCountInString(s) >= n
}

func MaxChars(s string, n int) bool {
	return utf8.RuneCountInString(s) <= n
}

func PermittedValues[T comparable](value T, permittedValues ...T) bool {
	return slices.Contains(permittedValues, value)
}

func ValidEmailAddr(s string) bool {
	_, err := mail.ParseAddress(s)
	if err != nil {
		return false
	}
	return true
}
