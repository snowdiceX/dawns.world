package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"

	"github.com/snowdiceX/dawns.world/chaincode/log"
	"github.com/snowdiceX/dawns.world/chaincode/util"
)

/**
* args:
      0 blockchain network
      1 block data
*/
func (w *WalletChaincode) registerBlock(
	stub shim.ChaincodeStubInterface, args []string) pb.Response {
	network, blockdata := args[0], args[1]
	return handleBlock(network, blockdata)
}

func handleBlock(network, blockdata string) pb.Response {
	log.Info("register chain block: ", network, "; ", blockdata)
	ret := &util.ChainResult{Code: 200, Message: "OK"}
	return util.Success(ret)
}
