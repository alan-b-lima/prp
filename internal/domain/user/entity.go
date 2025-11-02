package user

import (
	"unicode/utf8"

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
	level    int
}

type Scratch struct {
	Name     string
	Login    string
	Password string
	Level    int
}

func New(us *Scratch) (User, error) {
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
		return User{}, xerrors.ErrUserCreation.New(err)
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
		return xerrors.ErrNameEmpty
	}

	u.name = name
	return nil
}

func (u *User) Login() string {
	return u.login
}

func (u *User) SetLogin(login string) error {
	if login == "" {
		return xerrors.ErrLoginNameEmpty
	}

	u.login = login
	return nil
}

func (u *User) Password() [60]byte {
	return u.password
}

func (u *User) SetPassword(password string) error {
	if len(password) < 8 {
		return xerrors.ErrPasswordTooShort
	}

	if len(password) > 64 {
		return xerrors.ErrPasswordTooLong
	}

	switch password[0] {
	case ' ', '\t', '\n', '\r':
		return xerrors.ErrPasswordLeadOrTrailWhitespace
	}

	switch password[len(password)-1] {
	case ' ', '\t', '\n', '\r':
		return xerrors.ErrPasswordLeadOrTrailWhitespace
	}

	for _, rune := range password {
		if rune < ' ' || !utf8.ValidRune(rune) {
			return xerrors.ErrPasswordIllegalCharacters
		}
	}

	hash, err := hash.Hash([]byte(password))
	if err != nil {
		return xerrors.ErrFailedToHashPassword.New(err)
	}

	u.password = hash
	return nil
}

var AuthLevel = struct {
	Admin int
	User  int
}{
	Admin: 0,
	User:  1,
}
