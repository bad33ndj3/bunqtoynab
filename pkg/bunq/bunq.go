package bunq

import (
	"context"
	"github.com/OGKevin/go-bunq/bunq"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"go.uber.org/ratelimit"
	"time"
)

const name = "cli"
const layout = "2006-01-02 15:04:05.000000"

type Client struct {
	apiKey      string
	description string

	bunqClient *bunq.Client
	rt         ratelimit.Limiter
}

func NewClient(apiKey string) (*Client, error) {
	c := &Client{
		apiKey:      apiKey,
		description: name,
	}
	err := c.init(context.Background())
	if err != nil {
		return nil, errors.Wrap(err, "initializing client")
	}
	return c, nil
}

func (c *Client) init(ctx context.Context) error {
	key, err := bunq.CreateNewKeyPair()
	if err != nil {
		return errors.Wrap(err, "creating new key pair")
	}

	c.bunqClient = bunq.NewClient(ctx, bunq.BaseURLProduction, key, c.apiKey, c.description)
	err = c.bunqClient.Init()
	if err != nil {
		return errors.Wrap(err, "initializing bunq client")
	}

	c.rt = ratelimit.New(1)

	return nil
}

func (c *Client) AllPayments(accountID uint) ([]*Transaction, error) {
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

		transactions = append(transactions, &Transaction{
			ID:          payment.ID,
			Description: payment.Description,
			Amount:      amount,
			AccountID:   accountID,
			Date:        date,
			Payee:       payment.MerchantReference,
		})
	}

	return transactions, nil
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
		accounts = append(accounts, &Account{
			ID:          acc.ID,
			Description: acc.Description,
			AccountType: AccountTypeSaving,
		})

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
		accounts = append(accounts, &Account{
			ID:          acc.ID,
			Description: acc.Description,
			AccountType: AccountTypeBank,
		})
	}

	return accounts, nil
}

type AccountType string

const (
	AccountTypeBank   AccountType = "BANK"
	AccountTypeSaving AccountType = "SAVING"
)

type AccountStatus string

func (s AccountStatus) String() string {
	return string(s)
}

const (
	AccountStatusActive AccountStatus = "ACTIVE"
)

type Account struct {
	ID          int
	Description string
	AccountType AccountType
}

type Transaction struct {
	ID          int
	Description string
	Amount      decimal.Decimal
	AccountID   uint
	Date        time.Time
	Payee       string
}
