package errors

import "fmt"

type gen struct {
	kind  Kind
	title string
}

func Gen(kind Kind, title string) *gen {
	return &gen{
		kind:  kind,
		title: title,
	}
}

func (gen *gen) New(message string, cause error) error {
	return &Error{
		Kind:    gen.kind,
		Title:   gen.title,
		Message: message,
		Cause:   cause,
	}
}

type imp struct {
	kind    Kind
	title   string
	message string
}

func Imp(kind Kind, title, message string) *imp {
	return &imp{
		kind:    kind,
		title:   title,
		message: message,
	}
}

func (gen *imp) New(cause error) error {
	return &Error{
		Kind:    gen.kind,
		Title:   gen.title,
		Message: gen.message,
		Cause:   cause,
	}
}

type fmt_ struct {
	kind   Kind
	title  string
	format string
}

func Fmt(kind Kind, title, format string) *fmt_ {
	return &fmt_{
		kind:   kind,
		title:  title,
		format: format,
	}
}

func (gen *fmt_) New(v ...any) error {
	err := fmt.Errorf(gen.format, v...)

	return &Error{
		Kind:    gen.kind,
		Title:   gen.title,
		Message: err.Error(),
		Cause:   nil,
	}
}
