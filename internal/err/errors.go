package err

import (
	"errors"
	"fmt"
	"net/http"
)

type ErrStatusCode interface {
	error
	Code() int
}

type errStatus struct {
	error
	code int
}

func (e errStatus) Unwrap() error {
	return e.error
}

func (e errStatus) Error() string {
	return fmt.Sprintf("[%d] %v", e.Code(), e.error)
}

func (e errStatus) Code() int {
	return e.code
}

func New(code int, err error) error {
	if err == nil {
		err = errors.New(http.StatusText(code))
	}

	return errStatus{err, code}
}
