package errors

import (
	"encoding/json"
	"errors"
	"fmt"
)

func AsType[T error](err error) (T, bool) {
	var target T
	ok := errors.As(err, &target)
	return target, ok
}

type wrapped struct{ error }

func (e *wrapped) Error() string {
	return e.error.Error()
}

func (e *wrapped) String() string {
	return e.error.Error()
}

func (e *wrapped) Unwrap() error {
	return e.error
}

func (e wrapped) MarshalJSON() ([]byte, error) {
	switch err := e.error.(type) {
	case json.Marshaler:
		return err.MarshalJSON()

	case fmt.Stringer:
		return json.Marshal(err.String())

	case interface{ Unwrap() []error }:
		joined := Join(err.Unwrap()...)
		return json.Marshal(joined)

	}

	return json.Marshal(e.Error())
}
