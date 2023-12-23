package customerrors

import (
	"errors"
	"net/http"
)

// ErrUnknownType ошбика сигнализирует о том, что был пойман неизвестный тип
var ErrUnknownType = errors.New("unkown type, register type in middleware registration")

type ErrCodes struct {
	Err  error
	Code int
}

func (e ErrCodes) Error() string {
	return e.Err.Error()
}

func (e ErrCodes) StatusCode() int {
	return e.Code
}

func CodesNotFound(err error) error {
	return ErrCodes{
		Err:  err,
		Code: http.StatusNotFound,
	}
}

func CodesBadRequest(err error) error {
	return ErrCodes{
		Err:  err,
		Code: http.StatusBadRequest,
	}
}
