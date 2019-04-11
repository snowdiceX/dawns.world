package main

import (
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"strconv"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/snowdiceX/dawns.world/chaincode/log"
	"github.com/snowdiceX/dawns.world/chaincode/util"
)

type registerTx struct {
	Chain    string `json:"chain,omitempty"`
	Token    string `json:"token,omitempty"`
	Contract string `json:"contract,omitempty"`
	From     string `json:"from,omitempty"`
	To       string `json:"to,omitempty"`
	Amount   string `json:"amount,omitempty"`
	Txhash   string `json:"txhash,omitempty"`
	Height   string `json:"height,omitempty"`
}

func (w *WalletChaincode) registerTransaction(
	stub shim.ChaincodeStubInterface, args []string) pb.Response {
	tx := &registerTx{}
	err := json.Unmarshal([]byte(args[0]), tx)
	if err != nil {
		log.Errorf("parse error: %v\n    json: %s\n", err, args[0])
		return util.Error(http.StatusBadRequest, fmt.Sprintf(
			"register failed: %v", err))
	}
	logKey := util.BuildLogTransactionKey(tx.Chain, tx.Token, tx.Height, tx.Txhash)
	bytes, ccErr := checkState(stub, logKey, true)
	if ccErr != nil {
		return util.Error(ccErr.Code(),
			fmt.Sprintf("register failed: %v", ccErr))
	}
	if err = stub.PutState(logKey, []byte(args[0])); err != nil {
		log.Errorf("put state error: %s: %v", logKey, err)
		return util.Error(http.StatusInternalServerError,
			fmt.Sprintf("register failed: %v", err))
	}
	walletKey := util.BuildWalletKey(tx.Chain, tx.Token, tx.To)
	bytes, ccErr = checkState(stub, walletKey, false)
	if ccErr != nil {
		return util.Error(ccErr.Code(),
			fmt.Sprintf("register failed: %v", ccErr))
	}
	if bytes == nil {
		log.Errorf("wallet not found: %s", walletKey)
		return util.Error(http.StatusInternalServerError,
			fmt.Sprintf("wallet not found: %s", walletKey))
	}
	wallet := &util.Wallet{}
	if err = json.Unmarshal(bytes, wallet); err != nil {
		log.Errorf("parse state error: %v\n    json: %s", err, string(bytes))
		return util.Error(http.StatusInternalServerError, fmt.Sprintf(
			"register failed: %v", err))
	}
	wallet, err = sumTransaction(wallet, tx)
	if bytes, err = json.Marshal(wallet); err != nil {
		log.Errorf("json marshal error: %s: %v", walletKey, err)
		return util.Error(http.StatusInternalServerError,
			fmt.Sprintf("register failed: %v", err))
	}
	if err = stub.PutState(walletKey, bytes); err != nil {
		log.Errorf("put state error: %s: %v", logKey, err)
		return util.Error(http.StatusInternalServerError,
			fmt.Sprintf("register failed: %v", err))
	}
	ret := &util.ChainResult{Code: 200, Message: "OK"}
	ret.Result = tx
	return util.Success(ret)
}

func (w *WalletChaincode) queryTransaction(
	stub shim.ChaincodeStubInterface, args []string) pb.Response {
	what := args[0]
	switch what {
	case "sequence":
		{
			if len(args) < 2 {
				return shim.Error("insufficient parameters")
			}
			// queries a transaction by sequence
			return w.queryTransactionBySequence(stub, args[1])
		}
	case "page":
		{
			if len(args) < 7 {
				return shim.Error("insufficient parameters")
			}
			// query transactions
			return w.queryTransactions(stub, args[1:])
		}
	}
	return util.Error(http.StatusBadRequest, fmt.Sprintf(
		"query transaction failed: Query what? %s", what))
}

func (w *WalletChaincode) queryTransactions(
	stub shim.ChaincodeStubInterface, args []string) pb.Response {
	typeTx, chain, token, pageNumHex, pageSizeHex, walletAddress :=
		args[0], args[1], args[2], args[3], args[4], args[5]
	log.Debug(chain, token, pageNumHex, pageSizeHex, walletAddress)
	pageSize, _ := strconv.ParseInt(pageSizeHex, 10, 32)
	startKey := util.BuildLogTransactionStartKey(chain, token, "")
	logIterator, meta, err := stub.GetStateByRangeWithPagination(
		startKey, startKey+"a", int32(pageSize), "")
	if err != nil {
		log.Errorf("paging %s transactions error: %v", typeTx, err)
		return util.Error(http.StatusInternalServerError,
			fmt.Sprintf("paging %s transactions failed: %v", typeTx, err))
	}
	defer logIterator.Close()
	buffer, err := constructPageJSON(logIterator, meta)
	if err != nil {
		log.Errorf("paging %s transactions error: %v", typeTx, err)
		return util.Error(http.StatusInternalServerError,
			fmt.Sprintf("paging %s transactions failed: %v", typeTx, err))
	}
	log.Debugf("paging %s transactions result:\n%s\n",
		typeTx, buffer.String())
	ret := &util.ChainResult{Code: 200, Message: "OK"}
	ret.Result = buffer.String()
	return util.Success(ret)
}

// queryTransactionBySequence queries a transaction by sequence
func (w *WalletChaincode) queryTransactionBySequence(
	stub shim.ChaincodeStubInterface, sequence string) pb.Response {
	seq, _ := strconv.ParseUint(sequence, 10, 64)
	sequenceKey := util.BuildSequenceKey(seq)

	// Get the state from the ledger
	txBytes, err := stub.GetState(sequenceKey)
	if err != nil {
		jsonResp := "{\"Error\":\"failed to get state for " + sequenceKey + "\"}"
		return shim.Error(jsonResp)
	}

	if txBytes == nil {
		jsonResp := "{\"Error\":\"nil amount for " + sequenceKey + "\"}"
		return shim.Error(jsonResp)
	}

	fmt.Printf("query Tx by sequence response:%s\n", string(txBytes))
	return shim.Success(txBytes)
}

func sumTransaction(wallet *util.Wallet, tx *registerTx) (*util.Wallet, *ChaincodeError) {
	if wallet == nil || tx == nil ||
		!strings.EqualFold(wallet.Chain, tx.Chain) ||
		!strings.EqualFold(wallet.Token, tx.Token) ||
		!strings.EqualFold(wallet.Address, tx.To) {
		return nil, &ChaincodeError{
			code: http.StatusInternalServerError,
			errString: fmt.Sprintf("wrong wallet: %s, %s, %s/ %s, %s, %s",
				wallet.Chain, wallet.Token, wallet.Address,
				tx.Chain, tx.Token, tx.To)}
	}
	if !strings.HasPrefix(tx.Amount, "0x") {
		return nil, &ChaincodeError{
			code: http.StatusBadRequest,
			errString: fmt.Sprintf(
				`amount "%s" should be prefixed with 0x`, tx.Amount)}
	}
	balance := new(big.Int)
	x := new(big.Int)
	x.SetString(wallet.Balance[2:], 16)
	y := new(big.Int)
	y.SetString(tx.Amount[2:], 16)
	balance = balance.Add(x, y)
	wallet.Balance = fmt.Sprintf("0x%s", balance.Text(16))
	return wallet, nil
}
