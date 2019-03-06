package main

import "C"
import (
	"fmt"
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
func ChaincodeInvoke(chainID, chaincodeID, function, args string) string {
	fmt.Println("chaincode invoke...")
	return fmt.Sprintf("%s:%s:%s", chainID, chaincodeID, function)
}

//export initJNI
func initJNI() {
	InitJNI()
}

//export chaincodeInvoke
func chaincodeInvoke(chainID, chaincodeID, function, args *C.char) *C.char {
	return C.CString(ChaincodeInvoke(C.GoString(chainID),
		C.GoString(chaincodeID), C.GoString(function), C.GoString(args)))
}

func main() {
	//InitJNI()
}
