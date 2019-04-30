package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync/atomic"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/ledger/queryresult"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/snowdiceX/dawns.world/chaincode/wallet/log"
	"github.com/snowdiceX/dawns.world/chaincode/wallet/util"
)

func (w *WalletChaincode) registerBlock(
	stub shim.ChaincodeStubInterface, args []string) pb.Response {
	log.Info("register block: ", args[0])
	var ccErr *util.ChaincodeError
	var block *util.BlockRegister
	if block, ccErr = checkBlock(args[0]); ccErr != nil {
		return util.Error(ccErr.Code, ccErr.Error())
	}
	if ccErr = w.checkInSequence(block); ccErr != nil {
		return util.Error(ccErr.Code, ccErr.Error())
	}
	var count int
	var wallet *util.Wallet
	for _, tx := range block.Txs {
		tx.Height = block.Height
		if ccErr = checkTransactionLog(stub, tx); ccErr != nil {
			continue
		}
		if wallet, ccErr = checkWallet(stub, tx); ccErr != nil {
			continue
		}
		ccErr = wallet.Sum(tx)
		if ccErr != nil {
			log.Errorf("wallet sum error: %s: %v", wallet.Key, ccErr)
			return util.Error(ccErr.Code,
				fmt.Sprintf("register failed: %v", ccErr))
		}
		if ccErr = saveWallet(stub, wallet); ccErr != nil {
			continue
		}
		count++
	}
	if ccErr != nil && count == 0 {
		return util.Error(ccErr.Code, ccErr.Error())
	}
	ret := &util.ChainResult{Code: 200, Message: "OK"}
	ret.Result = count
	return util.Success(ret)
}

func checkBlock(blockJSON string) (*util.BlockRegister, *util.ChaincodeError) {
	block := &util.BlockRegister{}
	err := json.Unmarshal([]byte(blockJSON), block)
	if err != nil {
		errString := fmt.Sprintf("parse error: %v\n    json: %s",
			err, blockJSON)
		log.Error(errString)
		return nil, &util.ChaincodeError{
			Code:      http.StatusBadRequest,
			ErrString: errString}
	}
	return block, nil
}

func (w *WalletChaincode) checkInSequence(
	block *util.BlockRegister) *util.ChaincodeError {
	h, err := strconv.ParseUint(block.Height[2:], 16, 64)
	if err != nil {
		errString := fmt.Sprintf("parse block height failed: %s, %v",
			block.Height, err)
		log.Error(errString)
		return &util.ChaincodeError{
			Code:      http.StatusInternalServerError,
			ErrString: errString}
	}
	if !atomic.CompareAndSwapUint64(&w.InSequence, h-1, h) {
		errString := fmt.Sprintf("update in-sequence failed: %d; %d",
			w.InSequence, h)
		log.Error(errString)
		return &util.ChaincodeError{
			Code:      http.StatusInternalServerError,
			ErrString: errString}
	}
	return nil
}

func checkTransactionLog(
	stub shim.ChaincodeStubInterface,
	tx *util.TxRegister) *util.ChaincodeError {
	logKey := util.BuildLogTransactionKey(tx.Chain, tx.Token, tx.Height, tx.Txhash)
	bytes, ccErr := checkState(stub, logKey, true)
	if ccErr != nil {
		return ccErr
	}
	var err error
	if bytes, err = json.Marshal(tx); err != nil {
		errString := fmt.Sprintf("transaction marshal failed: %s %v",
			logKey, err)
		log.Error(errString)
		return &util.ChaincodeError{
			Code:      http.StatusInternalServerError,
			ErrString: errString}
	}
	if err = stub.PutState(logKey, bytes); err != nil {
		errString := fmt.Sprintf("transaction log put state error: %s: %v",
			logKey, err)
		log.Error(errString)
		return &util.ChaincodeError{
			Code:      http.StatusInternalServerError,
			ErrString: errString}
	}
	return nil
}

func checkWallet(stub shim.ChaincodeStubInterface,
	tx *util.TxRegister) (*util.Wallet, *util.ChaincodeError) {
	walletKey := util.BuildWalletKey(tx.Chain, tx.Token, tx.To)
	bytes, ccErr := checkState(stub, walletKey, false)
	if ccErr != nil {
		return nil, ccErr
	}
	if bytes == nil {
		walletKey = util.BuildWalletKey(tx.Chain, tx.Token, tx.From)
		bytes, ccErr = checkState(stub, walletKey, false)
		if ccErr != nil {
			return nil, ccErr
		}
	}
	if bytes == nil {
		errString := fmt.Sprintf("wallet not found: %s", walletKey)
		log.Error(errString)
		return nil, &util.ChaincodeError{
			Code:      http.StatusInternalServerError,
			ErrString: errString}
	}
	var err error
	wallet := &util.Wallet{Key: walletKey}
	if err = json.Unmarshal(bytes, wallet); err != nil {
		errString := fmt.Sprintf("wallet unmarshal error: %v\n    json: %s",
			err, string(bytes))
		log.Error(errString)
		return nil, &util.ChaincodeError{
			Code:      http.StatusInternalServerError,
			ErrString: errString}
	}
	return wallet, nil
}

func saveWallet(stub shim.ChaincodeStubInterface,
	wallet *util.Wallet) *util.ChaincodeError {
	var bytes []byte
	var err error
	if bytes, err = json.Marshal(wallet); err != nil {
		errString := fmt.Sprintf("wallet marshal error: %s: %v",
			wallet.Key, err)
		log.Error(errString)
		return &util.ChaincodeError{
			Code:      http.StatusInternalServerError,
			ErrString: errString}
	}
	if err = stub.PutState(wallet.Key, bytes); err != nil {
		errString := fmt.Sprintf("wallet put state error: %s: %v",
			wallet.Key, err)
		log.Error(errString)
		return &util.ChaincodeError{
			Code:      http.StatusInternalServerError,
			ErrString: errString}
	}
	log.Infof("wallet update: %s: %s", wallet.Key, string(bytes))
	return nil
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
			return w.paginationTransactions(stub, args[1:])
		}
	}
	return util.Error(http.StatusBadRequest, fmt.Sprintf(
		"query transaction failed: Query what? %s", what))
}

func (w *WalletChaincode) paginationTransactions(
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
	page, err := constructPage(logIterator, meta)
	if err != nil {
		log.Errorf("paging %s transactions error: %v", typeTx, err)
		return util.Error(http.StatusInternalServerError,
			fmt.Sprintf("paging %s transactions failed: %v", typeTx, err))
	}
	ret := &util.ChainResult{Code: 200, Message: "OK"}
	ret.Result = page
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

// construct a page struct from iterator
func constructPage(iterator shim.StateQueryIteratorInterface,
	metadata *pb.QueryResponseMetadata) (*util.Pagination, error) {
	var recs []interface{}
	var rec *queryresult.KV
	var err error
	for iterator.HasNext() {
		rec, err = iterator.Next()
		if err != nil {
			return nil, err
		}
		tx := &util.TxRegister{}
		err = json.Unmarshal(rec.Value, tx)
		if err != nil {
			return nil, err
		}
		tx.Key = rec.Key
		recs = append(recs, tx)
	}
	page := &util.Pagination{}
	page.Metadata = metadata
	page.Records = recs
	return page, nil
}
