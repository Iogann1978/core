package extensions

import "github.com/hyperledger/fabric/core/chaincode/shim"

type Keyable interface {
	GetKey(stub shim.ChaincodeStubInterface) (string, error)
}
