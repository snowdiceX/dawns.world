package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync/atomic"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/snowdiceX/dawns.world/chaincode/log"
	"github.com/snowdiceX/dawns.world/chaincode/util"
)

/**
* args:
      0 accountID
      1 address
      2 blockchain network
      3 token name
      4 blockchain height in hex
* register key: Register_[userID]_[network]_[token name]
* wallet key: [network]-[token name]-[address]
*/
func (w *WalletChaincode) registerWallet(
	stub shim.ChaincodeStubInterface, args []string) pb.Response {
	_, address, chain, token, heightHex :=
		args[0], args[1], args[2], args[3], args[4]
	log.Debug("chaincode[wallet] register", fmt.Sprint(args))
	var err error
	if _, err = w.checkInSequence(stub, heightHex); err != nil {
		log.Errorf("in sequence check error: %v", err)
		return util.Error(http.StatusBadRequest, fmt.Sprintf(
			"register failed: %v", err))
	}
	seq := atomic.AddUint64(&w.OutSequence, 1)
	walletSequence := util.WalletSequence{
		Version:  "v1.0.0",
		Sequence: strconv.FormatUint(seq, 10),
		TxID:     string(stub.GetTxID()),
		Func:     "register",
		Address:  address,
		Network:  chain,
		Token:    token,
		Height:   heightHex}
	var bytes []byte
	bytes, err = json.Marshal(walletSequence)
	if err != nil {
		return util.Error(http.StatusInternalServerError,
			fmt.Sprintf("register failed: %v", err))
	}
	sequenceKey := util.BuildSequenceKey(seq)
	log.Debug(`chaincode[wallet] sequence: %s, data: %s`+"\n",
		sequenceKey, string(bytes))
	if err = stub.PutState(sequenceKey, bytes); err != nil {
		return util.Error(http.StatusInternalServerError,
			fmt.Sprintf("register failed: %v", err))
	}
	wallet := util.Wallet{
		Version: "v1.0.0",
		Chain:   chain,
		Token:   token,
		Balance: "0x0",
		Height:  heightHex,
		Address: address,
		TxID:    string(stub.GetTxID()),
		Agent:   "peer0.org0"}
	bytes, err = json.Marshal(wallet)
	if err != nil {
		return util.Error(http.StatusInternalServerError,
			fmt.Sprintf("register failed: %v", err))
	}
	walletKey := util.BuildWalletKey(chain, token, address)
	log.Debug(`chaincode[wallet] register: %s, data: %s`+"\n",
		walletKey, string(bytes))
	if err = stub.PutState(walletKey, bytes); err != nil {
		return util.Error(http.StatusInternalServerError,
			fmt.Sprintf("register failed: %v", err))
	}
	ret := &util.ChainResult{Code: 200, Message: "OK"}
	ret.Result = walletSequence
	return util.Success(ret)
}

func (w *WalletChaincode) checkInSequence(
	stub shim.ChaincodeStubInterface, height string) (h uint64, err error) {
	if strings.HasPrefix(height, "0x") {
		height = height[2:]
	}
	h, err = strconv.ParseUint(height, 16, 64)
	if err != nil {
		msg := fmt.Sprintf("parse height 0x%s error: %v",
			height, err)
		err = errors.New(msg)
		return
	}
	if w.InSequence <= 0 {
		inSeq := atomic.LoadUint64(&w.InSequence)
		if inSeq > 0 {
			return
		}
		if atomic.CompareAndSwapUint64(&w.InSequence, 0, h) {
			if err = stub.PutState("Wallet-InSequence",
				[]byte(height)); err != nil {
				atomic.StoreUint64(&w.InSequence, 0)
				err = errors.New("update in sequence error")
				return
			}
		}
	}
	return
}
