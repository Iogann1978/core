package base

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"s7ab-platform-hyperledger/platform/core/utils"
)

type Ownable struct{}

const OWNER_KEY = "OWNER"

func (o *Ownable) SetOwner(stub shim.ChaincodeStubInterface, owner []byte) error {
	return stub.PutState(OWNER_KEY, owner)
}

func (o *Ownable) HasOwner(stub shim.ChaincodeStubInterface) bool {
	return len(o.GetOwner(stub)) != 0
}

func (o *Ownable) GetOwner(stub shim.ChaincodeStubInterface) []byte {
	owner, _ := stub.GetState(OWNER_KEY)
	return owner
}

func (o *Ownable) IsCallByOwner(stub shim.ChaincodeStubInterface) (*utils.Creator, bool) {
	creator, _ := stub.GetCreator()
	identity, _ := utils.NewCreator(creator)

	return identity, string(o.GetOwner(stub)) == identity.MspID
}
