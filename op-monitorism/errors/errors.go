package errors

import (
	"errors"
	"fmt"
)

type ErrorCode string

const (
	ErrCodeNetwork       ErrorCode = "NETWORK"
	ErrCodeValidation    ErrorCode = "VALIDATION"
	ErrCodeConfiguration ErrorCode = "CONFIG"
	ErrCodeInternal      ErrorCode = "INTERNAL"
	ErrCodeBlockchain    ErrorCode = "BLOCKCHAIN"
	ErrCodeNotFound      ErrorCode = "NOT_FOUND"
	ErrCodeTimeout       ErrorCode = "TIMEOUT"
	ErrCodePermission    ErrorCode = "PERMISSION"
)

type MonitorError struct {
	Code    ErrorCode
	Message string
	Cause   error
	Details map[string]interface{}
}

func (e *MonitorError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

func (e *MonitorError) Unwrap() error {
	return e.Cause
}

func New(code ErrorCode, message string) *MonitorError {
	return &MonitorError{
		Code:    code,
		Message: message,
		Details: make(map[string]interface{}),
	}
}

func Wrap(err error, code ErrorCode, message string) *MonitorError {
	if err == nil {
		return nil
	}
	return &MonitorError{
		Code:    code,
		Message: message,
		Cause:   err,
		Details: make(map[string]interface{}),
	}
}

func (e *MonitorError) WithDetail(key string, value interface{}) *MonitorError {
	if e.Details == nil {
		e.Details = make(map[string]interface{})
	}
	e.Details[key] = value
	return e
}

func Is(err error, code ErrorCode) bool {
	var me *MonitorError
	if errors.As(err, &me) {
		return me.Code == code
	}
	return false
}

func GetCode(err error) ErrorCode {
	var me *MonitorError
	if errors.As(err, &me) {
		return me.Code
	}
	return ErrCodeInternal
}

/*
Predefined errors for common scenarios
*/
var (
	ErrConnectionFailed = New(ErrCodeNetwork, "connection failed")
	ErrInvalidBlock     = New(ErrCodeValidation, "invalid block")
	ErrMissingConfig    = New(ErrCodeConfiguration, "missing configuration")
	ErrInvalidAddress   = New(ErrCodeValidation, "invalid address")
	ErrRPCTimeout       = New(ErrCodeTimeout, "RPC timeout")
)
