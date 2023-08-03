package bunq

import "github.com/bad33ndj3/bunqtoynab/pkg/filter"

// WithPaymentTypes returns a TransactionFilterFunc that filters out transactions without the given PaymentType's.
func WithPaymentTypes(types ...PaymentType) filter.Func[*Transaction] {
	return func(transaction *Transaction) bool {
		for _, t := range types {
			if transaction.Type == t {
				return true
			}
		}

		return false
	}
}

// WithPaymentSubTypes returns a filter.FilterFunc[] that filters out transactions without the given PaymentSubType's.
func WithPaymentSubTypes(types ...PaymentSubType) filter.Func[*Transaction] {
	return func(transaction *Transaction) bool {
		for _, t := range types {
			if transaction.SubType == t {
				return true
			}
		}

		return false
	}
}

// CombinedPaymentType is a combination of a PaymentType and a set of PaymentSubType's.
type CombinedPaymentType struct {
	PaymentType     PaymentType
	PaymentSubTypes []PaymentSubType
}

// WithCombinedPaymentTypes returns a filter.FilterFunc[] that filters out transactions
// without the given CombinedPaymentType's.
func WithCombinedPaymentTypes(types ...CombinedPaymentType) filter.Func[*Transaction] {
	return func(transaction *Transaction) bool {
		for _, t := range types {
			if WithPaymentTypes(t.PaymentType)(transaction) {
				for _, subType := range t.PaymentSubTypes {
					if WithPaymentSubTypes(subType)(transaction) {
						return true
					}
				}
			}
		}

		return false
	}
}

// WithPayeeIBAN returns a filter.FilterFunc[] that filters out transactions without the given IBAN's.
func WithPayeeIBAN(iban ...string) filter.Func[*Transaction] {
	return func(transaction *Transaction) bool {
		for _, i := range iban {
			if transaction.PayeeIBAN == i {
				return true
			}
		}
		return false
	}
}
