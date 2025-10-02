package main

import (
	"encoding/json"
	"fmt"
	"regexp"
)

type Opt[T any] struct {
	Val  T
	Some bool
}

func New[T any](content T) Opt[T] {
	return Opt[T]{Val: content, Some: true}
}

func NewNone[T any]() Opt[T] {
	return Opt[T]{}
}

var (
	_JsonWhiteSpace = `[\x20\x09\x0A\x0D]*`
	_JsonNull       = `null`

	_JsonNullPattern = regexp.MustCompile(`^` + _JsonWhiteSpace + _JsonNull + _JsonWhiteSpace + `$`)
)

func (o Opt[T]) String() string {
	if !o.Some {
		return "<none>"
	}

	return fmt.Sprint(o.Val)
}

func (o Opt[T]) MarshalJSON() ([]byte, error) {
	if !o.Some {
		return []byte("null"), nil
	}

	return json.Marshal(o.Val)
}

func (o *Opt[T]) UnmarshalJSON(b []byte) error {
	o.Some = !_JsonNullPattern.Match(b)
	if !o.Some {
		return nil
	}

	return json.Unmarshal(b, &o.Val)
}
