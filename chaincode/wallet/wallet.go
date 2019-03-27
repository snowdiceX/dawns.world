/**
 *
 */
package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync/atomic"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

const (
	version    = "Version"    // VERSION key of chaincode versio
	createtime = "Createtime" // CREATETIME key of init time of chaincode
	tagWallet  = "Wallet"     // WALLET key prefix of wallet address
	sequence   = "Sequence"   // SEQUENCE key prefix of transaction sequence

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
	if function == "create" {
		// create a wallet
		return w.create(stub, args)
	}

	if function == "query" {
		// queries an entity state
		return w.query(stub, args)
	}

	if function == "register" {
		// register a wallet
		return w.register(stub, args)
	}

	return shim.Error(fmt.Sprintf("Unknown function call: %s", function))
}

/**
* args:
      0 network
      1 token name
      2 address
      3 height
      4 tx id
      5 token amount
* account key: Wallet_[address]
* wallet key: [network]+[token name]+[address]
*/
func (w *WalletChaincode) create(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) < 6 {
		return shim.Error(fmt.Sprintf("expecting at least 6, %d", len(args)))
	}
	address := args[2]
	accountKey := buildAccountKey("asd", address)
	if err := stub.PutState(accountKey, []byte(address)); err != nil {
		return shim.Error(fmt.Sprintf("Error putting data for key [%s]: %s", accountKey, err))
	}
	fmt.Println("create an account: ", accountKey)

	walletKey := buildWalletKey(args[0], args[1], address)
	if err := stub.PutState(walletKey, []byte(args[5])); err != nil {
		return shim.Error(fmt.Sprintf("Error putting data for key [%s]: %s", walletKey, err))
	}
	fmt.Println("create a wallet: ", walletKey)

	// seqBytes, err := stub.GetState(SEQUENCE)
	// if err != nil {
	// 	return shim.Error("Failed to get state")
	// }
	// if seqBytes == nil {
	// 	return shim.Error("Entity not found")
	// }
	// seq, _ := strconv.ParseInt(string(seqBytes), 10, 64)
	seq := atomic.AddUint64(&w.OutSequence, 1)
	sequenceKey := buildSequenceKey(seq)
	jsonTx := "{\"sequence\":\"" + strconv.FormatUint(seq, 10) + "\",\"txid\":\"" + string(stub.GetTxID()) + "\"}"
	if err := stub.PutState(sequenceKey, []byte(jsonTx)); err != nil {
		return shim.Error(fmt.Sprintf("Error putting data for key [%s]: %s", walletKey, err))
	}

	fmt.Println("create success: ", stub.GetTxID())
	return shim.Success([]byte(fmt.Sprintf("{\"wallet\":\"%s\", \"txid\":\"%s\"}", walletKey, stub.GetTxID())))
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
		return w.querySequence(stub, args)
	}
	if strings.EqualFold("transaction", args[0]) {
		// queries a transaction by sequence
		return w.queryTransactionBySequence(stub, args)
	}

	return shim.Error(fmt.Sprintf("Unknown query function call: %s", args[0]))
}

func (w *WalletChaincode) querySequence(
	stub shim.ChaincodeStubInterface, args []string) pb.Response {
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
	sequenceKey := buildSequenceKey(seq)

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
      1 key
      2 network
      3 token name
* register key: Register_[userID]_[network]_[token name]
* wallet key: [network]-[token name]-[address]
*/
func (w *WalletChaincode) register(
	stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) < 6 {
		return shim.Error(fmt.Sprintf("expecting at least 6 arguments, %d", len(args)))
	}
	accountID, address, key, chain, token, heightHex :=
		args[0], args[1], args[2], args[3], args[4], args[5]

	fmt.Printf(`chaincode[wallet] register
		account: %s,
		address: %s,
		key: %s,
		chain: %s,
		token: %s,
		height: %s`+"\n",
		accountID, address, key, chain, token, heightHex)
	var height uint64
	var err error
	if height, err = w.checkInSequence(stub, heightHex); err != nil {
		fmt.Println("in sequence check error: ", err.Error())
		return shim.Error(err.Error())
	}
	seq := atomic.AddUint64(&w.OutSequence, 1)
	sequenceKey := buildSequenceKey(seq)
	jsonTx := fmt.Sprintf(`{"code": 0, "message": "OK", "sequence": %s`+
		`, "txid":"%s", "func":"register", "address":"%s"`+
		`, "network":"%s", "token":"%s", "height": %d}`,
		strconv.FormatUint(seq, 10),
		string(stub.GetTxID()), address,
		chain, token, height)
	if err := stub.PutState(sequenceKey, []byte("v1.0.0:"+jsonTx)); err != nil {
		return shim.Error(fmt.Sprintf("Error putting data for key [%s]: %s", sequenceKey, err))
	}
	jsonWallet := `v1.0.0:{"chain":"` + chain +
		`", "token":"` + token +
		`", "height":"` + heightHex +
		`", "accountID":"` + accountID +
		`", "address":"` + address +
		`", "key":"` + key +
		`", "txid":"` + string(stub.GetTxID()) +
		`", "agent":"node?ip:port?agentPubKey?"}`
	accountKey := buildAccountKey(args[0], args[1])
	fmt.Printf(`chaincode[wallet] register: %s, data: %s`+"\n", accountKey, jsonWallet)
	if err := stub.PutState(accountKey, []byte(jsonWallet)); err != nil {
		return shim.Error(fmt.Sprintf("Error putting data for key [%s]: %s", sequenceKey, err))
	}
	return shim.Success([]byte(jsonTx))
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

func buildAccountKey(accountID, address string) string {
	return fmt.Sprintf("%s-%s-%s", tagWallet, accountID, address)
}

func buildWalletKey(network, token, address string) string {
	return fmt.Sprintf("%s-%s-%s", network, token, address)
}

func buildSequenceKey(seq uint64) string {
	return fmt.Sprintf("%s-%d", sequence, seq)
}

func main() {
	err := shim.Start(new(WalletChaincode))
	if err != nil {
		fmt.Printf("Error starting WalletChaincode: %s", err)
	}
}
