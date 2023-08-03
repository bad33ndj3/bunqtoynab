package bunq

type PaymentType string

func (p PaymentType) String() string {
	return string(p)
}

func PaymentTypeFromString(src string) PaymentType {
	switch src {
	case string(PaymentTypePayment):
		return PaymentTypePayment
	case string(PaymentTypeIDEAL):
		return PaymentTypeIDEAL
	case string(PaymentTypeBUNQ):
		return PaymentTypeBUNQ
	case string(PaymentTypeMASTERCARD):
		return PaymentTypeMASTERCARD
	case string(PaymentTypeSWIFT):
		return PaymentTypeSWIFT
	case string(PaymentTypeSAVINGS):
		return PaymentTypeSAVINGS
	case string(PaymentTypePAYDAY):
		return PaymentTypePAYDAY
	case string(PaymentTypeINTEREST):
		return PaymentTypeINTEREST
	default:
		return PaymentTypeUnknown
	}
}

const (
	PaymentTypeUnknown    PaymentType = ""
	PaymentTypePayment    PaymentType = "PAYMENT"
	PaymentTypeIDEAL      PaymentType = "IDEAL"
	PaymentTypeBUNQ       PaymentType = "BUNQ"
	PaymentTypeMASTERCARD PaymentType = "MASTERCARD"
	PaymentTypeSWIFT      PaymentType = "SWIFT"
	PaymentTypeSAVINGS    PaymentType = "SAVINGS"
	PaymentTypePAYDAY     PaymentType = "PAYDAY"
	PaymentTypeINTEREST   PaymentType = "INTEREST"
)

type PaymentSubType string

func (p PaymentSubType) String() string {
	return string(p)
}

func PaymentSubTypeFromString(src string) PaymentSubType {
	switch src {
	case string(PaymentSubTypePayment):
		return PaymentSubTypePayment
	default:
		return PaymentSubTypeUnknown
	}
}

const (
	PaymentSubTypeUnknown PaymentSubType = ""
	PaymentSubTypePayment PaymentSubType = "PAYMENT"
)

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
