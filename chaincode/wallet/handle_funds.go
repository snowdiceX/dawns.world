package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/ledger/queryresult"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/snowdiceX/dawns.world/chaincode/wallet/log"
	"github.com/snowdiceX/dawns.world/chaincode/wallet/util"
)

// TokenBalance records balance of token
type TokenBalance struct {
	FundsHash string `json:"funds_hash,omitempty"`
	Chain     string `json:"chain,omitempty"`
	Token     string `json:"token,omitempty"`
	Balance   string `json:"balance,omitempty"`
}

// fundsState store in hyperledger fabric chain
type fundsState struct {
	Key       string `json:"key,omitempty"`
	Version   string `json:"version,omitempty"`
	Hash      string `json:"hash,omitempty"`
	BaseKey   string `json:"base,omitempty"`
	AcceptKey string `json:"accept,omitempty"`
}

// Funds with base and accept tokens
type Funds struct {
	Key     string       `json:"key,omitempty"`
	Version string       `json:"version,omitempty"`
	Hash    string       `json:"hash,omitempty"`
	Base    TokenBalance `json:"base,omitempty"`
	Accept  TokenBalance `json:"accept,omitempty"`
}

func (w *WalletChaincode) registerFunds(
	stub shim.ChaincodeStubInterface, args []string) pb.Response {
	baseChain, baseToken, chain, token :=
		args[0], args[1], args[2], args[3]
	baseChain = strings.ToUpper(baseChain)
	baseToken = strings.ToUpper(baseToken)
	chain = strings.ToUpper(chain)
	token = strings.ToUpper(token)
	key := util.BuildFundsKey(baseChain, baseToken, chain, token)
	_, ccErr := util.CheckState(stub, key, true)
	if ccErr != nil {
		log.Errorf("check state error: %s %v", key, ccErr)
		return util.Error(ccErr.Code,
			fmt.Sprintf("check state error: %s %v", key, ccErr))
	}
	fundsHash := strings.ToUpper(util.Hash(key))
	fundsState := &fundsState{
		Version:   util.ChaincodeVersion,
		Hash:      fundsHash,
		BaseKey:   util.BuildFundsBaseKey(baseChain, baseToken, fundsHash),
		AcceptKey: util.BuildFundsBaseKey(chain, token, fundsHash)}
	var bytes []byte
	var err error
	if bytes, err = json.Marshal(fundsState); err != nil {
		if err != nil {
			log.Errorf("register funds json marshal error: %s %v", key, err)
			return util.Error(http.StatusInternalServerError, err.Error())
		}
	}
	if err = stub.PutState(key, bytes); err != nil {
		log.Errorf("register funds put state error: %s %v", key, err)
		return util.Error(http.StatusInternalServerError, err.Error())
	}
	base := &TokenBalance{
		Chain: baseChain, Token: baseToken, Balance: util.ZeroBalance}
	if bytes, err = json.Marshal(base); err != nil {
		if err != nil {
			log.Errorf("register funds base json marshal error: %s %v",
				fundsState.BaseKey, err)
			return util.Error(http.StatusInternalServerError, err.Error())
		}
	}
	if err = stub.PutState(fundsState.BaseKey, bytes); err != nil {
		log.Errorf("register funds base put state error: %s %v",
			fundsState.BaseKey, err)
		return util.Error(http.StatusInternalServerError, err.Error())
	}
	accept := &TokenBalance{
		Chain: chain, Token: token, Balance: util.ZeroBalance}
	if bytes, err = json.Marshal(accept); err != nil {
		if err != nil {
			log.Errorf("register funds accept json marshal error: %s %v",
				fundsState.AcceptKey, err)
			return util.Error(http.StatusInternalServerError, err.Error())
		}
	}
	if err = stub.PutState(fundsState.AcceptKey, bytes); err != nil {
		log.Errorf("register funds accept put state error: %s %v",
			fundsState.AcceptKey, err)
		return util.Error(http.StatusInternalServerError, err.Error())
	}
	funds := &Funds{
		Version: fundsState.Version,
		Hash:    fundsState.Hash,
		Base:    *base,
		Accept:  *accept}
	ret := &util.ChainResult{Code: 200, Message: "OK"}
	ret.Result = funds
	return util.Success(ret)
}

