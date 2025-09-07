package user

import (
	"errors"

	"github.com/alan-b-lima/prp/pkg/uuid"
)

type User struct {
	uuid     uuid.UUID
	name     string
	login    string
}

type UserScrath struct {
	Name  string
	Login string
}

var ErrLoginEmptyString = errors.New("user: login cannot be empty")

func NewUser(us *UserScrath) (*User, error) {
	var u User

	errs := [...]error{
		u.SetName(us.Name),
		u.SetLogin(us.Login),
	}

	if err := errors.Join(errs[:]...); err != nil {
		return nil, err
	}

	u.uuid = uuid.NewUUIDv7()
	return &u, nil
}

func (u *User) UUID() uuid.UUID {
	return u.uuid
}

func (u *User) Name() string {
	return u.name
}

func (u *User) SetName(name string) error {
	u.name = name
	return nil
}

func (u *User) Login() string {
	return u.login
}

func (u *User) SetLogin(login string) error {
	if login == "" {
		return ErrLoginEmptyString
	}

	u.login = login
	return nil
}
