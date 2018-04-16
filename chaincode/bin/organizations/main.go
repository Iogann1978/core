package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"s7ab-platform-hyperledger/platform/core/chaincode"
	"s7ab-platform-hyperledger/platform/core/logger"
)

func main() {
	l := logger.NewZapLogger(nil)
	cc := chaincode.NewOrganization(l)

	if err := shim.Start(cc); err != nil {
		l.Warn(`chaincode`, logger.KV(`error`, err))
	}
}
