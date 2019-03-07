package main

import "C"
import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/snowdiceX/dawns.world/cazimi/fabric"
	"github.com/snowdiceX/dawns.world/cazimi/jni"
	"github.com/snowdiceX/dawns.world/cazimi/log"
)

// ----------------------------------------------------------------------------
// source code for so file generation with "go build " command,
// e.g.
// go build -o cazimi.so -buildmode=c-shared cazimi.go
// ----------------------------------------------------------------------------

// InitJNI init jni
func InitJNI() {
	log.Info("init jni...")
}

// ChaincodeInvoke call chaincode invoke of hyperledger fabric
func ChaincodeInvoke(chainID, chaincodeID, argsStr string) string {
	log.Info("chaincode invoke...")
	result := &jni.CallResult{
		Code: http.StatusOK, Message: "OK"}
	argsArray, err := fabric.ArgsArray(argsStr)
	if err == nil {
		for _, args := range argsArray {
			result.Message = fmt.Sprintf("OK; %s:%s:%s", chainID, chaincodeID, args.Func)
		}
		ret, err := fabric.ChaincodeInvoke(chainID, chaincodeID, argsArray)
		if err != nil {
			result.Code = http.StatusInternalServerError
			result.Message = fmt.Sprintf("chaincode invoke error: %v", err)
		} else {
			log.Info("chaincode invoke result: ", ret)
		}
	} else {
		result.Code = http.StatusInternalServerError
		result.Message = fmt.Sprintf("args JSON parsing err: %v", err)
	}
	bytes, err := json.Marshal(result)
	if err == nil {
		log.Info(string(bytes))
		return string(bytes)
	}
	log.Errorf("%s %v", jni.DefaultResultJSON, err)
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
