package util

import (
	"encoding/json"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/snowdiceX/dawns.world/chaincode/log"
)

const (
	okJSON  = `{"code": 200, "message": "OK"}`
	errJSON = `{"code": 500, "message": "Failed"}`
)

// ChainResult response of chaincode
type ChainResult struct {
	Code      int         `json:"code,omitempty"`
	Message   string      `json:"message,omitempty"`
	ErrString string      `json:"error,omitempty"`
	Result    interface{} `json:"result,omitempty"`
}

// Success response for chaincode call
func Success(ret *ChainResult) pb.Response {
	json, err := json.Marshal(ret)
	if err == nil {
		log.Info("response: ", string(json))
		return shim.Success(json)
	}
	log.Error("json marshal error: ", err)
	log.Info("response: ", okJSON)
	return shim.Success([]byte(okJSON))
}

// Error response for chaincode call
func Error(code int, errMsg string) pb.Response {
	ret := ChainResult{}
	ret.Code = code
	ret.Message = "Failed"
	ret.ErrString = errMsg
	json, err := json.Marshal(ret)
	if err == nil {
		log.Info("response: ", string(json))
		return shim.Error(string(json))
	}
	log.Error("json marshal error: ", err)
	log.Info("response: ", errJSON)
	return shim.Error(errJSON)
}
