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
	// TagFundsBase key prefix of funds base token
	TagFundsBase = "Base-Funds"
	// TagFundsAccept key prefix of funds accept token
	TagFundsAccept = "Accept-Funds"
	// TagToken key prefix of token
	TagToken = "Token"
	// TagLogRegisteredTx key prefix of registered transaction log
	TagLogRegisteredTx = "Log-RegisteredTx"
	// TagRecordFunds key prefix of funds deposit and withdraw record
	TagRecordFunds = "Record-Funds"
)

// BuildAccountKey key of user account
func BuildAccountKey(chain, token, accountID, address string) string {
	return fmt.Sprintf("%s-%s-%s-%s-%s",
		TagWallet, chain, token, accountID, address)
}

// BuildWalletKey key of wallet address
func BuildWalletKey(network, token, address string) string {
	return fmt.Sprintf("%s-%s-%s-%s", TagWallet,
		strings.ToUpper(network),
		strings.ToUpper(token),
		strings.ToUpper(address))
}

// BuildSequenceKey key of tx sequence
func BuildSequenceKey(seq uint64) string {
	return fmt.Sprintf("%s-%d", TagSequence, seq)
}

// BuildFundsKey key of funds
func BuildFundsKey(baseChain, baseToken, chain, token string) string {
	return BuildFundsStartKey(baseChain, baseToken, chain, token)
}

// BuildFundsStartKey key for query funds list
func BuildFundsStartKey(args ...string) string {
	return fmt.Sprintf("%s-%s",
		TagFunds, strings.Join(args, "-"))
}

// BuildFundsBaseKey key of funds base token
func BuildFundsBaseKey(chain, token, fundsHash string) string {
	return fmt.Sprintf("%s-%s-%s-%s",
		TagFundsBase, chain, token, fundsHash)
}

// BuildFundsAcceptKey key of funds accept token
func BuildFundsAcceptKey(chain, token, fundsHash string) string {
	return fmt.Sprintf("%s-%s-%s-%s",
		TagFundsAccept, chain, token, fundsHash)
}

// BuildTokenKey key of token,
// store in form:
// {"chain":"chain name",
//  "token":"token name",
//  "symbol":"token symbol",
//  "decimals":"decimals of token",
//  "contract":"contract address",
//  "txhash":"transaction hash in fabric"}
func BuildTokenKey(network, token string) string {
	return BuildTokenStartKey(network, token)
}

// BuildTokenStartKey key for query token list
func BuildTokenStartKey(args ...string) string {
	return fmt.Sprintf("%s-%s",
		TagToken, strings.Join(args, "-"))
}

// BuildFundsAddressKey key of transfer amount from a wallet address
func BuildFundsAddressKey(baseNetwork, baseToken, network, token,
	address string) string {
	return fmt.Sprintf("%s-%s",
		BuildFundsKey(baseNetwork, baseToken, network, token), address)
}

// BuildLogTransactionKey key of transfer transaction
func BuildLogTransactionKey(network, token, height, txhash string) string {
	return BuildLogTransactionStartKey(network, token,
		strings.ToUpper(height),
		strings.ToUpper(txhash))
}

// BuildLogTransactionStartKey key for query registered transaction list
func BuildLogTransactionStartKey(network, token string, args ...string) string {
	return fmt.Sprintf("%s-%s-%s-%s",
		TagLogRegisteredTx,
		strings.ToUpper(network),
		strings.ToUpper(token),
		strings.Join(args, "-"))
}

// BuildRecordFundsKey key of funds deposit and withdraw record
func BuildRecordFundsKey(tokenKey, walletAddress string) string {
	return fmt.Sprintf("%s-%s-%s",
		TagRecordFunds, tokenKey, walletAddress)
}
