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

// RegisterWallet call register chaincode of hyperledger fabric
func RegisterWallet(account, token, network string) string {
	fmt.Println("register wallet...")
	return fmt.Sprintf("%s:%s", network, token)
}

//export initJNI
func initJNI() {
	InitJNI()
}

//export registerWallet
func registerWallet(account, token, network *C.char) *C.char {
	return C.CString(RegisterWallet(C.GoString(account), C.GoString(token), C.GoString(network)))
}

func main() {
	//InitJNI()
}
