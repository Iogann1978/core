package chaincode

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"

	"s7ab-platform-hyperledger/platform/core/logger"
	pt "s7ab-platform-hyperledger/platform/s7platform/testing"
	"s7ab-platform-hyperledger/platform/s7platform/tests/fixture"
)

func TestOrganizations(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Organizations Suite")
}

var _ = Describe("Organizations", func() {

	l := logger.NewZapLogger(nil)

	var operator, bank fixture.OrganizationFixture
	var org1, org2 fixture.MemberFixture

	var orgs *pt.FullMockStub

	BeforeSuite(func() {

		operator, _ = fixture.GetOrgFixture("Org1MSP.json")
		bank, _ = fixture.GetOrgFixture("Org2MSP.json")
		org1, _ = fixture.GetMemberFixture("Org3MSP.json", bank.OrganizationId, false)
		org2, _ = fixture.GetMemberFixture("Org4MSP.json", bank.OrganizationId, false)

		//fmt.Printf("Org: %+v\n", org1)

		orgs = pt.NewFullMockStub(`organizations`, NewOrganization(l))
		orgs.MockInit("1", orgs.ArgsToBytes(operator.OrganizationId))

	})

	Describe("Initialization", func() {

		It("Disallow add organizations", func() {
			pt.ExpectResponseError(orgs.MockInvokeFunc("/create", bank.GetBytes()), `Chaincode owner required`)
		})

		It("Allow to add organizations", func() {
			//start working from operator account
			orgs.MockCreator(operator.OrganizationId, operator.OrganizationCACert)

			for _, o := range [][]byte{bank.GetBytes(), org1.GetBytes(), org2.GetBytes()} {
				pt.ExpectResponseOk(orgs.MockInvokeFunc("/create", o))
			}

			bankFromChaincode, _ := fixture.GetOrgFromBytes(orgs.MockInvokeFunc("/get", bank.OrganizationId).Payload)
			Expect(bankFromChaincode.OrganizationId).To(Equal(bank.OrganizationId))
		})

		It("Disallow to add organizations with same OrganizationId or Itn ", func() {
			pt.ExpectResponseError(orgs.MockInvokeFunc("/create", bank.GetBytes()), `organization with id already in chaincode`)

			orgCopy := org1
			orgCopy.OrganizationId = "SomeIdThatNotExists"
			pt.ExpectResponseError(orgs.MockInvokeFunc("/create", orgCopy.GetBytes()), `organization with itn already in chaincode`)
		})

		It("Disallow to add organizations with empty parameters", func() {
			orgEmpty, _ := fixture.GetOrgFixture("OrgEmpty.json")
			pt.ExpectResponseError(orgs.MockInvokeFunc("/create", orgEmpty.GetBytes()))
		})

		It("Disallow add organizations without required parameters", func() {
			orgCopy := org1
			orgCopy.Requisites.ITN = ""

			//fmt.Printf("Org: %+v\n", orgCopy)

			pt.ExpectResponseError(orgs.MockInvokeFunc("/create", orgCopy.GetBytes()))

		})

	})

})
