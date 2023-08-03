package bunq

import (
	"testing"
)

func TestPaymentType_String(t *testing.T) {
	tests := []struct {
		name   string
		pType  PaymentType
		expect string
	}{
		{"Payment Type Unknown", PaymentTypeUnknown, ""},
		{"Payment Type Payment", PaymentTypePayment, "PAYMENT"},
		{"Payment Type Mastercard", PaymentTypeMASTERCARD, "MASTERCARD"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := tt.pType.String()
			if assert != tt.expect {
				t.Errorf("Got %v, expect %v", assert, tt.expect)
			}
		})
	}
}

func TestPaymentTypeFromString(t *testing.T) {
	tests := []struct {
		name   string
		src    string
		expect PaymentType
	}{
		{"Empty source", "", PaymentTypeUnknown},
		{"Unknown source", "unknown", PaymentTypeUnknown},
		{"Valid source: PAYMENT", "PAYMENT", PaymentTypePayment},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := PaymentTypeFromString(tt.src)
			if assert != tt.expect {
				t.Errorf("Got %v, expect %v", assert, tt.expect)
			}
		})
	}
}

func TestPaymentSubTypeFromString(t *testing.T) {
	tests := []struct {
		name   string
		src    string
		expect PaymentSubType
	}{
		{"Empty source", "", PaymentSubTypeUnknown},
		{"Unknown source", "unknown", PaymentSubTypeUnknown},
		{"Valid source: PAYMENT", "PAYMENT", PaymentSubTypePayment},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := PaymentSubTypeFromString(tt.src)
			if assert != tt.expect {
				t.Errorf("Got %v, expect %v", assert, tt.expect)
			}
		})
	}
}

func TestAccountStatus_String(t *testing.T) {
	tests := []struct {
		name   string
		status AccountStatus
		expect string
	}{
		{"Account Status Active", AccountStatusActive, "ACTIVE"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := tt.status.String()
			if assert != tt.expect {
				t.Errorf("Got %v, expect %v", assert, tt.expect)
			}
		})
	}
}
