package base

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"s7ab-platform-hyperledger/platform/core/chaincodes"
	"s7ab-platform-hyperledger/platform/core/logger"
)

type Chaincode struct {
	Log logger.Logger
}

func (cc Chaincode) WriteError(err interface{}) pb.Response {
	return shim.Error(fmt.Sprintf("%s", err))
}

func (cc Chaincode) WriteSuccess(data interface{}) pb.Response {
	switch data.(type) {
	case string:
		return shim.Success([]byte(data.(string)))
	case []byte:
		return shim.Success(data.([]byte))
	default:
		b, err := json.Marshal(data)
		if err != nil {
			cc.Log.Warn(`json marshal error`, logger.KV(`err`, err.Error()))
			return shim.Success(nil)
		} else {
			return shim.Success(b)
		}
	}
}

func (cc Chaincode) GetCreator(stub shim.ChaincodeStubInterface) (*chaincodes.Creator, error) {
	payload, err := stub.GetCreator()
	if err != nil {
		return nil, err
	}

	return chaincodes.NewCreator(payload)
}

func (cc Chaincode) ToChaincodeArgs(args ...string) [][]byte {
	bargs := make([][]byte, len(args))
	for i, arg := range args {
		bargs[i] = []byte(arg)
	}
	return bargs
}
