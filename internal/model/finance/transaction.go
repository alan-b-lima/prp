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

package finance

import (
	"errors"
	"time"

	"github.com/alan-b-lima/prp/internal/model/user"
	"github.com/alan-b-lima/prp/internal/pkg/uuid"
)

type TransactionUUID = uuid.UUID

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

type TransactionBuilder struct {
	User        user.UserUUID
	Description string
	Amount      float64

	DebitAccount  AccountUUID
	CreditAccount AccountUUID

	TransactionDate time.Time
	SettlementDate  time.Time
}

var (
	errNegativeAmount = errors.New("amount cannot be negative")
	errNilUserUUID    = errors.New("user UUID cannot be a Nil UUID")
)

func NewTransaction(tb *TransactionBuilder) (*Transaction, error) {
	var t Transaction

	errs := [...]error{
		t.SetUser(tb.User),
		t.SetDescription(tb.Description),
		t.SetAmount(tb.Amount),
		t.SetDebitAccount(tb.DebitAccount),
		t.SetCreditAccount(tb.CreditAccount),
		t.SetTransactionDate(tb.TransactionDate),
		t.SetSettlementDate(tb.SettlementDate),
	}

	err := errors.Join(errs[:]...)
	if err != nil {
		return nil, err
	}

	t.uuid = uuid.NewUUIDv7()
	return &t, nil
}

func (t *Transaction) UUID() TransactionUUID {
	return t.uuid
}

func (t *Transaction) User() user.UserUUID {
	return t.user
}

func (t *Transaction) SetUser(user user.UserUUID) error {
	if user.IsNil() {
		return errNilUserUUID
	}

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
	if amount < 0 {
		return errNegativeAmount
	}

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
