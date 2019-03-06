package main

import "C"
import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/securekey/fabric-examples/fabric-cli/cmd/fabric-cli/action"
	"github.com/snowdiceX/dawns.world/cazimi/jni"
)

// ----------------------------------------------------------------------------
// source code for so file generation with "go build " command,
// e.g.
// go build -o cazimi.so -buildmode=c-shared cazimi.go
// ----------------------------------------------------------------------------

// InitJNI init jni
func InitJNI() {
	fmt.Println("init jni...")
}

// ChaincodeInvoke call chaincode invoke of hyperledger fabric
func ChaincodeInvoke(chainID, chaincodeID, argsStr string) string {
	fmt.Println("chaincode invoke...")
	result := &jni.CallResult{
		Code: http.StatusOK, Message: "OK"}
	args := &action.ArgStruct{}
	err := json.Unmarshal([]byte(argsStr), args)
	if err == nil {

		result.Message = fmt.Sprintf("OK; %s:%s:%s", chainID, chaincodeID, args.Func)
	} else {
		result.Code = http.StatusInternalServerError
		result.Message = fmt.Sprintf("chaincode invoke err: %v", err)
	}
	bytes, err := json.Marshal(result)
	if err == nil {
		return string(bytes)
	}
	return jni.DefaultResultJSON
}

//export initJNI
func initJNI() {
	InitJNI()
}

//export chaincodeInvoke
func chaincodeInvoke(chainID, chaincodeID, args *C.char) *C.char {
	return C.CString(ChaincodeInvoke(C.GoString(chainID),
		C.GoString(chaincodeID), C.GoString(args)))
}

func main() {
	//InitJNI()
}
