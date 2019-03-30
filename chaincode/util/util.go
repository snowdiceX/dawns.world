package util

import (
	"encoding/json"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/snowdiceX/dawns.world/chaincode/log"
)

const (
	okJSON = `{"code": 200, "message": "OK"}`
)

// ChainResult response of chaincode
type ChainResult struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
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
