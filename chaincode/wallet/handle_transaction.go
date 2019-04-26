package main

import (
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"sync/atomic"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/ledger/queryresult"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/snowdiceX/dawns.world/chaincode/wallet/log"
	"github.com/snowdiceX/dawns.world/chaincode/wallet/util"
)

type paginationTxs struct {
	Records  []interface{}             `json:"records,omitempty"`
	Metadata *pb.QueryResponseMetadata `json:"metadata,omitempty"`
}

// TxRegister registered Tx
type TxRegister struct {
	Key      string `json:"key,omitempty"`
	Chain    string `json:"chain,omitempty"`
	Token    string `json:"token,omitempty"`
	Contract string `json:"contract,omitempty"`
	From     string `json:"from,omitempty"`
	To       string `json:"to,omitempty"`
	Amount   string `json:"amount,omitempty"`
	GasUsed  string `json:"gasUsed,omitempty"`
	GasPrice string `json:"gasPrice,omitempty"`
	Txhash   string `json:"txhash,omitempty"`
	Height   string `json:"height,omitempty"`
	Status   string `json:"status,omitempty"`
}

// BlockRegister registered block
type BlockRegister struct {
	Height string        `json:"height,omitempty"`
	Txs    []*TxRegister `json:"transactions,omitempty"`
}

