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
	case *Multi:
		errs := make([]wrapped, 0, len(err.errs))
		for _, err := range err.errs {
			errs = append(errs, wrapped{err})
		}

		return json.Marshal(errs)
		
	case json.Marshaler:
		return err.MarshalJSON()

	case fmt.Stringer:
		return json.Marshal(err.String())
	}

	return json.Marshal(e.Error())
}
