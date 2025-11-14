package user

import (
	"unicode/utf8"

	"github.com/alan-b-lima/prp/internal/auth"
	"github.com/alan-b-lima/prp/internal/xerrors"
	"github.com/alan-b-lima/prp/pkg/errors"
	"github.com/alan-b-lima/prp/pkg/hash"
	"github.com/alan-b-lima/prp/pkg/uuid"
)

type User struct {
	uuid     uuid.UUID
	name     string
	login    string
	password [60]byte
	level    auth.Level
}

func New(name, login, password string, level auth.Level) (User, error) {
	var u User

	errpwd := u.SetPassword(password)
	if err, ok := errors.AsType[*errors.Error](errpwd); ok && err.IsInternal() {
		return User{}, err
	}

	err := errors.Join(
		u.SetName(name),
		u.SetLogin(login),
		errpwd,
		u.SetLevel(level),
	)
	if err != nil {
		return User{}, xerrors.ErrUserCreation.New(err)
	}

	u.uuid = uuid.NewUUIDv7()
	return u, nil
}

func (u *User) UUID() uuid.UUID    { return u.uuid }
func (u *User) Name() string       { return u.name }
func (u *User) Login() string      { return u.login }
func (u *User) Password() [60]byte { return u.password }
func (u *User) Level() auth.Level  { return u.level }

func (u *User) SetName(name string) error         { return set(&u.name, name, ProcessName) }
func (u *User) SetLogin(login string) error       { return set(&u.login, login, ProcessLogin) }
func (u *User) SetPassword(password string) error { return set(&u.password, password, ProcessPassword) }
func (u *User) SetLevel(level auth.Level) error   { return set(&u.level, level, ProcessLevel) }

func ProcessName(name string) (string, error) {
	if name == "" {
		return "", xerrors.ErrNameEmpty
	}

	return name, nil
}

func ProcessLogin(login string) (string, error) {
	if login == "" {
		return "", xerrors.ErrLoginNameEmpty
	}

	return login, nil
}

func ProcessPassword(password string) ([60]byte, error) {
	if len(password) < 8 {
		return [60]byte{}, xerrors.ErrPasswordTooShort
	}

	if len(password) > 64 {
		return [60]byte{}, xerrors.ErrPasswordTooLong
	}

	switch password[0] {
	case ' ', '\t', '\n', '\r':
		return [60]byte{}, xerrors.ErrPasswordLeadOrTrailWhitespace
	}

	switch password[len(password)-1] {
	case ' ', '\t', '\n', '\r':
		return [60]byte{}, xerrors.ErrPasswordLeadOrTrailWhitespace
	}

	for _, rune := range password {
		if rune < ' ' || !utf8.ValidRune(rune) {
			return [60]byte{}, xerrors.ErrPasswordIllegalCharacters
		}
	}

	hash, err := hash.Hash([]byte(password))
	if err != nil {
		return [60]byte{}, xerrors.ErrFailedToHashPassword.New(err)
	}

	return hash, nil
}

func ProcessLevel(level auth.Level) (auth.Level, error) {
	if !level.IsValid() {
		return 0, nil
	}

	return level, nil
}

func set[T, R any](dst *R, src T, fn func(T) (R, error)) error {
	val, err := fn(src)
	if err != nil {
		return err
	}

	*dst = val
	return nil
}
