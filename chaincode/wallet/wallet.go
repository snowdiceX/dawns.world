/**
 *
 */
package main

import (
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

const (
	VERSION    = "Version"
	CREATETIME = "Createtime"
)

// WalletChaincode is wallet Chaincode implementation
type WalletChaincode struct {
	Version    string
	Createtime string
}

// Init ...
func (w *WalletChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("WalletChaincode Init...")
	args := stub.GetStringArgs()

	var err error

	// Initialize the chaincode
	fmt.Printf("init: ", args...)

	// Write the state to the ledger
	err = stub.PutState(VERSION, []byte(w.Version))
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(CREATETIME, []byte(w.Createtime))
	if err != nil {
		return shim.Error(err.Error())
	}

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
// Transaction makes payment of X units from A to B
func (w *WalletChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("WalletChaincode Invoke...")
	args := stub.GetStringArgs()
	if len(args) == 0 {
		return shim.Error("Function not provided")
	}

	function := args[0]

	if function == "create" {
		// create a wallet
		return w.create(stub, args)
	}

	if function == "query" {
		// queries an entity state
		return w.query(stub, args)
	}

	if function == "move" {
		eventID := "testEvent"
		if len(args) >= 5 {
			eventID = args[4]
		}
		if err := stub.SetEvent(eventID, []byte("Test Payload")); err != nil {
			return shim.Error("Unable to set CC event: testEvent. Aborting transaction ...")
		}
		return w.move(stub, args)
	}

	if function == "put" {
		return w.put(stub, args[1:])
	}

	if function == "get" {
		return w.get(stub, args[1:])
	}

	return shim.Error(fmt.Sprintf("Unknown function call: %s", function))
}

/**
 * wallet address: [network]+[address]+[token name]
 */
func (w *WalletChaincode) create(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) < 5 {
		return shim.Error("Incorrect number of arguments. Expecting at least 2")
	}
	fmt.Sprintf()
}

func (w *WalletChaincode) move(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// must be an invoke
	var A, B string    // Entities
	var Aval, Bval int // Asset holdings
	var X int          // Transaction value
	var err error
	if len(args) < 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4, function followed by 2 names and 1 value")
	}

	A = args[1]
	B = args[2]

	// Get the state from the ledger
	// TODO: will be nice to have a GetAllState call to ledger
	Avalbytes, err := stub.GetState(A)
	if err != nil {
		return shim.Error("Failed to get state")
	}
	if Avalbytes == nil {
		return shim.Error("Entity not found")
	}
	Aval, _ = strconv.Atoi(string(Avalbytes))

	Bvalbytes, err := stub.GetState(B)
	if err != nil {
		return shim.Error("Failed to get state")
	}
	if Bvalbytes == nil {
		return shim.Error("Entity not found")
	}
	Bval, _ = strconv.Atoi(string(Bvalbytes))

	// Perform the execution
	X, err = strconv.Atoi(args[3])
	if err != nil {
		return shim.Error("Invalid transaction amount, expecting a integer value")
	}
	Aval = Aval - X
	Bval = Bval + X
	fmt.Printf("Aval = %d, Bval = %d\n", Aval, Bval)

	// Write the state back to the ledger
	err = stub.PutState(A, []byte(strconv.Itoa(Aval)))
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(B, []byte(strconv.Itoa(Bval)))
	if err != nil {
		return shim.Error(err.Error())
	}

	if transientMap, err := stub.GetTransient(); err == nil {
		if transientData, ok := transientMap["result"]; ok {
			fmt.Printf("Transient data in 'move' : %s\n", transientData)
			return shim.Success(transientData)
		}
	}
	return shim.Success(nil)
}

// Deletes an entity from state
func (w *WalletChaincode) delete(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	A := args[1]

	// Delete the key from the state in ledger
	err := stub.DelState(A)
	if err != nil {
		return shim.Error("Failed to delete state")
	}

	return shim.Success(nil)
}

// Query callback representing the query of a chaincode
func (w *WalletChaincode) query(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var A string // Entities
	var err error

	if len(args) != 2 {
		return shim.Error(fmt.Sprintf("Incorrect number of arguments: %v", args))
	}

	A = args[1]

	// Get the state from the ledger
	Avalbytes, err := stub.GetState(A)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + A + "\"}"
		return shim.Error(jsonResp)
	}

	if Avalbytes == nil {
		jsonResp := "{\"Error\":\"Nil amount for " + A + "\"}"
		return shim.Error(jsonResp)
	}

	jsonResp := "{\"Name\":\"" + A + "\",\"Amount\":\"" + string(Avalbytes) + "\"}"
	fmt.Printf("Query Response:%s\n", jsonResp)
	return shim.Success(Avalbytes)
}

func (w *WalletChaincode) put(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return shim.Error("Invalid args. Expecting key and value")
	}

	key := args[0]
	value := args[1]

	existingValue, err := stub.GetState(key)
	if err != nil {
		return shim.Error(fmt.Sprintf("Error getting data for key [%s]: %s", key, err))
	}
	if existingValue != nil {
		value = string(existingValue) + "-" + value
	}

	if err := stub.PutState(key, []byte(value)); err != nil {
		return shim.Error(fmt.Sprintf("Error putting data for key [%s]: %s", key, err))
	}

	return shim.Success([]byte(value))
}

func (w *WalletChaincode) get(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Invalid args. Expecting key")
	}

	key := args[0]

	value, err := stub.GetState(key)
	if err != nil {
		return shim.Error(fmt.Sprintf("Error getting data for key [%s]: %s", key, err))
	}

	return shim.Success([]byte(value))
}

func main() {
	err := shim.Start(new(WalletChaincode))
	if err != nil {
		fmt.Printf("Error starting WalletChaincode: %s", err)
	}
}
