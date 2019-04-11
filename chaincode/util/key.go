package util

import (
	"fmt"
	"strings"
)

const (
	// TagWallet key prefix of wallet address
	TagWallet = "Wallet"
	// TagSequence key prefix of transaction sequence
	TagSequence = "Sequence"
	// TagFunds key prefix of funds
	TagFunds = "Funds"
	// TagLogRegisteredTx key prefix of registered transaction log
	TagLogRegisteredTx = "LogRegisteredTx"
)

// BuildAccountKey key of user account
func BuildAccountKey(chain, token, accountID, address string) string {
	return fmt.Sprintf("%s-%s-%s-%s-%s",
		TagWallet, chain, token, accountID, address)
}

// BuildWalletKey key of wallet address
func BuildWalletKey(network, token, address string) string {
	return fmt.Sprintf("%s-%s-%s-%s", TagWallet, network, token, address)
}

// BuildSequenceKey key of tx sequence
func BuildSequenceKey(seq uint64) string {
	return fmt.Sprintf("%s-%d", TagSequence, seq)
}

// BuildFundsKey key of funds store in form:
// {"base":"base token amount", "token":"Acceptance token amount"}
func BuildFundsKey(baseNetwork, baseToken, network, token string) string {
	return fmt.Sprintf("%s-%s-%s-%s-%s",
		TagFunds, baseNetwork, baseToken, network, token)
}

// BuildFundsAddressKey key of transfer amount from a wallet address
func BuildFundsAddressKey(baseNetwork, baseToken, network, token,
	address string) string {
	return fmt.Sprintf("%s-%s-%s-%s-%s-%s",
		TagFunds, baseNetwork, baseToken, network, token, address)
}

// BuildLogTransactionKey key of transfer transaction
func BuildLogTransactionKey(network, token, height, txhash string) string {
	return BuildLogTransactionStartKey(network, token, height, txhash)
}

// BuildLogTransactionStartKey key for query registered transactions
func BuildLogTransactionStartKey(network, token string, args ...string) string {
	return fmt.Sprintf("%s-%s-%s-%s",
		TagLogRegisteredTx, network, token, strings.Join(args, "-"))
}
