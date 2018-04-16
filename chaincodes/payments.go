package chaincodes

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

type BasePayment struct {
	BaseSmartContract
	//base.Ownable
}

func (p *BasePayment) Init(stub shim.ChaincodeStubInterface) peer.Response {

	//args := stub.GetArgs()
	////protection from upgrade calls
	//if len(args) == 1 && !p.HasOwner(stub) {
	//	err := p.SetOwner(stub, args[0])
	//	if err != nil {
	//		return shim.Error(err.Error())
	//	}
	//}

	return shim.Success(nil)
}

func (p *BasePayment) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	return p.CallMethodByStubParameters(p, stub)
}

func (p *BasePayment) Add(stub shim.ChaincodeStubInterface) peer.Response {

	//if identity, passed := p.IsCallByOwner(stub); !passed {
	//	return shim.Error(fmt.Sprintf("Access to adding payment for org  %s denied ", identity.MspID))
	//}

	//fmt.Println(identity.MspID)

	return shim.Success(nil)
}
