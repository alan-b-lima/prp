package xerrors

import "github.com/alan-b-lima/prp/pkg/errors"

var (
	ErrBadUUID = errors.Imp(errors.InvalidInput, "bad-uuid", "given UUID could not be parsed")

	ErrBadOffsetOrLimit = errors.Imp(errors.InvalidInput, "bad-offset-or-limit", "bad offset or limit params")

	ErrNoContentType              = errors.New(errors.PreconditionFailed, "no-content-type", "content type must be informed", nil)
	ErrUnsupportedContentTypeJson = errors.New(errors.PreconditionFailed, "unsupported-content-type", "content type must be application/json", nil)

	ErrJsonSyntax        = errors.Gen(errors.InvalidInput, "json-syntax-error")
	ErrNotAcceptableJson = errors.New(errors.PreconditionFailed, "not-acceptable-type", "client does not accept application/json", nil)

	ErrSessionNotFound = errors.New(errors.NotFound, "session-not-found", "session not found", nil)

	ErrUserCreation = errors.Imp(errors.InvalidInput, "user-creation", "given data does not satisfy the user type")

	ErrNameEmpty                     = errors.New(errors.InvalidInput, "name-empty", "name cannot be empty", nil)
	ErrLoginNameEmpty                = errors.New(errors.InvalidInput, "login-empty", "login cannot be empty", nil)
	ErrPasswordTooShort              = errors.New(errors.InvalidInput, "password-short", "password must be at least 8 characters long", nil)
	ErrPasswordTooLong               = errors.New(errors.InvalidInput, "password-long", "password must be a maximum of 64 characters long", nil)
	ErrPasswordLeadOrTrailWhitespace = errors.New(errors.InvalidInput, "password-edge-whitespace", "password must not begin or end with whitespaces", nil)
	ErrPasswordIllegalCharacters     = errors.New(errors.InvalidInput, "password-illegal-chars", "password must not contain unprintable or invalid uft-8 characters", nil)

	ErrIncorrectPassword    = errors.New(errors.Unauthorized, "incorrect-password", "given password is incorrect", nil)
	ErrFailedToHashPassword = errors.Imp(errors.Internal, "hash-failure", "failed to hash the password")

	ErrUnauthenticatedUser = errors.New(errors.Unauthorized, "unauthenticated-user", "user is not logged in", nil)
	ErrUnauthorizedUser    = errors.Fmt(errors.Forbidden, "unauthorized-user", "auth level %v does not match any criteria in %v")

	ErrUserNotFound = errors.New(errors.NotFound, "user-not-found", "user not found", nil)
	ErrLoginTaken   = errors.New(errors.Conflict, "login-in-use", "login already taken", nil)
)
