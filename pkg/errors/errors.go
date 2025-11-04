package errors

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type Error struct {
	Kind    Kind
	Title   string
	Message string
	Cause   error
}

func New(kind Kind, title, message string, cause error) error {
	return &Error{
		Kind:    kind,
		Title:   title,
		Message: message,
		Cause:   cause,
	}
}

func (err *Error) Error() string {
	if err.Cause != nil {
		return err.Message + `: ` + err.Cause.Error()
	}

	return err.Message
}

func (err *Error) Unwrap() error {
	return err.Cause
}

func (err Error) MarshalJSON() ([]byte, error) {
	var efj errorForJSON
	efj = errorForJSON(err)

	if efj.Cause != nil {
		efj.Cause = &wrapped{efj.Cause}
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
	client_errors_start Kind = iota

	InvalidInput
	Unauthorized
	Forbidden
	PreconditionFailed
	NotFound
	Conflict

	client_errors_end
	internal_errors_start

	Internal
	Unavailable
	Timeout
	BadGateway

	internal_errors_end
)

var kindStrings = map[Kind]string{
	InvalidInput:       "invalid input",
	Unauthorized:       "unauthorized",
	Forbidden:          "forbidden",
	PreconditionFailed: "precondition failed",
	NotFound:           "not found",
	Conflict:           "conflict",

	Internal:    "internal error",
	Unavailable: "unavailable",
	Timeout:     "timeout",
	BadGateway:  "bad gateway",
}

var stringKinds = invert(kindStrings)

func (k Kind) IsClient() bool {
	return client_errors_start < k && k < client_errors_end
}

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

type errorForJSON struct {
	Kind    Kind   `json:"kind"`
	Title   string `json:"title"`
	Message string `json:"message"`
	Cause   error  `json:"cause,omitempty"`
}

func invert[K, V comparable](m map[K]V) map[V]K {
	nm := make(map[V]K, len(m))
	for k, v := range m {
		nm[v] = k
	}

	return nm
}
