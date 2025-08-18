// Copyright 2025 Alan Barbosa Lima
// Licensed under the Apache License, Version 2.0

package finance

import (
	"github.com/alan-b-lima/prp/server/internal/uuid"
	"github.com/alan-b-lima/prp/server/pkg/model/user"
)

type AccountUUID uuid.UUID

type Account struct {
	uuid         AccountUUID
	name         string
	account_type AccountType
	user         user.UserUUID
	parent       AccountUUID
}

func NewAccount(name string, account_type AccountType, user user.UserUUID, parent AccountUUID) *Account {
	return &Account{
		uuid:         AccountUUID(uuid.NewUUIDv7()),
		name:         name,
		account_type: account_type,
		user:         user,
		parent:       parent,
	}
}

func (a *Account) UUID() AccountUUID {
	return a.uuid
}

func (a *Account) Name() string {
	return a.name
}

func (a *Account) SetName(name string) error {
	a.name = name
	return nil
}

func (a *Account) AccountType() AccountType {
	return a.account_type
}

func (a *Account) SetAccountType(accountType AccountType) error {
	a.account_type = accountType
	return nil
}

func (a *Account) User() user.UserUUID {
	return a.user
}

func (a *Account) SetUser(user user.UserUUID) error {
	a.user = user
	return nil
}

func (a *Account) Parent() AccountUUID {
	return a.parent
}

func (a *Account) SetParent(parent AccountUUID) error {
	a.parent = parent
	return nil
}

type AccountType int

const (
	Asset AccountType = iota
	Liability
	Equity
	Revenue
	Expense
)

var accountTypeString = map[AccountType]string{
	Asset:     "ASSET",
	Liability: "LIABILITY",
	Equity:    "EQUITY",
	Revenue:   "REVENUE",
	Expense:   "EXPENSE",
}

func (a AccountType) String() string {
	return accountTypeString[a]
}
