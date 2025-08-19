// Copyright (C) 2025 Alan Barbosa Lima
//
// PRP is free software: you can redistribute it and/or modify it
// under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// PRP is distributed in the hope that it will be useful, but WITHOUT
// ANY WARRANTY; without even the implied warranty of MERCHANTABILITY
// or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public
// License for more details.
//
// You should have received a copy of the GNU General Public License
// along with PRP, located in LICENSE, at the root of the source
// tree. If not, see <https://www.gnu.org/licenses/>.

package user

import "github.com/alan-b-lima/prp/internal/pkg/uuid"

type UserUUID = uuid.UUID

type User struct {
	uuid UserUUID
	name string
}

func NewUser(name string) *User {
	return &User{
		uuid: uuid.NewUUIDv7(),
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
