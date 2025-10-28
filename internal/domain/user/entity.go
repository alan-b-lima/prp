package user

import (
	"unicode/utf8"

	"github.com/alan-b-lima/prp/pkg/errors"
	"github.com/alan-b-lima/prp/pkg/hash"
	"github.com/alan-b-lima/prp/pkg/uuid"
)

type User struct {
	uuid     uuid.UUID
	name     string
	login    string
	password [60]byte
}

type Scratch struct {
	Name     string
	Login    string
	Password string
}

func NewUser(us *Scratch) (User, error) {
	user := User{}

	pwderr := user.SetPassword(us.Password)
	if err, ok := errors.AsType[*errors.Error](pwderr); ok && !err.Kind.IsClient() {
		return User{}, err
	}

	err := errors.Join(
		user.SetName(us.Name),
		user.SetLogin(us.Login),
		pwderr,
	)
	if err != nil {
		return User{}, errors.New(
			errors.InvalidInput, "",
			"given data does not satisfy the user type",
			err,
		)
	}

	user.uuid = uuid.NewUUIDv7()
	return user, nil
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
		return ErrLoginNameEmpty
	}

	u.login = login
	return nil
}

func (u *User) SetPassword(password string) error {
	if len(password) < 8 {
		return ErrPasswordTooShort
	}

	if len(password) > 64 {
		return ErrPasswordTooLong
	}

	switch password[0] {
	case ' ', '\t', '\n', '\r':
		return ErrPasswordLeadOrTrailWhitespace
	}

	switch password[len(password)-1] {
	case ' ', '\t', '\n', '\r':
		return ErrPasswordLeadOrTrailWhitespace
	}

	for _, rune := range password {
		if rune < ' ' || !utf8.ValidRune(rune) {
			return ErrPasswordIllegalCharacters
		}
	}

	hash, err := hash.Hash([]byte(password))
	if err != nil {
		return ErrFailedToHashPassword.New(err)
	}

	u.password = hash
	return nil
}

func (u *User) ComparePassword(password string) bool {
	return hash.Compare(u.password[:], []byte(password))
}
