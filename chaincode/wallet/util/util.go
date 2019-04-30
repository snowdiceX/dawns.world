package util

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/snowdiceX/dawns.world/chaincode/wallet/log"
)

const (
	okJSON  = `{"code": 200, "message": "OK"}`
	errJSON = `{"code": 500, "message": "Failed"}`
)

// ChaincodeError error of chaincode
type ChaincodeError struct {
	Code      int
	ErrString string
}

// Error string
func (e ChaincodeError) Error() string {
	return e.ErrString
}

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

// Hash return hash string of args
func Hash(args ...interface{}) string {
	hash := md5.New()
	hash.Write([]byte(fmt.Sprint(args...)))
	return hex.EncodeToString(hash.Sum(nil))
}

// Pagination of data
type Pagination struct {
	Records  []interface{}             `json:"records,omitempty"`
	Metadata *pb.QueryResponseMetadata `json:"metadata,omitempty"`
}

// TxRegister registered Tx
type TxRegister struct {
	Key      string `json:"key,omitempty"`
	Chain    string `json:"chain,omitempty"`
	Token    string `json:"token,omitempty"`
	Contract string `json:"contract,omitempty"`
	From     string `json:"from,omitempty"`
	To       string `json:"to,omitempty"`
	Amount   string `json:"amount,omitempty"`
	GasUsed  string `json:"gasUsed,omitempty"`
	GasPrice string `json:"gasPrice,omitempty"`
	Txhash   string `json:"txhash,omitempty"`
	Height   string `json:"height,omitempty"`
	Status   string `json:"status,omitempty"`
}

// BlockRegister registered block
type BlockRegister struct {
	Height string        `json:"height,omitempty"`
	Txs    []*TxRegister `json:"transactions,omitempty"`
}
