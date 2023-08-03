package app

import (
	"testing"
	"time"

	"github.com/bad33ndj3/bunqtoynab/pkg/bunq"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestTransformBunqToYNABPayload(t *testing.T) {
	t.Run("when Payee is longer than maxPayeeLength", func(t *testing.T) {
		txn := &bunq.Transaction{
			ID:        1,
			Payee:     "0123456789",
			Amount:    decimal.NewFromInt(10),
			Date:      time.Now(),
			AccountID: 1,
		}
		actual := TransformBunqToYNABPayload(txn, "testAccountID")
		assert.Equal(t, "012345: ", *actual.Memo)
	})

	t.Run("when Payee is shorter than maxPayeeLength", func(t *testing.T) {
		txn := &bunq.Transaction{
			ID:        1,
			Payee:     "0123",
			Amount:    decimal.NewFromInt(10),
			Date:      time.Now(),
			AccountID: 1,
		}
		actual := TransformBunqToYNABPayload(txn, "testAccountID")
		assert.Equal(t, "0123: ", *actual.Memo)
	})
}

func TestImportID(t *testing.T) {
	// Mock Transaction object
	mockTransaction := bunq.Transaction{
		Amount: decimal.NewFromInt(1000),
		Date:   time.Date(2006, 1, 2, 15, 4, 5, 0, time.UTC),
	}

	want := "YNAB:1000:2006-01-02:1"

	got := importID(&mockTransaction)

	if got != want {
		t.Errorf("importID() = %q, want %q", got, want)
	}
}
