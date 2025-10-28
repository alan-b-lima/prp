package user

import "github.com/alan-b-lima/prp/pkg/errors"

// entity
var (
	ErrNameEmpty                     = errors.New(errors.InvalidInput, "name-empty", "name cannot be empty", nil)
	ErrLoginNameEmpty                = errors.New(errors.InvalidInput, "login-empty", "login cannot be empty", nil)
	ErrPasswordTooShort              = errors.New(errors.InvalidInput, "password-short", "password must be at least 8 characters long", nil)
	ErrPasswordTooLong               = errors.New(errors.InvalidInput, "password-long", "password must be a maximum of 64 characters long", nil)
	ErrPasswordLeadOrTrailWhitespace = errors.New(errors.InvalidInput, "password-edge-whitespace", "password must not begin or end with whitespaces", nil)
	ErrPasswordIllegalCharacters     = errors.New(errors.InvalidInput, "password-illegal-chars", "password must not contain unprintable or invalid uft-8 characters", nil)

	ErrFailedToHashPassword = errors.Imp(errors.Internal, "hash-failure", "failed to hash the password")
)

// repository
var (
	ErrUserNotFound = errors.New(errors.NotFound, "user-not-found", "user not found", nil)
	ErrLoginTaken   = errors.New(errors.Conflict, "login-in-use", "login already taken", nil)
)
