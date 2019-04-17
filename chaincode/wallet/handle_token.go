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

type registerToken struct {
	Key      string `json:"key,omitempty"`
	Contract string `json:"contract,omitempty"`
	Chain    string `json:"chain,omitempty"`
	Symbol   string `json:"symbol,omitempty"`
	Name     string `json:"name,omitempty"`
	Decimals string `json:"decimals,omitempty"`
}

func (w *WalletChaincode) registerToken(
	stub shim.ChaincodeStubInterface, args []string) pb.Response {
	contract, chain, symbol, name, decimals :=
		args[0], args[1], args[2], args[3], args[4]
	chain = strings.ToUpper(chain)
	symbol = strings.ToUpper(symbol)
	key := util.BuildTokenKey(chain, symbol)
	_, ccErr := checkState(stub, key, true)
	if ccErr != nil {
		log.Errorf("check state error: %s %v", key, ccErr)
		return util.Error(ccErr.Code(),
			fmt.Sprintf("check token state error: %s: %v", key, ccErr))
	}
	token := &registerToken{
		Contract: contract,
		Chain:    chain,
		Symbol:   symbol,
		Name:     name,
		Decimals: decimals}
	var bytes []byte
	var err error
	if bytes, err = json.Marshal(token); err != nil {
		if err != nil {
			log.Errorf("register token json marshal error: %s %s", chain, symbol)
			return util.Error(http.StatusInternalServerError, err.Error())
		}
	}
	if err = stub.PutState(key, bytes); err != nil {
		log.Errorf("register token put state error: %s %s", chain, symbol)
		return util.Error(http.StatusInternalServerError, err.Error())
	}
	ret := &util.ChainResult{Code: 200, Message: "OK"}
	ret.Result = token
	return util.Success(ret)
}

func (w *WalletChaincode) queryToken(
	stub shim.ChaincodeStubInterface, args []string) pb.Response {
	what := args[0]
	switch what {
	case "page":
		{
			// query tokens
			return w.paginationTokens(stub, args[1:])
		}
	}
	return util.Error(http.StatusBadRequest, fmt.Sprintf(
		"query token failed: Query what? %s", what))
}

func (w *WalletChaincode) paginationTokens(
	stub shim.ChaincodeStubInterface, args []string) pb.Response {
	_, pageSizeHex := args[0], args[1]
	pageSize, _ := strconv.ParseInt(pageSizeHex, 10, 32)
	startKey := util.BuildTokenStartKey()
	logIterator, meta, err := stub.GetStateByRangeWithPagination(
		startKey, startKey+"a", int32(pageSize), "")
	if err != nil {
		log.Errorf("paging tokens error: %v", err)
		return util.Error(http.StatusInternalServerError,
			fmt.Sprintf("paging tokens failed: %v", err))
	}
	defer logIterator.Close()
	page, err := constructTokenPage(logIterator, meta)
	if err != nil {
		log.Errorf("paging tokens error: %v", err)
		return util.Error(http.StatusInternalServerError,
			fmt.Sprintf("paging tokens failed: %v", err))
	}
	ret := &util.ChainResult{Code: 200, Message: "OK"}
	ret.Result = page
	return util.Success(ret)
}

// construct a page struct from iterator
func constructTokenPage(iterator shim.StateQueryIteratorInterface,
	metadata *pb.QueryResponseMetadata) (*paginationTxs, error) {
	var recs []interface{}
	var rec *queryresult.KV
	var err error
	for iterator.HasNext() {
		rec, err = iterator.Next()
		if err != nil {
			return nil, err
		}
		rt := &registerToken{}
		err = json.Unmarshal(rec.Value, rt)
		if err != nil {
			return nil, err
		}
		rt.Key = rec.Key
		recs = append(recs, rt)
	}
	page := &paginationTxs{}
	page.Metadata = metadata
	page.Records = recs
	return page, nil
}
