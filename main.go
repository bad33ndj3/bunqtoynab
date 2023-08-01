package main

import (
	"github.com/bad33ndj3/bunqtoynab/pkg/bunq"
	"github.com/brunomvsouza/ynab.go"
	"github.com/brunomvsouza/ynab.go/api"
	"github.com/brunomvsouza/ynab.go/api/transaction"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"gopkg.in/ffmt.v1"
	"log"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	ynabKey := os.Getenv("YNAB_KEY")
	ynabBudgetID := os.Getenv("YNAB_BUDGET_ID")
	bunqKey := os.Getenv("BUNQ_KEY")

	ynabClient := ynab.NewClient(ynabKey)
	bunqClient, err := bunq.NewClient(bunqKey)
	if err != nil {
		log.Fatalf("error creating bunq client: %v", err)
	}

	app := App{
		bunqClient: bunqClient,
		ynabClient: ynabClient,
	}

	err = app.Run(ynabBudgetID)
	if err != nil {
		log.Fatalf("error running program: %v", err)
	}
}

type App struct {
	bunqClient *bunq.Client
	ynabClient ynab.ClientServicer
}

func (s App) Run(budgetID string) error {
	accounts, err := s.bunqClient.AllAccounts()
	if err != nil {
		return errors.Wrap(err, "getting all accounts")
	}

	var transactions []*bunq.Transaction
	for _, acc := range accounts {
		trans, err := s.bunqClient.AllPayments(uint(acc.ID))
		if err != nil {
			return errors.Wrap(err, "getting all payments")
		}
		transactions = append(transactions, trans...)
	}

	ynabTransactions := make([]transaction.PayloadTransaction, len(transactions))
	for i, t := range transactions {
		ynabTransactions[i] = TransformBunqToYNABPayload(t)
	}

	// todo: this implementation would add all existing transactions to one account
	ffmt.Pjson(ynabTransactions)

	//_, err = s.ynabClient.Transaction().CreateTransactions(budgetID, ynabTransactions)
	//if err != nil {
	//	return errors.Wrap(err, "creating transactions")
	//}

	return nil
}

func TransformBunqToYNABPayload(t *bunq.Transaction) transaction.PayloadTransaction {
	importID := importID(t)
	return transaction.PayloadTransaction{
		AccountID: "",
		Date:      api.Date{Time: t.Date},
		Amount:    t.Amount.Mul(decimal.NewFromInt(1000)).IntPart(),
		PayeeName: &t.Payee,
		Memo:      &t.Description,
		FlagColor: nil,
		ImportID:  &importID,
	}

}

// YNAB:[milliunit_amount]:[iso_date]:[occurrence]'. For example,
// a transaction dated 2015-12-30 in the amount of -$294.23 USD
// would have an import_id of 'YNAB:-294230:2015-12-30:1'.
// If a second transaction on the same account was imported and
// had the same date and same amount,
// its import_id would be 'YNAB:-294230:2015-12-30:2'.
func importID(t *bunq.Transaction) string {
	return "YNAB:" + t.Amount.String() + ":" + t.Date.Format("2006-01-02") + ":1"
}
