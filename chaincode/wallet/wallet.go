/**
 *
 */
package main

import (
	"encoding/json"
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

const (
	version    = "Version"    // VERSION key of chaincode versio
	createtime = "Createtime" // CREATETIME key of init time of chaincode

	// ChaincodeVersion current version of chaincode
	ChaincodeVersion string = "0.0.1"
)

// ChaincodeError error of chaincode
type ChaincodeError struct {
	code      int
	errString string
}

// Code of error
func (e ChaincodeError) Code() int {
	return e.code
}

// Error string
func (e ChaincodeError) Error() string {
	return e.errString
}

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
		return w.register(stub, args)
	}

	return shim.Error(fmt.Sprintf("Unknown function call: %s", function))
}

// register function of the chaincode
func (w *WalletChaincode) register(
	stub shim.ChaincodeStubInterface, args []string) pb.Response {
	l := len(args)
	what := args[0]
	switch what {
	case "wallet":
		{
			// register a wallet
			if l < 5 {
				return util.Error(http.StatusBadRequest, fmt.Sprintf(
					"expecting at least 5 arguments, %d", l))
			}
			return w.registerWallet(stub, args[1:])
		}
	case "block":
		{
			// register chain's block
			return w.registerBlock(stub, args[1:])
		}
	case "transaction":
		{
			// register chain's transaction
			if l < 2 {
				return util.Error(http.StatusBadRequest, fmt.Sprintf(
					"expecting at least 2 arguments, %d", l))
			}
			return w.registerTransaction(stub, args[1:])
		}
	}
	return util.Error(http.StatusBadRequest, fmt.Sprintf(
		"register failed: Register what? %s", what))
}

// query function of the chaincode
func (w *WalletChaincode) query(
	stub shim.ChaincodeStubInterface, args []string) pb.Response {
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

func main() {
	err := shim.Start(new(WalletChaincode))
	if err != nil {
		fmt.Printf("Error starting WalletChaincode: %s", err)
	}
}

func checkState(stub shim.ChaincodeStubInterface,
	key string, returnExistError bool) ([]byte, *ChaincodeError) {
	bytes, err := stub.GetState(key)
	if err != nil {
		log.Errorf("check state error: %s: %v", key, err)
		ccErr := &ChaincodeError{
			code:      http.StatusInternalServerError,
			errString: fmt.Sprintf("check state error: %s: %v", key, err)}
		return nil, ccErr
	}
	if returnExistError && bytes != nil {
		log.Warnf("check state error: %s: state exist", key)
		ccErr := &ChaincodeError{
			code:      http.StatusConflict,
			errString: fmt.Sprintf("state exist: %s", key)}
		return bytes, ccErr
	}
	return bytes, nil
}
