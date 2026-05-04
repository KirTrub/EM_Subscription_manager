package apperrors

import (
	"errors"
	"fmt"
	"net/http"
)

type AppError struct {
	Code    int    `json:"-"`
	Message string `json:"error"`
	Err     error  `json:"-"`
}

func (e *AppError) Error() string {
	if e == nil {
		return "<nil>"
	}
	if e.Err == nil {
		return e.Message
	}
	return fmt.Sprintf("%s: %v", e.Message, e.Err)
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func New(code int, message string) *AppError {
	return &AppError{Code: code, Message: message}
}

func Wrap(err error, code int, message string) error {
	if err == nil {
		return nil
	}
	return &AppError{Code: code, Message: message, Err: err}
}

func NewBadRequest(message string) error {
	return New(http.StatusBadRequest, message)
}

func NewNotFound(message string) error {
	return New(http.StatusNotFound, message)
}

func NewConflict(message string) error {
	return New(http.StatusConflict, message)
}

func NewInternal(message string) error {
	return New(http.StatusInternalServerError, message)
}

func As(err error) (*AppError, bool) {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr, true
	}
	return nil, false
}

func GetCode(err error) int {
	if err == nil {
		return http.StatusOK
	}
	if appErr, ok := As(err); ok {
		return appErr.Code
	}
	return http.StatusInternalServerError
}

func GetMessage(err error) string {
	if err == nil {
		return ""
	}
	if appErr, ok := As(err); ok {
		return appErr.Message
	}
	return err.Error()
}
