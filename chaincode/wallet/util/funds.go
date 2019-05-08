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

// RecordFunds funds deposit record of wallet address
type RecordFunds struct {
	Key           string
	Version       string `json:"version,omitempty"`
	FundsTokenKey string `json:"fundsTokenKey,omitempty"`
	WalletAddress string `json:"address,omitempty"`
	Chain         string `json:"chain,omitempty"`
	Token         string `json:"token,omitempty"`
	FundsHash     string `json:"fundsHash,omitempty"`
	Balance       string `json:"balance,omitempty"`
}

// NewRecordFunds create a funds record
func NewRecordFunds(fundsTokenKey, walletAddress string) *RecordFunds {
	vs := strings.Split(fundsTokenKey, "-")
	if len(vs) < 4 {
		vs = []string{"wrong", "funds", "token", "key"}
	}
	rec := &RecordFunds{
		Version:       ChaincodeVersion,
		FundsTokenKey: fundsTokenKey,
		WalletAddress: walletAddress,
		Chain:         vs[1],
		Token:         vs[2],
		FundsHash:     vs[3]}
	rec.Key = rec.buildKey()
	return rec
}

// buildKey returns state key of this record
func (r *RecordFunds) buildKey() string {
	return BuildRecordFundsKey(r.Chain, r.Token, r.FundsHash, r.WalletAddress)
}

// Load returns state data of this record
func (r *RecordFunds) Load(stub shim.ChaincodeStubInterface) *ChaincodeError {
	bytes, ccErr := CheckState(stub, r.Key, false)
	if ccErr != nil {
		errString := fmt.Sprintf(
			"funds record load failed: %s; %v", r.Key, ccErr)
		log.Errorf(errString)
		return ccErr
	}
	if bytes == nil {
		errString := fmt.Sprintf(
			"funds record is nil: %s; %s", r.Key, string(bytes))
		log.Warn(errString)
	} else {
		log.Debugf("funds record load: %s; %s", r.Key, string(bytes))
		if err := json.Unmarshal(bytes, r); err != nil {
			errString := fmt.Sprintf(
				"funds record unmarshal failed: %s; %v", r.Key, err)
			log.Errorf(errString)
			return &ChaincodeError{Code: http.StatusInternalServerError,
				ErrString: errString}
		}
	}
	return nil
}

// Save state data of this record
func (r *RecordFunds) Save(stub shim.ChaincodeStubInterface) *ChaincodeError {
	bytes, err := json.Marshal(r)
	if err != nil {
		errString := fmt.Sprintf("funds record marshal error: %v", err)
		log.Errorf(errString)
		return &ChaincodeError{
			Code:      http.StatusInternalServerError,
			ErrString: errString}
	}
	if err = stub.PutState(r.Key, bytes); err != nil {
		errString := fmt.Sprintf("funds record save error: %v", err)
		log.Errorf(errString)
		return &ChaincodeError{
			Code:      http.StatusInternalServerError,
			ErrString: errString}
	}
	log.Debugf("funds record saved: %s: %s", r.Key, string(bytes))
	return nil
}

// Add sets funds record balance to the sum r.Balance + tx.Amount()
func (r *RecordFunds) Add(tx TransactionLog) *ChaincodeError {
	var err *ChaincodeError
	balance := new(big.Int)
	if r.Balance != "" {
		balance.SetString(r.Balance[2:], 16)
	}
	if err = addAmount(balance, tx.AmountHex()); err != nil {
		log.Error(err.Error())
		return nil
	}
	r.Balance = fmt.Sprintf("0x%s", balance.Text(16))
	log.Info("funds record balance: ", r.Balance)
	return nil
}

func addAmount(balance *big.Int, amount string) *ChaincodeError {
	a := new(big.Int)
	a.SetString(amount[2:], 16)
	if a.Sign() < 0 {
		return &ChaincodeError{
			Code: http.StatusBadRequest,
			ErrString: fmt.Sprintf(
				`amount "%s" should be positive`,
				amount)}
	}
	balance.Add(balance, a)
	return nil
}

// Sub sets funds record balance to the difference r.Balance - tx.Amount()
func (r *RecordFunds) Sub(tx TransactionLog) *ChaincodeError {
	var err *ChaincodeError
	balance := new(big.Int)
	balance.SetString(r.Balance[2:], 16)
	if err = subAmount(balance, tx.AmountHex()); err != nil {
		log.Error(err.Error())
		return nil
	}
	r.Balance = fmt.Sprintf("0x%s", balance.Text(16))
	log.Info("funds record balance: ", r.Balance)
	return nil
}
