package form

import "context"

type CustomError struct {
	Code   int
	Status int
	Err    string
}

func (e CustomError) Error() string {
	return e.Err
}

func NewCustomError(code, status int, err string) *CustomError {
	return &CustomError{
		Code:   code,
		Status: status,
		Err:    err,
	}
}

func IsCustomError(err error) (bool, *CustomError) {
	if e, ok := err.(*CustomError); ok {
		return true, e
	}
	return false, nil
}

type WithCodeError struct {
	Code int
	Err  string
}

func (e WithCodeError) Error() string {
	return e.Err
}

func NewWithCodeError(code int, err string) *WithCodeError {
	return &WithCodeError{
		Code: code,
		Err:  err,
	}
}

func IsWithCodeError(err error) (bool, *WithCodeError) {
	if e, ok := err.(*WithCodeError); ok {
		return true, e
	}
	return false, nil
}

type WithStatusError struct {
	Status int
	Err    string
}

func (e WithStatusError) Error() string {
	return e.Err
}

func NewWithStatusError(status int, err string) *WithStatusError {
	return &WithStatusError{
		Status: status,
		Err:    err,
	}
}

func IsWithStatusError(err error) (bool, *WithStatusError) {
	if e, ok := err.(*WithStatusError); ok {
		return true, e
	}
	return false, nil
}

type UnauthorizedError struct {
	err string
}

func (e UnauthorizedError) Error() string {
	return e.err
}

func NewUnauthorizedError(err string) *UnauthorizedError {
	return &UnauthorizedError{
		err: err,
	}
}

func IsUnauthorizedError(err error) bool {
	if _, ok := err.(*UnauthorizedError); ok {
		return true
	}
	return false
}

type UnprocessableError struct {
	err string
}

func (e UnprocessableError) Error() string {
	return e.err
}

func NewUnprocessableError(err string) *UnprocessableError {
	return &UnprocessableError{
		err: err,
	}
}

func IsUnprocessableError(err error) bool {
	if _, ok := err.(*UnprocessableError); ok {
		return true
	}
	return false
}

type InternalStateError struct {
	err string
}

func (e InternalStateError) Error() string {
	return e.err
}

func NewInternalStateError(err string) *InternalStateError {
	return &InternalStateError{
		err: err,
	}
}

func IsInternalStateError(err error) bool {
	if _, ok := err.(*InternalStateError); ok {
		return true
	}
	return false
}

type NotFoundError struct {
	err string
}

func (e NotFoundError) Error() string {
	return e.err
}

func NewNotFoundError(err string) *NotFoundError {
	return &NotFoundError{
		err: err,
	}
}

func IsNotFoundError(err error) bool {
	if _, ok := err.(*NotFoundError); ok {
		return true
	}
	return false
}

func IsContextCancel(err error) bool {
	return err == context.Canceled
}
