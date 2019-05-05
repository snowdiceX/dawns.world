package util

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/snowdiceX/dawns.world/chaincode/wallet/log"
)

// RecordFunds funds deposit record of wallet address
type RecordFunds struct {
	Key           string
	Version       string `json:"version,omitempty"`
	FundsTokenKey string `json:"FundsTokenKey,omitempty"`
	WalletAddress string `json:"walletAddress,omitempty"`
	Chain         string `json:"chain,omitempty"`
	Token         string `json:"token,omitempty"`
	FundsHash     string `json:"FundsHash,omitempty"`
	Balance       string `json:"balance,omitempty"`
}

// NewRecordFunds create a funds record
func NewRecordFunds(fundsTokenKey, walletAddress string) *RecordFunds {
	vs := strings.Split(fundsTokenKey, "-")
	rec := &RecordFunds{
		Version:       ChaincodeVersion,
		FundsTokenKey: fundsTokenKey,
		WalletAddress: walletAddress,
		Chain:         vs[0],
		Token:         vs[1],
		FundsHash:     vs[2]}
	rec.Key = rec.buildKey()
	return rec
}

// buildKey returns state key of this record
func (r *RecordFunds) buildKey() string {
	return BuildRecordFundsKey(r.FundsTokenKey, r.WalletAddress)
}

// Load returns state data of this record
func (r *RecordFunds) Load(stub shim.ChaincodeStubInterface) *ChaincodeError {
	bytes, ccErr := CheckState(stub, r.Key, false)
	if ccErr != nil {
		return ccErr
	}
	rec := &RecordFunds{}
	if err := json.Unmarshal(bytes, rec); err != nil {
		errString := fmt.Sprintf("funds record load failed: %v", err)
		log.Errorf(errString)
		return &ChaincodeError{Code: http.StatusInternalServerError,
			ErrString: errString}
	}
	if !strings.EqualFold(r.Hash, rec.Hash) ||
		!strings.EqualFold(r.TokenKey, rec.TokenKey) ||
		!strings.EqualFold(r.WalletAddress, rec.WalletAddress) {
		errString := fmt.Sprintf("wrong data: %s", string(bytes))
		log.Errorf(errString)
		return &ChaincodeError{Code: http.StatusInternalServerError,
			ErrString: errString}
	}
	r.Balance = rec.Balance
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
