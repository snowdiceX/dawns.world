package util

// TransactionLog trasaction of wallet
type TransactionLog interface {
	Chain() string
	Token() string
	Address() string
	AmountHex() string
	GasUsedHex() string
	GasPriceHex() string
}
