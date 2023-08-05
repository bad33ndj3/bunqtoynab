// Package bunq provides a client for the bunq API.
package bunq

import (
	"context"
	"time"

	"github.com/OGKevin/go-bunq/bunq"
	"github.com/bad33ndj3/bunqtoynab/pkg/filter"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"go.uber.org/ratelimit"
)

const (
	name   = "cli"
	layout = "2006-01-02 15:04:05.000000"
)

// Client is a client for the bunq API.
type Client struct {
	client *bunq.Client
	rt     ratelimit.Limiter
}

// NewClient creates a new Client.
func NewClient(ctx context.Context, apiKey string, rt ratelimit.Limiter) (*Client, error) {
	key, err := bunq.CreateNewKeyPair()
	if err != nil {
		return nil, errors.Wrap(err, "creating new key pair")
	}

	bunqClient := bunq.NewClient(ctx, bunq.BaseURLProduction, key, apiKey, name)

	err = bunqClient.Init()
	if err != nil {
		return nil, errors.Wrap(err, "initializing bunq client")
	}

	return &Client{
		client: bunqClient,
		rt:     rt,
	}, nil
}

// AllPayments returns all payments for the given account.
func (c *Client) AllPayments(
	accountID uint,
	filters ...filter.Func[*Transaction],
) ([]*Transaction, error) {
	var transactions []*Transaction

	c.rt.Take()

	allPaymentResponse, err := c.client.PaymentService.GetAllPayment(accountID)
	if err != nil {
		return nil, errors.Wrap(err, "getting all payments")
	}

	for _, r := range allPaymentResponse.Response {
		payment := r.Payment

		amount, err := decimal.NewFromString(payment.Amount.Value)
		if err != nil {
			return nil, errors.Wrap(err, "converting amount to decimal")
		}

		date, err := time.Parse(layout, payment.Created)
		if err != nil {
			return nil, errors.Wrap(err, "parsing date")
		}

		transaction := &Transaction{
			ID:          payment.ID,
			Description: payment.Description,
			Amount:      amount,
			AccountID:   accountID,
			Date:        date,
			Type:        PaymentTypeFromString(payment.Type),
			SubType:     PaymentSubTypeFromString(payment.SubType),
			Payee:       payment.CounterpartyAlias.DisplayName,
			PayeeIBAN:   payment.CounterpartyAlias.IBAN,
		}

		for _, f := range filters {
			if f(transaction) {
				transactions = append(transactions, transaction)
			}
		}
	}

	return transactions, nil
}

// AccountsIBAN returns the IBANs of the given accounts.
func AccountsIBAN(accounts []*Account) []string {
	var ibans []string
	for _, account := range accounts {
		ibans = append(ibans, account.IBAN)
	}

	return ibans
}

// AllAccounts returns all accounts.
func (c *Client) AllAccounts() ([]*Account, error) {
	var accounts []*Account
	c.rt.Take()

	savingAccounts, err := c.client.AccountService.GetAllMonetaryAccountSaving()
	if err != nil {
		return nil, errors.Wrap(err, "getting all saving accounts")
	}

	for _, r := range savingAccounts.Response {
		acc := r.MonetaryAccountSaving
		if acc.Status != AccountStatusActive.String() {
			continue
		}

		account := &Account{
			ID:          acc.ID,
			Description: acc.Description,
			AccountType: AccountTypeSaving,
		}
		if len(acc.Alias) > 0 {
			account.IBAN = acc.Alias[0].Value
		}
		accounts = append(accounts, account)
	}

	c.rt.Take()
	bankAccounts, err := c.client.AccountService.GetAllMonetaryAccountBank()
	if err != nil {
		return nil, errors.Wrap(err, "getting all bank accounts")
	}

	for _, r := range bankAccounts.Response {
		acc := r.MonetaryAccountBank
		if acc.Status != AccountStatusActive.String() {
			continue
		}
		account := &Account{
			ID:          acc.ID,
			Description: acc.Description,
			AccountType: AccountTypeBank,
		}
		if len(acc.Alias) > 0 {
			account.IBAN = acc.Alias[0].Value
		}
		accounts = append(accounts, account)
	}

	return accounts, nil
}

// Account returns the account with the given IBAN.
type Account struct {
	ID          int
	Description string
	AccountType AccountType
	IBAN        string
}

// Transaction represents a transaction.
type Transaction struct {
	ID          int
	Description string
	Amount      decimal.Decimal
	AccountID   uint
	Date        time.Time
	Payee       string
	Type        PaymentType
	SubType     PaymentSubType
	PayeeIBAN   string
}