func (w *WalletChaincode) registerBlock(
	stub shim.ChaincodeStubInterface, args []string) pb.Response {
	log.Info("register block: ", args[0])
	block := &BlockRegister{}
	err := json.Unmarshal([]byte(args[0]), block)
	if err != nil {
		log.Errorf("parse error: %v\n    json: %s\n", err, args[0])
		return util.Error(http.StatusBadRequest, fmt.Sprintf(
			"register failed: %v", err))
	}
	h, err := strconv.ParseUint(block.Height[2:], 16, 64)
	if err != nil {
		log.Errorf("register failed: %v", err)
		return util.Error(http.StatusInternalServerError,
			fmt.Sprintf("register failed: %v", err))
	}
	if !atomic.CompareAndSwapUint64(&w.InSequence, h-1, h) {
		log.Errorf("register block update InSequence failed: %v", err)
		return util.Error(http.StatusInternalServerError,
			fmt.Sprintf("register block update InSequence failed:: %d; %d",
				w.InSequence, h))
	}
	var bytes []byte
	var ccErr *ChaincodeError
	for _, tx := range block.Txs {
		tx.Height = block.Height
		logKey := util.BuildLogTransactionKey(tx.Chain, tx.Token, tx.Height, tx.Txhash)
		bytes, ccErr = checkState(stub, logKey, true)
		if ccErr != nil {
			if ccErr.Code() == http.StatusConflict {
				continue
			}
			log.Errorf("register failed: %v", ccErr)
			return util.Error(ccErr.Code(),
				fmt.Sprintf("register failed: %v", ccErr))
		}
		if bytes, err = json.Marshal(tx); err != nil {
			log.Errorf("register failed: %v", err)
			return util.Error(http.StatusInternalServerError,
				fmt.Sprintf("register failed: %v", err))
		}
		if err = stub.PutState(logKey, bytes); err != nil {
			log.Errorf("put state error: %s: %v", logKey, err)
			return util.Error(http.StatusInternalServerError,
				fmt.Sprintf("register failed: %v", err))
		}
		walletKey := util.BuildWalletKey(tx.Chain, tx.Token, tx.To)
		bytes, ccErr = checkState(stub, walletKey, false)
		if ccErr != nil {
			log.Errorf("register failed: %v", ccErr)
			return util.Error(ccErr.Code(),
				fmt.Sprintf("register failed: %v", ccErr))
		}
		if bytes == nil {
			walletKey = util.BuildWalletKey(tx.Chain, tx.Token, tx.From)
			bytes, ccErr = checkState(stub, walletKey, false)
			if ccErr != nil {
				log.Errorf("register failed: %v", ccErr)
				return util.Error(ccErr.Code(),
					fmt.Sprintf("register failed: %v", ccErr))
			}
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
		wallet, ccErr = sumTransaction(wallet, tx)
		if ccErr != nil {
			log.Errorf("wallet sum error: %s: %v", walletKey, err)
			return util.Error(ccErr.Code(),
				fmt.Sprintf("register failed: %v", ccErr))
		}
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
		log.Infof("register block wallet update: %s: %s", walletKey, string(bytes))
	}
	ret := &util.ChainResult{Code: 200, Message: "OK"}
	ret.Result = len(block.Txs)
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

func sumTransaction(wallet *util.Wallet, tx *TxRegister) (*util.Wallet, *ChaincodeError) {
	if wallet == nil || tx == nil ||
		!strings.EqualFold(wallet.Chain, tx.Chain) ||
		!strings.EqualFold(wallet.Token, tx.Token) {
		return nil, &ChaincodeError{
			code: http.StatusInternalServerError,
			errString: fmt.Sprintf("wrong wallet: %s, %s, %s/ %s, %s, %s",
				wallet.Chain, wallet.Token, wallet.Address,
				tx.Chain, tx.Token, tx.To)}
	}
	if !strings.HasPrefix(tx.Amount, "0x") && tx.Amount != "" {
		return nil, &ChaincodeError{
			code: http.StatusBadRequest,
			errString: fmt.Sprintf(
				`amount "%s" should be prefixed with 0x`, tx.Amount)}
	}
	if strings.EqualFold(wallet.Address, tx.To) &&
		!strings.EqualFold(wallet.Address, tx.From) {
		balance := new(big.Int)
		balance.SetString(wallet.Balance[2:], 16)
		if len(tx.Amount) > 0 {
			y := new(big.Int)
			y.SetString(tx.Amount[2:], 16)
			balance = balance.Add(balance, y)
		}
		wallet.Balance = fmt.Sprintf("0x%s", balance.Text(16))
		log.Info("wallet balance: ", wallet.Balance)
		return wallet, nil
	}
	if !strings.EqualFold(wallet.Address, tx.To) &&
		strings.EqualFold(wallet.Address, tx.From) {
		balance := new(big.Int)
		balance.SetString(wallet.Balance[2:], 16)
		if len(tx.Amount) > 0 {
			y := new(big.Int)
			y.SetString(tx.Amount[2:], 16)
			balance = balance.Sub(balance, y)
		}
		if len(tx.GasUsed) > 0 && len(tx.GasPrice) > 0 {
			g := new(big.Int)
			g.SetString(tx.GasUsed[2:], 16)
			gp := new(big.Int)
			gp.SetString(tx.GasPrice[2:], 16)
			g = g.Mul(g, gp)
			balance = balance.Sub(balance, g)
		}
		wallet.Balance = fmt.Sprintf("0x%s", balance.Text(16))
		log.Info("wallet balance: ", wallet.Balance)
		return wallet, nil
	}
	return nil, &ChaincodeError{
		code: http.StatusBadRequest,
		errString: fmt.Sprintf(
			`incorrect address, wallet: %s from: %s, to: %s`,
			wallet.Address, tx.From, tx.To)}
}

// construct a page struct from iterator
func constructPage(iterator shim.StateQueryIteratorInterface,
	metadata *pb.QueryResponseMetadata) (*paginationTxs, error) {
	var recs []interface{}
	var rec *queryresult.KV
	var err error
	for iterator.HasNext() {
		rec, err = iterator.Next()
		if err != nil {
			return nil, err
		}
		tx := &TxRegister{}
		err = json.Unmarshal(rec.Value, tx)
		if err != nil {
			return nil, err
		}
		tx.Key = rec.Key
		recs = append(recs, tx)
	}
	page := &paginationTxs{}
	page.Metadata = metadata
	page.Records = recs
	return page, nil
}
