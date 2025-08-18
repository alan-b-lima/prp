// Copyright 2025 Alan Barbosa Lima
// Licensed under the Apache License, Version 2.0

package finance

import (
	"time"

	"github.com/alan-b-lima/prp/server/internal/uuid"
	"github.com/alan-b-lima/prp/server/pkg/model/user"
)

type TransactionUUID uuid.UUID

type Transaction struct {
	uuid        TransactionUUID
	user        user.UserUUID
	description string
	amount      float64

	debit_account  AccountUUID
	credit_account AccountUUID

	transaction_date time.Time
	settlement_date  time.Time
}

func NewTransaction(user user.UserUUID, description string, amount float64, debit_account, credit_account AccountUUID, transaction_date, settlement_date time.Time) *Transaction {
	return &Transaction{
		uuid:             TransactionUUID(uuid.NewUUIDv7()),
		user:             user,
		description:      description,
		amount:           amount,
		debit_account:    debit_account,
		credit_account:   credit_account,
		transaction_date: transaction_date,
		settlement_date:  settlement_date,
	}
}

func (t *Transaction) UUID() TransactionUUID {
	return t.uuid
}

func (t *Transaction) User() user.UserUUID {
	return t.user
}

func (t *Transaction) SetUser(user user.UserUUID) error {
	t.user = user
	return nil
}

func (t *Transaction) Description() string {
	return t.description
}

func (t *Transaction) SetDescription(description string) error {
	t.description = description
	return nil
}

func (t *Transaction) Amount() float64 {
	return t.amount
}

func (t *Transaction) SetAmount(amount float64) error {
	t.amount = amount
	return nil
}

func (t *Transaction) DebitAccount() AccountUUID {
	return t.debit_account
}

func (t *Transaction) SetDebitAccount(debit_account AccountUUID) error {
	t.debit_account = debit_account
	return nil
}

func (t *Transaction) CreditAccount() AccountUUID {
	return t.credit_account
}

func (t *Transaction) SetCreditAccount(credit_account AccountUUID) error {
	t.credit_account = credit_account
	return nil
}

func (t *Transaction) TransactionDate() time.Time {
	return t.transaction_date
}

func (t *Transaction) SetTransactionDate(transaction_date time.Time) error {
	t.transaction_date = transaction_date
	return nil
}

func (t *Transaction) SettlementDate() time.Time {
	return t.settlement_date
}

func (t *Transaction) SetSettlementDate(settlement_date time.Time) error {
	t.settlement_date = settlement_date
	return nil
}
