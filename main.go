package main

import (
	"context"
	"log"
	"os"

	"github.com/bad33ndj3/bunqtoynab/pkg/app"
	"github.com/bad33ndj3/bunqtoynab/pkg/bunq"
	"github.com/brunomvsouza/ynab.go"
	"github.com/joho/godotenv"
	"go.uber.org/ratelimit"
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

	bunqClient, err := bunq.NewClient(context.Background(), bunqKey, ratelimit.New(1))
	if err != nil {
		log.Fatalf("error creating bunq client: %v", err)
	}

	a := app.NewApp(bunqClient, ynabClient)
	err = a.ImportAsOne(ynabBudgetID, ynabAccountName)
	if err != nil {
		log.Fatalf("error running program: %v", err)
	}
}
