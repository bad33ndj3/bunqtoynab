package bunq

// TransactionFilterFunc is a function that filters transactions.
// If the function returns true, the transaction is included in the result.
type TransactionFilterFunc func(transaction *Transaction) bool

func InverseFilterFunc(filter TransactionFilterFunc) TransactionFilterFunc {
	return func(transaction *Transaction) bool {
		return !filter(transaction)
	}
}

// WithPaymentTypes returns a TransactionFilterFunc that filters out transactions without the given PaymentType's.
func WithPaymentTypes(types ...PaymentType) TransactionFilterFunc {
	return func(transaction *Transaction) bool {
		for _, t := range types {
			if transaction.Type == t {
				return true
			}
		}

		return false
	}
}

// WithPaymentSubTypes returns a TransactionFilterFunc that filters out transactions without the given PaymentSubType's.
func WithPaymentSubTypes(types ...PaymentSubType) TransactionFilterFunc {
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

// WithCombinedPaymentTypes returns a TransactionFilterFunc that filters out transactions
// without the given CombinedPaymentType's.
func WithCombinedPaymentTypes(types ...CombinedPaymentType) TransactionFilterFunc {
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

func WithPayeeIBAN(iban ...string) TransactionFilterFunc {
	return func(transaction *Transaction) bool {
		for _, i := range iban {
			if transaction.PayeeIBAN == i {
				return true
			}
		}
		return false
	}
}
