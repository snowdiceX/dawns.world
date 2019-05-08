package util

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"

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

// TxInfo info of Tx
type TxInfo struct {
	Contract string `json:"contract,omitempty"`
	From     string `json:"from,omitempty"`
	To       string `json:"to,omitempty"`
	Amount   string `json:"amount,omitempty"`
	GasUsed  string `json:"gasUsed,omitempty"`
	GasPrice string `json:"gasPrice,omitempty"`
	TxHash   string `json:"txHash,omitempty"`
	Height   string `json:"height,omitempty"`
	Status   string `json:"status,omitempty"`
}

// TxRegister registered Tx
type TxRegister struct {
	Key       string  `json:"key,omitempty"`
	ChainName string  `json:"chain,omitempty"`
	TokenName string  `json:"token,omitempty"`
	Addr      string  `json:"address,omitempty"`
	Amount    string  `json:"amount,omitempty"`
	GasUsed   string  `json:"gasUsed,omitempty"`
	GasPrice  string  `json:"gasPrice,omitempty"`
	Info      *TxInfo `json:"info,omitempty"`
}

// Chain returns chain name
func (t *TxRegister) Chain() string {
	return t.ChainName
}

// Token returns token name
func (t *TxRegister) Token() string {
	return t.TokenName
}

// Address returns the address of the modified change
func (t *TxRegister) Address() string {
	return t.Addr
}

// AmountHex returns the amount of the transaction in hex format
func (t *TxRegister) AmountHex() string {
	return t.Amount
}

// GasUsedHex returns the gas used of the transaction in hex format
func (t *TxRegister) GasUsedHex() string {
	return t.GasUsed
}

// GasPriceHex returns the gas price of the transaction in hex format
func (t *TxRegister) GasPriceHex() string {
	return t.GasPrice
}

// BlockRegister registered block
type BlockRegister struct {
	Height string        `json:"height,omitempty"`
	Txs    []*TxRegister `json:"transactions,omitempty"`
}

// CheckState check world state
func CheckState(stub shim.ChaincodeStubInterface,
	key string, returnExistError bool) ([]byte, *ChaincodeError) {
	bytes, err := stub.GetState(key)
	if err != nil {
		log.Errorf("check state error: %s: %v", key, err)
		ccErr := &ChaincodeError{
			Code:      http.StatusInternalServerError,
			ErrString: fmt.Sprintf("check state error: %s: %v", key, err)}
		return nil, ccErr
	}
	if returnExistError && bytes != nil {
		log.Warnf("check state error: %s: state exist ", key)
		ccErr := &ChaincodeError{
			Code:      http.StatusConflict,
			ErrString: fmt.Sprintf("state exist: %s", key)}
		return bytes, ccErr
	}
	return bytes, nil
}
