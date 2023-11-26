package errors

import (
	"errors"
	"fmt"
	"strings"
)

// Error list.
var (
	ErrInternalDB           = errors.New("internal database error")
	ErrInternalCache        = errors.New("internal cache error")
	ErrInternalServer       = errors.New("internal server error")
	ErrInvalidDBFormat      = errors.New("invalid db address")
	ErrInvalidRequestFormat = errors.New("invalid request format")
	ErrInvalidRequestData   = errors.New("invalid request data")
	ErrInvalidMessageType   = errors.New("invalid message type")
	ErrInvalidToken         = errors.New("invalid token or already expired")
	ErrWikiaPageNotFound    = errors.New("wikia page not found")
	ErrInvalidID            = errors.New("invalid id")
	ErrVtuberNotFound       = errors.New("vtuber not found")
	ErrAgencyNotFound       = errors.New("agency not found")
	ErrChannelNotFound      = errors.New("channel not found")
	ErrUserNotFound         = errors.New("user not found")
)

// ErrRequiredField is error for missing field.
func ErrRequiredField(str string) error {
	return fmt.Errorf("required field %s", str)
}

// ErrGTField is error for greater than field.
func ErrGTField(str, value string) error {
	return fmt.Errorf("field %s must be greater than %s", str, value)
}

// ErrGTEField is error for greater than or equal field.
func ErrGTEField(str, value string) error {
	return fmt.Errorf("field %s must be greater than or equal %s", str, value)
}

// ErrLTField is error for lower than field.
func ErrLTField(str, value string) error {
	return fmt.Errorf("field %s must be lower than %s", str, value)
}

// ErrLTEField is error for lower than or equal field.
func ErrLTEField(str, value string) error {
	return fmt.Errorf("field %s must be lower than or equal %s", str, value)
}

// ErrOneOfField is error for oneof field.
func ErrOneOfField(str, value string) error {
	return fmt.Errorf("field %s must be one of %s", str, strings.Join(strings.Split(value, " "), "/"))
}
