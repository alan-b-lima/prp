// Copyright 2025 Alan Barbosa Lima
// Licensed under the Apache License, Version 2.0

package user

import "github.com/alan-b-lima/prp/server/internal/uuid"

type UserUUID uuid.UUID

type User struct {
	uuid UserUUID
	name string
}

func NewUser(name string) *User {
	return &User{
		uuid: UserUUID(uuid.NewUUIDv7()),
		name: name,
	}
}

func (u *User) UUID() UserUUID {
	return u.uuid
}

func (u *User) Name() string {
	return u.name
}

func (u *User) SetName(name string) error {
	u.name = name
	return nil
}
