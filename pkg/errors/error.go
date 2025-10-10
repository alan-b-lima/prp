package errors

import (
	"strings"
)

type Error struct {
	Message string
	Package string
	Origin  string
	Causes  []error
}

type ErrorGen struct {
	Package string
	Origin  string
}

const _Indent = "\t"

func New(msg, pkg, origin string, causes ...error) error {
	if len(causes) > 0 {
		var n int
		for _, err := range causes {
			if err != nil {
				n++
			}
		}
		if n == 0 {
			return nil
		}

		s := make([]error, 0, n)
		for _, err := range causes {
			if err != nil {
				s = append(s, err)
			}
		}

		causes = s
	}

	return Error{
		Message: msg,
		Package: pkg,
		Origin:  origin,
		Causes:  causes,
	}
}

// Error returns the error string, it makes the Error type be complient to the
// [error] interface.
//
// The error string given by this function follows the pattern:
//
//	<Package>(<Origin>): <Message>
//		<Causes>...
//
// If package is empty, <Package>(<Origin>): is omitted. If origin is empty,
// (<Origin>) is ommited. If causes is empty, the cause clause is ommited.
func (e Error) Error() string {
	var b strings.Builder

	if e.Package != "" {
		b.WriteString(e.Package)

		if e.Origin != "" {
			b.WriteByte('(')
			b.WriteString(e.Origin)
			b.WriteByte(')')
		}

		b.WriteString(": ")
	}

	b.WriteString(e.Message)

	for _, err := range e.Causes {
		s := err.Error()
		s = strings.ReplaceAll(s, "\n", "\n"+_Indent)

		b.WriteString("\n" + _Indent)
		b.WriteString(s)
	}

	return b.String()
}

func NewGen(pkg, origin string) ErrorGen {
	return ErrorGen{Package: pkg, Origin: origin}
}

func (eg *ErrorGen) New(msg string, cause ...error) error {
	return New(msg, eg.Package, eg.Origin, cause...)
}
