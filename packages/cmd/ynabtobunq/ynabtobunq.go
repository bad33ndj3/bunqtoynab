package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/bad33ndj3/bunqtoynab/pkg/app"
	"github.com/bad33ndj3/bunqtoynab/pkg/bunq"
	"github.com/brunomvsouza/ynab.go"
	"go.uber.org/ratelimit"
)

func Main(_ map[string]interface{}) map[string]interface{} {
	ynabKey := os.Getenv("YNAB_KEY")
	ynabBudgetID := os.Getenv("YNAB_BUDGET_ID")
	ynabAccountName := os.Getenv("YNAB_ACCOUNT_NAME")
	bunqKey := os.Getenv("BUNQ_KEY")

	ynabClient := ynab.NewClient(ynabKey)

	rt := ratelimit.New(3, ratelimit.Per(time.Second*3))
	bunqClient, err := bunq.NewClient(context.Background(), bunqKey, rt)
	if err != nil {
		log.Panicf("error creating bunq client: %v", err)
	}

	a := app.NewApp(bunqClient, ynabClient)
	err = a.ImportAsOne(ynabBudgetID, ynabAccountName)
	if err != nil {
		log.Panicf("error running program: %v", err)
	}

	log.Println("done")
	return map[string]interface{}{}
}
