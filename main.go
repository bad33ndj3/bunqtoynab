package main

import (
	"context"
	"log"
	"os"

	"github.com/bad33ndj3/bunqtoynab/pkg/bunq"
	"github.com/brunomvsouza/ynab.go"
	"github.com/brunomvsouza/ynab.go/api"
	"github.com/brunomvsouza/ynab.go/api/transaction"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	ynabKey := os.Getenv("YNAB_KEY")
	ynabBudgetID := os.Getenv("YNAB_BUDGET_ID")
	ynabAccountName := os.Getenv("YNAB_ACCOUNT_NAME")
	bunqKey := os.Getenv("BUNQ_KEY")

	ynabClient := ynab.NewClient(ynabKey)

	bunqClient, err := bunq.NewClient(context.Background(), bunqKey)
	if err != nil {
		log.Fatalf("error creating bunq client: %v", err)
	}

	app := App{
		bunqClient: bunqClient,
		ynabClient: ynabClient,
	}

	err = app.ImportAsOne(ynabBudgetID, ynabAccountName)
	if err != nil {
		log.Fatalf("error running program: %v", err)
	}
}

type App struct {
	bunqClient *bunq.Client
	ynabClient ynab.ClientServicer
}

// ImportAsOne imports all transactions from bunq to ynab as if they were all from the same account.
// This is useful if you have multiple bunq accounts and want to import them all as one ynab account.
// To make this less confusing, all internal bunq transactions are ignored.
func (s App) ImportAsOne(budgetID string, accountName string) error {
	accounts, err := s.bunqClient.AllAccounts()
	if err != nil {
		return errors.Wrap(err, "getting all accounts")
	}

	withoutInternalTransactions := bunq.InverseFilterFunc(
		bunq.WithPayeeIBAN(bunq.AccountsIBAN(accounts)...),
	)

	var transactions []*bunq.Transaction
	for _, acc := range accounts {
		trans, err := s.bunqClient.AllPayments(uint(acc.ID), withoutInternalTransactions)
		if err != nil {
			return errors.Wrap(err, "getting all payments")
		}

		transactions = append(transactions, trans...)
	}

	ynabAccountsResp, err := s.ynabClient.Account().GetAccounts(budgetID, nil)
	if err != nil {
		return errors.Wrap(err, "getting ynab accounts")
	}

	var accountID string
	for _, acc := range ynabAccountsResp.Accounts {
		if acc.Name == accountName {
			accountID = acc.ID
		}
	}

	ynabTransactions := make([]transaction.PayloadTransaction, len(transactions))
	for i, t := range transactions {
		ynabTransactions[i] = TransformBunqToYNABPayload(t, accountID)
	}

	_, err = s.ynabClient.Transaction().CreateTransactions(budgetID, ynabTransactions)
	if err != nil {
		return errors.Wrap(err, "creating transactions")
	}

	return nil
}

func TransformBunqToYNABPayload(
	t *bunq.Transaction,
	accountID string,
) transaction.PayloadTransaction {
	importID := importID(t)

	const maxPayeeLenght = 6

	var shortPayee string
	if len(t.Payee) > maxPayeeLenght {
		shortPayee = t.Payee[:maxPayeeLenght]
	} else {
		shortPayee = t.Payee
	}

	description := shortPayee + ": " + t.Description

	return transaction.PayloadTransaction{
		ID:         "",
		AccountID:  accountID,
		Date:       api.Date{Time: t.Date},
		Amount:     t.Amount.Mul(decimal.NewFromInt(1000)).IntPart(),
		Memo:       &description,
		Cleared:    transaction.ClearingStatusUncleared,
		Approved:   false,
		PayeeID:    nil,
		PayeeName:  &t.Payee,
		CategoryID: nil,
		FlagColor:  nil,
		ImportID:   &importID,
	}
}

// importID generates an importID for a transaction.
// This is used by YNAB to prevent duplicate imports.
// If you want to import the same transaction multiple times, you can change the importIteration.
func importID(t *bunq.Transaction) string {
	const importIteration = "1"

	return "YNAB:" + t.Amount.String() + ":" + t.Date.Format("2006-01-02") + ":" + importIteration
}
