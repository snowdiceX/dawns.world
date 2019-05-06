package util

type TransactionLog interface {
	Chain() string
	Token() string
	Address() string
	AmountHex() string
	GasUsedHex() string
	GasPriceHex() string
}