func (w *WalletChaincode) queryFunds(
	stub shim.ChaincodeStubInterface, args []string) pb.Response {
	what := args[0]
	switch what {
	case "page":
		{
			// query tokens
			return w.paginationFunds(stub, args[1:])
		}
	}
	return util.Error(http.StatusBadRequest, fmt.Sprintf(
		"query funds failed: query what? %s", what))
}

func (w *WalletChaincode) funds(
	stub shim.ChaincodeStubInterface, args []string) pb.Response {
	what := args[0]
	switch what {
	case "deposit":
		{
			return w.fundsDeposit(stub, args[1:])
		}
	case "withdraw":
		{
			return w.fundsWithdraw(stub, args[1:])
		}
	}
	return util.Error(http.StatusBadRequest, fmt.Sprintf(
		"funds execute failed: execute what? %s", what))
}

func (w *WalletChaincode) fundsDeposit(
	stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fundsTokenKey, walletAddress, amount := args[0], args[1], args[2]
	rec := util.NewRecordFunds(fundsTokenKey, walletAddress)
	var err *util.ChaincodeError
	if err = rec.Load(stub); err != nil {
		return util.Error(err.Code, err.Error())
	}
	wallet := util.NewWallet(
		rec.Chain,
		rec.Token,
		rec.WalletAddress)
	if err := wallet.Load(stub); err != nil {
		return util.Error(err.Code, err.Error())
	}
	if err = wallet.Sub(amount); err != nil {
		return util.Error(err.Code, err.Error())
	}
	if err = wallet.Save(stub); err != nil {
		return util.Error(err.Code, err.Error())
	}
	ret := &util.ChainResult{Code: 200, Message: "OK"}
	ret.Result = rec
	return util.Success(ret)
}

func (w *WalletChaincode) fundsWithdraw(
	stub shim.ChaincodeStubInterface, args []string) pb.Response {
	rec := util.NewRecordFunds(args[0], args[1], args[2])
	if err := rec.Load(stub); err != nil {
		return util.Error(err.Code, err.Error())
	}
	ret := &util.ChainResult{Code: 200, Message: "OK"}
	return util.Success(ret)
}

func (w *WalletChaincode) paginationFunds(
	stub shim.ChaincodeStubInterface, args []string) pb.Response {
	_, pageSizeHex := args[0], args[1]
	pageSize, _ := strconv.ParseInt(pageSizeHex, 10, 32)
	startKey := util.BuildFundsStartKey()
	logIterator, meta, err := stub.GetStateByRangeWithPagination(
		startKey, startKey+"a", int32(pageSize), "")
	if err != nil {
		log.Errorf("paging funds error: %v", err)
		return util.Error(http.StatusInternalServerError,
			fmt.Sprintf("paging funds failed: %v", err))
	}
	defer logIterator.Close()
	page, err := constructFundsPage(logIterator, meta)
	if err != nil {
		log.Errorf("paging funds error: %v", err)
		return util.Error(http.StatusInternalServerError,
			fmt.Sprintf("paging funds failed: %v", err))
	}
	ret := &util.ChainResult{Code: 200, Message: "OK"}
	ret.Result = page
	return util.Success(ret)
}

// construct a page struct from iterator
func constructFundsPage(iterator shim.StateQueryIteratorInterface,
	metadata *pb.QueryResponseMetadata) (*util.Pagination, error) {
	var recs []interface{}
	var rec *queryresult.KV
	var err error
	for iterator.HasNext() {
		rec, err = iterator.Next()
		if err != nil {
			return nil, err
		}
		log.Debug(string(rec.Value))
		rt := &fundsState{}
		err = json.Unmarshal(rec.Value, rt)
		if err != nil {
			return nil, err
		}
		rt.Key = rec.Key
		recs = append(recs, rt)
	}
	page := &util.Pagination{}
	page.Metadata = metadata
	page.Records = recs
	return page, nil
}
