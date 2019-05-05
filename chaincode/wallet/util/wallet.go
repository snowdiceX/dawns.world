package util

import (
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/snowdiceX/dawns.world/chaincode/wallet/log"
)

// Wallet store in hyperledger fabric chain
type Wallet struct {
	Key     string
	Version string `json:"version,omitempty"`
	Chain   string `json:"chain,omitempty"`
	Token   string `json:"token,omitempty"`
	Balance string `json:"balance,omitempty"`
	Height  string `json:"height,omitempty"`
	Address string `json:"address,omitempty"`
	TxID    string `json:"txid,omitempty"`
	Agent   string `json:"agent,omitempty"`
}

// WalletSequence record of wallet chaincode sequence
type WalletSequence struct {
	Version  string `json:"version,omitempty"`
	Sequence string `json:"sequence,omitempty"`
	TxID     string `json:"txid,omitempty"`
	Func     string `json:"func,omitempty"`
	Address  string `json:"address,omitempty"`
	Network  string `json:"network,omitempty"`
	Token    string `json:"token,omitempty"`
	Height   string `json:"height,omitempty"`
}

// NewWallet create a user's lock wallet
func NewWallet(chain, token, address string) *Wallet {
	wallet := &Wallet{
		Chain:   chain,
		Token:   token,
		Address: address}
	wallet.Key = BuildWalletKey(chain, token, address)
	return wallet
}

// Load returns state data of this wallet
func (w *Wallet) Load(stub shim.ChaincodeStubInterface) *ChaincodeError {
	bytes, ccErr := CheckState(stub, w.Key, false)
	if ccErr != nil {
		return ccErr
	}
	if bytes == nil {
		errString := fmt.Sprintf("wallet not found: %s", w.Key)
		log.Error(errString)
		return &ChaincodeError{
			Code:      http.StatusInternalServerError,
			ErrString: errString}
	}
	if err := json.Unmarshal(bytes, w); err != nil {
		errString := fmt.Sprintf("wallet unmarshal error: %v\n  %s\n  json: %s",
			err, w.Key, string(bytes))
		log.Error(errString)
		return &ChaincodeError{
			Code:      http.StatusInternalServerError,
			ErrString: errString}
	}
	return nil
}

// Save state data of this wallet
func (w *Wallet) Save(stub shim.ChaincodeStubInterface) *ChaincodeError {
	bytes, err := json.Marshal(w)
	if err != nil {
		errString := fmt.Sprintf("wallet marshal failed: %v", err)
		log.Errorf(errString)
		return &ChaincodeError{
			Code:      http.StatusInternalServerError,
			ErrString: errString}
	}
	if err = stub.PutState(w.Key, bytes); err != nil {
		errString := fmt.Sprintf("wallet save failed: %v", err)
		return &ChaincodeError{
			Code:      http.StatusInternalServerError,
			ErrString: errString}
	}
	log.Debugf("wallet saved: %s: %s", w.Key, string(bytes))
	return nil
}

// Sum transactions to wallet
func (w *Wallet) Sum(tx *TxRegister) *ChaincodeError {
	var err *ChaincodeError
	if err = preCheck(w, tx); err != nil {
		log.Error(err.Error())
		return err
	}
	balance := new(big.Int)
	balance.SetString(w.Balance[2:], 16)
	if err = sumAmount(balance, tx.Amount); err != nil {
		log.Error(err.Error())
		return err
	}
	if err = sumGas(balance, tx.GasUsed, tx.GasPrice); err != nil {
		log.Error(err.Error())
		return err
	}
	w.Balance = fmt.Sprintf("0x%s", balance.Text(16))
	log.Info("wallet balance: ", w.Balance)
	return nil
}

func preCheck(w *Wallet, tx *TxRegister) *ChaincodeError {
	if tx == nil {
		return &ChaincodeError{
			Code:      http.StatusBadRequest,
			ErrString: "Tx is nil"}
	}
	if !strings.EqualFold(w.Chain, tx.Chain) ||
		!strings.EqualFold(w.Token, tx.Token) ||
		!strings.EqualFold(w.Address, tx.WalletAddress) {
		return &ChaincodeError{
			Code: http.StatusInternalServerError,
			ErrString: fmt.Sprintf(
				"wrong wallet: %s, %s, %s/ %s, %s, %s",
				w.Chain, w.Token, w.Address,
				tx.Chain, tx.Token, tx.WalletAddress)}
	}
	if !strings.HasPrefix(tx.Amount, "0x") && tx.Amount != "" {
		return &ChaincodeError{
			Code: http.StatusBadRequest,
			ErrString: fmt.Sprintf(
				`amount "%s" should be prefixed with 0x`,
				tx.Amount)}
	}
	return nil
}

// enough determines if wallet balance enough for the amount needed
func enough(balance *big.Int, amount *big.Int) bool {
	if balance.CmpAbs(amount) < 0 {
		return false
	}
	return true
}

func sumAmount(balance *big.Int, amount string) *ChaincodeError {
	if len(amount) > 0 {
		a := new(big.Int)
		a.SetString(amount[2:], 16)
		if a.Sign() < 0 && enough(balance, a) {
			return &ChaincodeError{
				Code: http.StatusBadRequest,
				ErrString: fmt.Sprintf(
					`balance "%s" not enough: %s`,
					w.Balance, amount)}
		}
		balance.Add(balance, a)
	}
	return nil
}

func sumGas(balance *big.Int, gasUsed, gasPrice string) *ChaincodeError {
	if len(tx.GasUsed) > 0 && len(tx.GasPrice) > 0 {
		g := new(big.Int)
		g.SetString(tx.GasUsed[2:], 16)
		if g.Sign() <= 0 {
			return &ChaincodeError{
				Code: http.StatusBadRequest,
				ErrString: fmt.Sprintf(
					`gas used "%s" should be positive`,
					tx.GasUsed)}
		}
		gp := new(big.Int)
		gp.SetString(tx.GasPrice[2:], 16)
		if gp.Sign() <= 0 {
			return &ChaincodeError{
				Code: http.StatusBadRequest,
				ErrString: fmt.Sprintf(
					`gas price "%s" should be positive`,
					tx.GasPrice)}
		}
		g.Mul(g, gp)
		balance.Sub(balance, g)
	}
}
