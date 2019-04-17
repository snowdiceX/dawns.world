#! /bin/sh

./fabric-cli channel create --cid orgchannel \
  --txfile $GOPATH/src/github.com/securekey/fabric-examples/fabric-cli/fabric-sdk-go/test/fixtures/fabric/v1.4/channel/orgchannel.tx \
  --config $GOPATH/src/github.com/securekey/fabric-examples/fabric-cli/test/fixtures/config/config_test_local.yaml

./fabric-cli channel create --cid orgchannel \
  --txfile $GOPATH/src/github.com/securekey/fabric-examples/fabric-cli/fabric-sdk-go/test/fixtures/fabric/v1.4/channel/orgchannelOrg1MSPanchors.tx \
  --config $GOPATH/src/github.com/securekey/fabric-examples/fabric-cli/test/fixtures/config/config_test_local.yaml --orgid org1

./fabric-cli channel create --cid orgchannel \
  --txfile $GOPATH/src/github.com/securekey/fabric-examples/fabric-cli/fabric-sdk-go/test/fixtures/fabric/v1.4/channel/orgchannelOrg2MSPanchors.tx \
  --config $GOPATH/src/github.com/securekey/fabric-examples/fabric-cli/test/fixtures/config/config_test_local.yaml --orgid org2

./fabric-cli channel join --cid orgchannel \
  --config $GOPATH/src/github.com/securekey/fabric-examples/fabric-cli/test/fixtures/config/config_test_local.yaml

./fabric-cli chaincode install --ccp=github.com/snowdiceX/dawns.world/chaincode/wallet \
  --ccid=wallet --v v0 --gopath $GOPATH \
  --config $GOPATH/src/github.com/securekey/fabric-examples/fabric-cli/test/fixtures/config/config_test_local.yaml

./fabric-cli chaincode instantiate --cid orgchannel \
  --ccp=github.com/snowdiceX/dawns.world/chaincode/wallet \
  --ccid wallet --v v0 --args '{"Args":["A","1","B","2"]}' \
  --policy "AND('Org1MSP.member','Org2MSP.member')" \
  --config $GOPATH/src/github.com/securekey/fabric-examples/fabric-cli/test/fixtures/config/config_test_local.yaml

#./fabric-cli chaincode invoke --cid orgchannel --ccid wallet --args '{"Func":"create","Args":["ethereum", "eth", "addr_asdasdasd", "11", "tx_asdasdasdasd", "123"]}' --peer localhost:7051,localhost:8051 --payload --config $GOPATH/src/github.com/securekey/fabric-examples/fabric-cli/test/fixtures/config/config_test_local.yaml

./fabric-cli chaincode query --cid orgchannel --ccid wallet \
  --args '{"Func":"query","Args":["sequence", "1"]}' \
  --peer localhost:7051,localhost:8051 --payload \
  --config $GOPATH/src/github.com/securekey/fabric-examples/fabric-cli/test/fixtures/config/config_test_local.yaml
