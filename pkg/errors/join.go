package errors

import (
	"encoding/json"
	"strings"
)

type Multi struct{ errs []error }

func Join(errs ...error) error {
	var n int
	for _, err := range errs {
		if err != nil {
			n++
		}
	}
	if n == 0 {
		return nil
	}

	merr := Multi{errs: make([]error, 0, n)}
	for _, err := range errs {
		if err != nil {
			merr.errs = append(merr.errs, err)
		}
	}

	return &merr
}

func (e *Multi) Error() string {
	if e == nil || len(e.errs) == 0 {
		return "<nil>"
	}

	if len(e.errs) == 1 {
		return e.errs[0].Error()
	}

	var b strings.Builder
	b.WriteString(e.errs[0].Error())
	for _, err := range e.errs[1:] {
		b.WriteByte('\n')
		b.WriteString(err.Error())
	}

	return b.String()
}

func (e Multi) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.errs)
}

func (e *Multi) UnmarshalJSON(buf []byte) error {
	var errs []error
	if err := json.Unmarshal(buf, &errs); err != nil {
		return err
	}

	e.errs = errs
	return nil
}
