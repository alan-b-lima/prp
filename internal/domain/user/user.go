package user

import (
	"errors"
	"regexp"

	"github.com/alan-b-lima/prp/pkg/hash"
	"github.com/alan-b-lima/prp/pkg/uuid"
)

type User struct {
	uuid     uuid.UUID
	name     string
	login    string
	password [32]byte
	salt     [16]byte
}

type UserScrath struct {
	Name     string
	Login    string
	Password string
}

var (
	ErrNameEmpty                       = errors.New("user(entity): name cannot be empty")
	ErrLoginEmpty                      = errors.New("user(entity): login cannot be empty")
	ErrPasswordEmpty                   = errors.New("user(entity): password cannot be empty")
	ErrPasswordTooShort                = errors.New("user(entity): password must be at least 8 characters long")
	ErrPasswordTrailingOrLeadingSpaces = errors.New("user(entity): password must not begin or end with a space")
	ErrPasswordMalformed               = errors.New("user(entity): password contains unallowed characters")
)

func NewUser(us *UserScrath) (*User, error) {
	u := new(User)

	err := errors.Join(
		u.SetName(us.Name),
		u.SetLogin(us.Login),
		u.SetPassword(us.Password),
	)
	if err != nil {
		return nil, err
	}

	u.uuid = uuid.NewUUIDv7()
	return u, nil
}

func (u *User) UUID() uuid.UUID {
	return u.uuid
}

func (u *User) Name() string {
	return u.name
}

func (u *User) SetName(name string) error {
	if name == "" {
		return ErrNameEmpty
	}

	u.name = name
	return nil
}

func (u *User) Login() string {
	return u.login
}

func (u *User) SetLogin(login string) error {
	if login == "" {
		return ErrLoginEmpty
	}

	u.login = login
	return nil
}

var PasswordWellformedPattern = regexp.MustCompile(`^[^\x00-\x1F]*$`)

func (u *User) SetPassword(password string) error {
	if password == "" {
		return ErrPasswordEmpty
	}

	if len(password) < 8 {
		return ErrPasswordTooShort
	}

	if password[0] == ' ' || password[len(password)-1] == ' ' {
		return ErrPasswordTrailingOrLeadingSpaces
	}

	if !PasswordWellformedPattern.MatchString(password) {
		return ErrPasswordMalformed
	}

	u.salt = hash.NewSalt()
	u.password = hash.Hash(append([]byte(password), u.salt[:]...))
	return nil
}
