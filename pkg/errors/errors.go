package errors

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

type Error struct {
	Kind     Kind
	Message  string
	Cause    error
	Metadata map[string]any
}

func New(kind Kind, message string, cause error, metadata map[string]any) error {
	return &Error{
		Kind:     kind,
		Message:  message,
		Cause:    cause,
		Metadata: metadata,
	}
}

func (err *Error) Error() string {
	return err.Message + `: ` + err.Cause.Error()
}

func (err *Error) Unwrap() error {
	return err.Cause
}

func (err Error) MarshalJSON() ([]byte, error) {
	var efj errorForJSON
	efj = errorForJSON(err)

	if efj.Metadata == nil {
		efj.Metadata = map[string]any{}
	}

	return json.Marshal(efj)
}

func (err *Error) UnmarshalJSON(buf []byte) error {
	var efj errorForJSON
	if err := json.Unmarshal(buf, &efj); err != nil {
		return err
	}

	*err = Error(efj)
	return nil
}

type Kind int

const (
	InvalidInput Kind = iota
	Unauthorized
	PreconditionFailed
	NotFound
	Conflict

	Internal
	Unavailable
	Timeout
	BadGateway
)

var kindStrings = map[Kind]string{
	InvalidInput:       "invalid input",
	Unauthorized:       "unauthorized",
	PreconditionFailed: "precondition failed",
	NotFound:           "not found",
	Conflict:           "conflict",

	Internal:    "internal error",
	Unavailable: "unavailable",
	Timeout:     "timeout",
	BadGateway:  "bad gateway",
}

var stringKinds = invert(kindStrings)

func (k Kind) String() string {
	return kindStrings[k]
}

func (k Kind) MarshalJSON() ([]byte, error) {
	quoted := strconv.Quote(kindStrings[k])
	return []byte(quoted), nil
}

func (k *Kind) UnmarshalJSON(buf []byte) error {
	unquoted, err := strconv.Unquote(string(buf))
	if err != nil {
		return fmt.Errorf("kind must be a valid JSON string: %w", err)
	}

	*k = stringKinds[unquoted]
	return nil
}

func Join(errs ...error) error {
	return errors.Join(errs...)
}

type errorForJSON struct {
	Kind     Kind           `json:"kind"`
	Message  string         `json:"message"`
	Cause    error          `json:"cause,omitempty"`
	Metadata map[string]any `json:"metadata,omitempty"`
}

func invert[K, V comparable](m map[K]V) map[V]K {
	nm := make(map[V]K, len(m))
	for k, v := range m {
		nm[v] = k
	}

	return nm
}
