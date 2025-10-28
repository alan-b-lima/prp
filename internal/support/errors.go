package router

import (
	"github.com/alan-b-lima/prp/pkg/errors"
)

var (
	ErrBadUUID = errors.Imp(errors.InvalidInput, "bad-uuid", "given UUID could not be parsed")

	ErrBadOffsetOrLimit = errors.Imp(errors.InvalidInput, "bad-offset-or-limit", "bad offset or limit params")

	ErrNoContentType              = errors.New(errors.PreconditionFailed, "no-content-type", "content type must be informed", nil)
	ErrUnsupportedContentTypeJson = errors.New(errors.PreconditionFailed, "unsupported-content-type", "content type must be application/json", nil)

	ErrJsonSyntax        = errors.Gen(errors.InvalidInput, "json-syntax-error")
	ErrNotAcceptableJson = errors.New(errors.PreconditionFailed, "not-acceptable-type", "client does not accept application/json", nil)
)
