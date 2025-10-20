// Copyright (C) 2025 Alan Barbosa Lima.
//
// PRP is licensed under the GNU General Public License
// version 3. You should have received a copy of the
// license, located in LICENSE, at the root of the source
// tree. If not, see <https://www.gnu.org/licenses/>.

// Package errors implements functionally for internal standardization of error
// handling.
package errors

import (
	"strings"
)

// Error is a [error] implementer that caries information about its package and
// origin (may also be referred as a escope), with the goal of providing more
// standardized internal error messages.
type Error struct {
	Message string
	Package string
	Origin  string
	Causes  []error
}

const _Indent = "\t"

// New generates a new error given all its parameters. If no causes are
// present, it will always generate a non-nil error. If one or more causes are
// present, then a non-nil error will be generated if at least one of the
// causes is non-nil.
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
// If package is empty, "<Package>(<Origin>):" is omitted, regardless of the
// emptyness of origin. If origin is empty, "(<Origin>)" is ommited. If causes
// is empty, the cause clause, as well as any newlines inserted by it, are
// ommited.
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

// Unwrap returns all the causes of the error, all causes are non-nil errors.
// The returned slice may be empty.
func (e Error) Unwrap() []error {
	return e.Causes
}

// ErrorGen is a factory of [Error], this avoids repetition an possible
// mistakes when a package/origin has multiple errors.
type ErrorGen struct {
	Package string
	Origin  string
}

// NewGen creates a new [ErrorGen].
func NewGen(pkg, origin string) ErrorGen {
	return ErrorGen{Package: pkg, Origin: origin}
}

// New is a wrapper over [New], the package level function, that instead of
// receiving a package and an origin, it uses its [ErrorGen] fields.
func (eg *ErrorGen) New(msg string, cause ...error) error {
	return New(msg, eg.Package, eg.Origin, cause...)
}
