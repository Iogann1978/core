package chaincodes

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/stretchr/testify/assert"
	"s7ab-platform-hyperledger/platform/core/utils"
	"testing"
)

func TestNewChainCode(t *testing.T) {
	baseCc := new(BaseSmartContract)
	cc := NewChainCode(baseCc, nil)
	if cc == nil {
		t.Fatal(`chaincode is nil`)
	}
	ms := shim.NewMockStub(`basechaincode`, cc)
	if r := ms.MockInit(`basechaincode`, nil); r.Status != shim.OK {
		t.Fatal(`chaincode error:`, r.Message)
	} else {
		assert.Equal(t, Version, string(r.Payload))
	}
}

func TestBaseSmartContract_Invoke(t *testing.T) {
	ms := shim.NewMockStub(`basechaincode`, NewChainCode(new(BaseSmartContract), nil))
	if ms == nil {
		t.Fatal(`chaincode is nil`)
	}

	if r := ms.MockInvoke(`basechaincode`, utils.ToChaincodeArgs(`TestInvoke`)); r.Status == shim.ERROR {
		t.Fatal(`chaincode error:`, r.Message)
	} else {
		assert.Equal(t, int32(shim.OK), r.Status)
		assert.Equal(t, TestSmartContractResponse, string(r.Payload))
	}
}
