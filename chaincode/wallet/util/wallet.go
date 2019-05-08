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

// Sum sets wallet balance to the sum w.Balance + tx.Amount()
func (w *Wallet) Sum(tx TransactionLog) *ChaincodeError {
	var err *ChaincodeError
	if err = w.preCheck(tx); err != nil {
		log.Error(err.Error())
		return err
	}
	balance := new(big.Int)
	balance.SetString(w.Balance[2:], 16)
	if err = sumGasFee(balance, tx.GasUsedHex(), tx.GasPriceHex()); err != nil {
		log.Error(err.Error())
		return err
	}
	if err = sumAmount(balance, tx.AmountHex()); err != nil {
		log.Error(err.Error())
		return err
	}
	w.Balance = fmt.Sprintf("0x%s", balance.Text(16))
	log.Info("wallet balance: ", w.Balance)
	return nil
}

// Sub sets wallet balance to the difference w.Balance - tx.Amount()
func (w *Wallet) Sub(tx TransactionLog) *ChaincodeError {
	var err *ChaincodeError
	if err = w.preCheck(tx); err != nil {
		log.Error(err.Error())
		return err
	}
	balance := new(big.Int)
	balance.SetString(w.Balance[2:], 16)
	if err = sumGasFee(balance, tx.GasUsedHex(), tx.GasPriceHex()); err != nil {
		log.Error(err.Error())
		return err
	}
	if err = subAmount(balance, tx.AmountHex()); err != nil {
		log.Error(err.Error())
		return nil
	}
	w.Balance = fmt.Sprintf("0x%s", balance.Text(16))
	log.Info("wallet balance: ", w.Balance)
	return nil
}

func (w *Wallet) preCheck(tx TransactionLog) *ChaincodeError {
	if tx == nil {
		return &ChaincodeError{
			Code:      http.StatusBadRequest,
			ErrString: "Tx is nil"}
	}
	if !strings.EqualFold(w.Chain, tx.Chain()) ||
		!strings.EqualFold(w.Token, tx.Token()) ||
		!strings.EqualFold(w.Address, tx.Address()) {
		return &ChaincodeError{
			Code: http.StatusInternalServerError,
			ErrString: fmt.Sprintf(
				"wrong wallet: %s, %s, %s / %s, %s, %s",
				w.Chain, w.Token, w.Address,
				tx.Chain(), tx.Token(), tx.Address())}
	}
	if len(tx.AmountHex()) > 0 &&
		!strings.HasPrefix(tx.AmountHex(), "0x") {
		return &ChaincodeError{
			Code: http.StatusBadRequest,
			ErrString: fmt.Sprintf(
				`amount "%s" should be prefixed with 0x`,
				tx.AmountHex())}
	}
	if len(tx.GasUsedHex()) > 0 &&
		!strings.HasPrefix(tx.GasUsedHex(), "0x") {
		return &ChaincodeError{
			Code: http.StatusBadRequest,
			ErrString: fmt.Sprintf(
				`gas used "%s" should be prefixed with 0x`,
				tx.GasUsedHex())}
	}
	if len(tx.GasPriceHex()) > 0 &&
		!strings.HasPrefix(tx.GasPriceHex(), "0x") {
		return &ChaincodeError{
			Code: http.StatusBadRequest,
			ErrString: fmt.Sprintf(
				`gas price "%s" should be prefixed with 0x`,
				tx.GasPriceHex())}
	}
	return nil
}

// enough determines if wallet balance enough for the amount needed
func enough(balance *big.Int, amount *big.Int) bool {
	return balance.CmpAbs(amount) >= 0
}

func sumAmount(balance *big.Int, amount string) *ChaincodeError {
	if len(amount) > 0 {
		a := new(big.Int)
		a.SetString(amount[2:], 16)
		if a.Sign() < 0 && !enough(balance, a) {
			return &ChaincodeError{
				Code: http.StatusBadRequest,
				ErrString: fmt.Sprintf(
					`balance not enough for amount: %s`,
					amount)}
		}
		balance.Add(balance, a)
	}
	return nil
}

func subAmount(balance *big.Int, amount string) *ChaincodeError {
	a := new(big.Int)
	a.SetString(amount[2:], 16)
	if a.Sign() < 0 {
		return &ChaincodeError{
			Code: http.StatusBadRequest,
			ErrString: fmt.Sprintf(
				`amount "%s" should be positive`,
				amount)}
	}
	if !enough(balance, a) {
		return &ChaincodeError{
			Code: http.StatusBadRequest,
			ErrString: fmt.Sprintf(
				`balance not enough for amount: %s`,
				amount)}
	}
	balance.Sub(balance, a)
	return nil
}

func sumGasFee(balance *big.Int, gasUsed, gasPrice string) *ChaincodeError {
	if len(gasUsed) > 0 && len(gasPrice) > 0 {
		g := new(big.Int)
		g.SetString(gasUsed[2:], 16)
		if g.Sign() <= 0 {
			return &ChaincodeError{
				Code: http.StatusBadRequest,
				ErrString: fmt.Sprintf(
					`gas used "%s" should be positive`,
					gasUsed)}
		}
		gp := new(big.Int)
		gp.SetString(gasPrice[2:], 16)
		if gp.Sign() <= 0 {
			return &ChaincodeError{
				Code: http.StatusBadRequest,
				ErrString: fmt.Sprintf(
					`gas price "%s" should be positive`,
					gasPrice)}
		}
		g.Mul(g, gp)
		if !enough(balance, g) {
			return &ChaincodeError{
				Code: http.StatusBadRequest,
				ErrString: fmt.Sprintf(
					`balance 0x%s not enough for fee: %s * %s`,
					balance.Text(16), gasUsed, gasPrice)}
		}
		balance.Sub(balance, g)
	}
	return nil
}
