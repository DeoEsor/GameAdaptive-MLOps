package customerrors

import (
	"errors"
	"net/http"
)

var (
	// ErrUnknownType ошбика сигнализирует о том, что был пойман неизвестный тип
	ErrUnknownType = errors.New("unkown type")
)

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
