package main

import (
	pb "github.com/hyperledger/fabric/protos/peer"

	"github.com/snowdiceX/dawns.world/chaincode/log"
	"github.com/snowdiceX/dawns.world/chaincode/util"
)

func registerChainBlock(network, blockdata string) pb.Response {
	log.Info("register chain block: ", network, "; ", blockdata)
	ret := &util.ChainResult{Code: 200, Message: "OK"}
	return util.Success(ret)
}
