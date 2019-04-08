/**
 *
 */
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync/atomic"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/snowdiceX/dawns.world/chaincode/log"
	"github.com/snowdiceX/dawns.world/chaincode/util"
)

const (
	version    = "Version"    // VERSION key of chaincode versio
	createtime = "Createtime" // CREATETIME key of init time of chaincode

	// ChaincodeVersion current version of chaincode
	ChaincodeVersion string = "0.0.1"
)

// WalletChaincode is wallet Chaincode implementation
type WalletChaincode struct {
	Createtime  string
	OutSequence uint64
	InSequence  uint64
}

// Init ...
func (w *WalletChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("WalletChaincode Init...")
	// args := stub.GetStringArgs()

	var err error

	// Initialize the chaincode
	// fmt.Println("init: ", args...)
	// var sequence int64

	// Write the state to the ledger
	err = stub.PutState(version, []byte(ChaincodeVersion))
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(createtime, []byte(w.Createtime))
	if err != nil {
		return shim.Error(err.Error())
	}

	// err = stub.PutState(SEQUENCE, []byte(strconv.FormatInt(sequence, 10)))
	// if err != nil {
	// 	return shim.Error(err.Error())
	// }

	if transientMap, err := stub.GetTransient(); err == nil {
		if transientData, ok := transientMap["result"]; ok {
			fmt.Printf("Transient data in 'init' : %s\n", transientData)
			return shim.Success(transientData)
		}
	}
	return shim.Success(nil)

}

// Query ...
func (w *WalletChaincode) Query(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Error("Unknown supported call")
}

// Invoke ...
func (w *WalletChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	// args := stub.GetStringArgs()
	//
	// if len(args) == 0 {
	// 	return shim.Error("Function not provided")
	// }
	//
	// function := args[0]
	// args = args[1:]

	function, args := stub.GetFunctionAndParameters()
	fmt.Println("WalletChaincode Invoke: ", function)
	if function == "query" {
		// queries an entity state
		return w.query(stub, args)
	}
	if function == "register" {
		// register a wallet
		return w.register(stub, args)
	}
	if function == "registerBlock" {
		// register chain's block
		return w.registerBlock(stub, args)
	}

	return shim.Error(fmt.Sprintf("Unknown function call: %s", function))
}

// Query callback representing the query of a chaincode
func (w *WalletChaincode) query(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// if len(args) != 3 {
	// 	return shim.Error(fmt.Sprintf("Incorrect number of arguments: %v", args))
	// }
	// walletKey := buildWalletKey(args[0], args[1], args[2])

	// // Get the state from the ledger
	// walletBytes, err := stub.GetState(walletKey)
	// if err != nil {
	// 	jsonResp := "{\"Error\":\"Failed to get state for " + walletKey + "; " +
	// 		strconv.FormatUint(w.Sequence, 10) + "\"}"
	// 	return shim.Error(jsonResp)
	// }

	// if walletBytes == nil {
	// 	jsonResp := "{\"Error\":\"Nil amount for " + walletKey + "; " +
	// 		strconv.FormatUint(w.Sequence, 10) + "\"}"
	// 	return shim.Error(jsonResp)
	// }

	// jsonResp := "{\"wallet\":\"" + walletKey + "; " +
	// 	strconv.FormatUint(w.Sequence, 10) + "\",\"amount\":\"" + string(walletBytes) + "\"}"

	if len(args) == 0 || strings.EqualFold("sequence", args[0]) {
		return w.querySequence(stub)
	}
	if strings.EqualFold("transaction", args[0]) {
		// queries a transaction by sequence
		return w.queryTransactionBySequence(stub, args)
	}
	if strings.EqualFold("wallet", args[0]) {
		// queries a transaction by sequence
		return w.queryWallet(stub, args[1], args[2], args[3])
	}

	return shim.Error(fmt.Sprintf("Unknown query function call: %s", args[0]))
}

func (w *WalletChaincode) queryWallet(stub shim.ChaincodeStubInterface,
	chain, token, address string) pb.Response {
	walletKey := util.BuildWalletKey(chain, token, address)
	bytes, err := stub.GetState(walletKey)
	if err != nil {
		log.Debug("query wallet error: ", err.Error())
		return util.Error(500, fmt.Sprintf("query wallet failed: %v", err))
	}
	log.Debug("query wallet: ", walletKey, "; ", string(bytes))
	wallet := &util.Wallet{}
	err = json.Unmarshal(bytes, wallet)
	if err != nil {
		log.Error("query wallet error: ", err.Error())
		return util.Error(500, fmt.Sprintf("query wallet failed: %v", err))
	}
	ret := &util.ChainResult{Code: 200, Message: "OK"}
	ret.Result = wallet
	return util.Success(ret)
}

func (w *WalletChaincode) querySequence(
	stub shim.ChaincodeStubInterface) pb.Response {
	in := atomic.LoadUint64(&w.InSequence)
	out := atomic.LoadUint64(&w.OutSequence)
	respJSON := fmt.Sprintf(
		`{"chaincode": "wallet", "InSequence": %d, "OutSequence": %d}`,
		in, out)
	fmt.Printf("Query Response:%s\n", respJSON)
	return shim.Success([]byte(respJSON))
}

// queryTransactionBySequence queries a transaction by sequence
func (w *WalletChaincode) queryTransactionBySequence(
	stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 3 {
		return shim.Error(fmt.Sprintf(
			"incorrect number of arguments: %v", args))
	}
	seq, _ := strconv.ParseUint(args[2], 10, 64)
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
func (w *WalletChaincode) register(
	stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) < 5 {
		return util.Error(500, fmt.Sprintf(
			"expecting at least 5 arguments, %d", len(args)))
	}
	_, address, chain, token, heightHex :=
		args[0], args[1], args[2], args[3], args[4]
	log.Debug("chaincode[wallet] register", fmt.Sprint(args))
	var err error
	if _, err = w.checkInSequence(stub, heightHex); err != nil {
		log.Debug("in sequence check error: ", err.Error())
		return util.Error(500, fmt.Sprintf("register failed: %v", err))
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
		return util.Error(500, fmt.Sprintf("register failed: %v", err))
	}
	sequenceKey := util.BuildSequenceKey(seq)
	log.Debug(`chaincode[wallet] sequence: %s, data: %s`+"\n",
		sequenceKey, string(bytes))
	if err = stub.PutState(sequenceKey, bytes); err != nil {
		return util.Error(500, fmt.Sprintf("register failed: %v", err))
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
		return util.Error(500, fmt.Sprintf("register failed: %v", err))
	}
	walletKey := util.BuildWalletKey(chain, token, address)
	log.Debug(`chaincode[wallet] register: %s, data: %s`+"\n",
		walletKey, string(bytes))
	if err = stub.PutState(walletKey, bytes); err != nil {
		return util.Error(500, fmt.Sprintf("register failed: %v", err))
	}
	ret := &util.ChainResult{Code: 200, Message: "OK"}
	ret.Result = walletSequence
	return util.Success(ret)
}

/**
* args:
      0 blockchain network
      1 block data
*/
func (w *WalletChaincode) registerBlock(
	stub shim.ChaincodeStubInterface, args []string) pb.Response {
	network, blockdata := args[0], args[1]
	return registerChainBlock(network, blockdata)
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

func main() {
	err := shim.Start(new(WalletChaincode))
	if err != nil {
		fmt.Printf("Error starting WalletChaincode: %s", err)
	}
}
