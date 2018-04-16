package chaincodes

import (
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/stretchr/testify/assert"
	"s7ab-platform-hyperledger/platform/core/utils"
	s7t "s7ab-platform-hyperledger/platform/s7platform/testing"
	"s7ab-platform-hyperledger/platform/s7platform/tests/fixture"
	"testing"
)

func TestBasePayment_Invoke(t *testing.T) {
	payments := new(BasePayment)

	ms := s7t.NewFullMockStub("payments", payments)

	//Org2MSP is administrator of chaincode
	ms.MockInit("0", utils.ToChaincodeArgs(`Org2MSP`))

	ms.MockCreator("Org2MSP", fixture.ORG2_CA_CERT)

	fmt.Println("Lets go")

	if r := ms.MockInvoke(`1`, utils.ToChaincodeArgs(`Add`)); r.Status == shim.ERROR {
		t.Fatal(`chaincode error:`, r.Message)
	} else {
		assert.Equal(t, int32(shim.OK), r.Status)

	}

	ms.MockCreator("Org3MSP", fixture.ORG3_CA_CERT)

	r := ms.MockInvoke(`2`, utils.ToChaincodeArgs(`Add`))
	assert.Equal(t, int32(shim.ERROR), r.Status)

}
