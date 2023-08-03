package bunq

import (
	"context"
	"time"

	"github.com/OGKevin/go-bunq/bunq"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"go.uber.org/ratelimit"
)

const (
	name   = "cli"
	layout = "2006-01-02 15:04:05.000000"
)

var ErrUnknownPaymentType = errors.New("unknown payment type")

type Client struct {
	bunqClient *bunq.Client
	rt         ratelimit.Limiter
}

func NewClient(ctx context.Context, apiKey string) (*Client, error) {
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
		bunqClient: bunqClient,
		rt:         ratelimit.New(1),
	}, nil
}

func (c *Client) AllPayments(
	accountID uint,
	filters ...TransactionFilterFunc,
) ([]*Transaction, error) {
	var transactions []*Transaction

	c.rt.Take()

	allPaymentResponse, err := c.bunqClient.PaymentService.GetAllPayment(accountID)
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

		for _, filter := range filters {
			if filter(transaction) {
				transactions = append(transactions, transaction)
			}
		}
	}

	return transactions, nil
}

func AccountsIBAN(accounts []*Account) []string {
	var ibans []string
	for _, account := range accounts {
		ibans = append(ibans, account.IBAN)
	}

	return ibans
}

func (c *Client) AllAccounts() ([]*Account, error) {
	var accounts []*Account
	c.rt.Take()

	savingAccounts, err := c.bunqClient.AccountService.GetAllMonetaryAccountSaving()
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
	bankAccounts, err := c.bunqClient.AccountService.GetAllMonetaryAccountBank()
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

type Account struct {
	ID          int
	Description string
	AccountType AccountType
	IBAN        string
}

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
