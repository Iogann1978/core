package chaincode

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"s7ab-platform-hyperledger/platform/core/chaincode/base"
	extCreator "s7ab-platform-hyperledger/platform/core/chaincode/base/extensions/creator"
	"s7ab-platform-hyperledger/platform/core/chaincode/base/extensions/owner"
	"s7ab-platform-hyperledger/platform/core/chaincode/base/extensions/router"
	"s7ab-platform-hyperledger/platform/core/entities"
	"s7ab-platform-hyperledger/platform/core/logger"
)

const OrganizationKey = `ORGANIZATIONS`

type Organization struct {
	base.Chaincode
	// prefix for state key with org json
	key   string
	r     *router.Group
	owner *owner.Owner
}

func NewOrganization(l logger.Logger) Organization {
	o := Organization{key: OrganizationKey}
	o.Log = l

	r := router.New()
	o.owner = owner.NewOwner(l)

	r.Add(`/create`, o.create)
	r.Add(`/get`, o.get)

	// add bank handlers
	bank := r.Group(`/bank`)
	bank.Add(`/list`, o.bankList)

	// add bank-> member handlers
	bankMembers := bank.Group(`/member`)
	bankMembers.Add(`/list`, o.bankMemberList)
	bankMembers.Add(`/confirm`, o.bankMemberConfirm)
	bankMembers.Add(`/unconfirm`, o.bankMemberUnConfirm)

	member := r.Group(`/member`)
	member.Add(`/list`, o.memberList)
	member.Add(`/byITN`, o.memberByITN)

	o.r = r

	return o
}

func (o Organization) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return o.owner.SetFromFirstArgOrCreator(stub)
}

func (o Organization) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	return o.r.Handle(stub)
}

func (o Organization) prepareOrgKey(stub shim.ChaincodeStubInterface, org entities.Organization) (key string, err error) {
	key, err = stub.CreateCompositeKey(o.key, []string{string(org.Type), org.OrganizationId})
	return
}

func (o Organization) prepareOrgKeyById(stub shim.ChaincodeStubInterface, id string, orgType entities.OrganizationType) (key string, err error) {
	key, err = stub.CreateCompositeKey(o.key, []string{string(orgType), id})
	return
}

func (o Organization) prepareOrgKeyByItn(itn string) (key string) {
	return "ORG_ID_BY_ITN_" + itn
}

func (o Organization) validateOrgData(org entities.Organization) pb.Response {

	if org.OrganizationId == "" {
		return shim.Error(fmt.Sprintf("missing property OrganizationId: %v", org))
	}

	if org.OrganizationCert == "" {
		return shim.Error(fmt.Sprintf("missing property OrganizationCert: %v", org))
	}

	if org.OrganizationCACert == "" {
		return shim.Error(fmt.Sprintf("missing property OrganizationCACert: %v", org))
	}

	if org.Requisites.BIC == "" {
		return shim.Error(fmt.Sprintf("missing property Requisites.BIC: %v", org))
	}

	if org.Requisites.CorrespondentAccount == "" {
		return shim.Error(fmt.Sprintf("Missing property Requisites.CorrespondentAccount: %v", org))
	}

	if org.Requisites.ITN == "" {
		return shim.Error(fmt.Sprintf("Missing property Requisites.ITN: %v", org))
	}

	if org.Requisites.IEC == "" {
		return shim.Error(fmt.Sprintf("Missing property Requisites.IEC: %v", org))
	}

	if org.Type != entities.BANK_TYPE && org.Type != entities.MEMBER_TYPE {
		return shim.Error(fmt.Sprintf("Provided payload does not contain any type of BANK or MEMBER: %s", org))
	}

	return shim.Success([]byte{})
}

func (o Organization) getOrgIdByItn(stub shim.ChaincodeStubInterface, itn string) (string, error) {
	key, err := stub.GetState(o.prepareOrgKeyByItn(itn))
	return string(key), err
}

func (o Organization) create(stub shim.ChaincodeStubInterface) pb.Response {

	if isOwner, _ := o.owner.IsOwner(stub); !isOwner {
		return shim.Error(`Chaincode owner required`)
	}

	_, args := stub.GetFunctionAndParameters()
	payload := []byte(args[0])

	var org entities.Organization
	err := json.Unmarshal(payload, &org)
	if err != nil {
		return shim.Error(err.Error())
	}

	key, err := o.prepareOrgKey(stub, org)
	val, err := stub.GetState(key)
	if val != nil {
		return shim.Error(fmt.Sprintf("organization with id already in chaincode: id=%s %s", org.OrganizationId, val))
	}

	if isValid := o.validateOrgData(org); isValid.Status == shim.ERROR {
		return isValid
	}

	if orgKey, _ := o.getOrgIdByItn(stub, org.Requisites.ITN); len(orgKey) > 0 {
		return shim.Error(fmt.Sprintf("organization with itn already in chaincode: id=%s %s", org.OrganizationId, val))
	}

	if org.Type == entities.MEMBER_TYPE {
		var member entities.Member

		err = json.Unmarshal(payload, &member)
		if err != nil {
			return shim.Error(err.Error())
		}

		if member.ConfirmedByBank == true {
			return shim.Error(fmt.Sprintf("You trying to add already confirmed member: %v", member))
		}

		if member.Requisites.SettlementAccount == "" {
			return shim.Error(fmt.Sprintf("Missing property Requisites.SettlementAccount: %v", member))
		}

	} else if org.Type == entities.BANK_TYPE {

		var bank entities.Bank
		err = json.Unmarshal(payload, &bank)
		if err != nil {
			return shim.Error(err.Error())
		}
	}

	if err = stub.PutState(key, payload); err != nil {
		return shim.Error(err.Error())
	}

	if err = stub.PutState(o.prepareOrgKeyByItn(org.Requisites.ITN), []byte(org.OrganizationId)); err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte("WriteSuccess"))
}

func (o Organization) bankList(stub shim.ChaincodeStubInterface) pb.Response {
	var creds []entities.Bank

	iter, err := stub.GetStateByPartialCompositeKey(o.key, []string{string(entities.BANK_TYPE)})
	if err != nil {
		return shim.Error(err.Error())
	}

	defer iter.Close()

	for iter.HasNext() {
		v, err := iter.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		o.Log.Debug(`bankList entry`, logger.KV(`key`, v.Key))

		var entry entities.Bank
		err = json.Unmarshal(v.Value, &entry)

		creds = append(creds, entry)
	}

	if len(creds) == 0 {
		return shim.Success([]byte{})
	}

	result, err := json.Marshal(creds)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(result)
}

func (o Organization) memberList(stub shim.ChaincodeStubInterface) pb.Response {
	var creds []entities.Member

	iter, err := stub.GetStateByPartialCompositeKey(o.key, []string{string(entities.MEMBER_TYPE)})
	if err != nil {
		return shim.Error(err.Error())
	}

	defer iter.Close()

	for iter.HasNext() {
		v, err := iter.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		var entry entities.Member
		err = json.Unmarshal(v.Value, &entry)

		creds = append(creds, entry)
	}

	if len(creds) == 0 {
		return shim.Success([]byte{})
	}

	result, err := json.Marshal(creds)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(result)
}

// get org data by org id
func (o Organization) get(stub shim.ChaincodeStubInterface) pb.Response {
	_, args := stub.GetFunctionAndParameters()

	if data, err := o.getOrgDataById(stub, args[0]); err != nil {
		return shim.Error(err.Error())
	} else {
		return shim.Success(data)
	}
}

func (o Organization) getOrgDataById(stub shim.ChaincodeStubInterface, id string) (data []byte, err error) {

	memberKey, err := o.prepareOrgKeyById(stub, id, entities.MEMBER_TYPE)
	if err != nil {
		return data, err
	}

	if data, err = stub.GetState(memberKey); err == nil && data != nil {
		return data, err
	}

	bankKey, err := o.prepareOrgKeyById(stub, id, entities.BANK_TYPE)
	if err != nil {
		return data, err
	}

	data, err = stub.GetState(bankKey)

	return data, err
}

func (o Organization) bankMemberList(stub shim.ChaincodeStubInterface) pb.Response {
	creator, err := stub.GetCreator()
	if err != nil {
		return shim.Error(err.Error())
	}

	identity, err := extCreator.NewCreator(creator)
	if err != nil {
		return shim.Error(err.Error())
	}

	bankKey, err := stub.CreateCompositeKey(o.key, []string{string(entities.BANK_TYPE), identity.MspID})
	if err != nil {
		return shim.Error(err.Error())
	}

	resultBytes, err := stub.GetState(bankKey)
	if err != nil {
		return shim.Error(err.Error())
	}

	if resultBytes == nil {
		return shim.Error(fmt.Sprintf("[MembersByBank]: Organization what invoke this does not found in chaincode, org name: %s", identity.MspID))
	}

	var org entities.Organization
	err = json.Unmarshal(resultBytes, &org)
	if err != nil {
		return shim.Error(err.Error())
	}

	if org.Type != entities.BANK_TYPE {
		return shim.Error(fmt.Sprintf("You are not a bank: %s", identity.MspID))
	}

	queryString := fmt.Sprintf("{\"selector\":{\"type\":\"%s\",\"bank_organization_id\":\"%s\"}}", entities.MEMBER_TYPE, identity.MspID)

	query, err := stub.GetQueryResult(queryString)

	if query == nil {
		return shim.Error(fmt.Sprintf("Query in nil: %s", queryString))
	}

	defer query.Close()

	var result []entities.Member

	for query.HasNext() {
		v, err := query.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		var member entities.Member

		err = json.Unmarshal(v.Value, &member)
		if err != nil {
			return shim.Error(err.Error())
		}

		result = append(result, member)
	}

	responseBytes, err := json.Marshal(result)

	return shim.Success(responseBytes)
}

func (o Organization) memberByITN(stub shim.ChaincodeStubInterface) pb.Response {
	_, args := stub.GetFunctionAndParameters()

	id, _ := o.getOrgIdByItn(stub, string(args[0]))

	if len(id) == 0 {
		return shim.Error(`itn not found`)
	}

	if data, err := o.getOrgDataById(stub, id); err == nil {
		return shim.Success(data)
	} else {
		return shim.Error(`org data not found`)
	}
}

func (o Organization) bankMemberConfirm(stub shim.ChaincodeStubInterface) pb.Response {
	creator, err := stub.GetCreator()
	if err != nil {
		return shim.Error(err.Error())
	}

	identity, err := extCreator.NewCreator(creator)
	if err != nil {
		return shim.Error(err.Error())
	}

	bankKey, err := stub.CreateCompositeKey(o.key, []string{string(entities.BANK_TYPE), identity.MspID})
	if err != nil {
		return shim.Error(err.Error())
	}

	resultBytes, err := stub.GetState(bankKey)
	if err != nil {
		return shim.Error(err.Error())
	}

	if resultBytes == nil {
		return shim.Error(fmt.Sprintf("[MembersByBank]: Organization what invoke this does not found in chaincode, org name: %s", identity.MspID))
	}

	var org entities.Organization
	err = json.Unmarshal(resultBytes, &org)
	if err != nil {
		return shim.Error(err.Error())
	}

	if org.Type != entities.BANK_TYPE {
		return shim.Error(fmt.Sprintf("You are not a bank: %s", identity.MspID))
	}

	_, args := stub.GetFunctionAndParameters()

	memberKey, err := stub.CreateCompositeKey(o.key, []string{string(entities.MEMBER_TYPE), args[0]})
	if err != nil {
		return shim.Error(err.Error())
	}

	memberBytes, err := stub.GetState(memberKey)
	if err != nil {
		return shim.Error(err.Error())
	}

	var member entities.Member
	err = json.Unmarshal(memberBytes, &member)
	if err != nil {
		return shim.Error(err.Error())
	}

	if member.BankOrganizationId != identity.MspID {
		return shim.Error(fmt.Sprintf("This member not in your Bank: %s", args[0]))
	}

	if member.ConfirmedByBank == true {
		return shim.Error(fmt.Sprintf("This member already confirmed: %s", args[0]))
	}

	member.ConfirmedByBank = true
	memberBytesConfirmed, err := json.Marshal(member)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(memberKey, memberBytesConfirmed)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (o Organization) bankMemberUnConfirm(stub shim.ChaincodeStubInterface) pb.Response {
	creator, err := stub.GetCreator()
	if err != nil {
		return shim.Error(err.Error())
	}

	identity, err := extCreator.NewCreator(creator)
	if err != nil {
		return shim.Error(err.Error())
	}

	bankKey, err := stub.CreateCompositeKey(o.key, []string{string(entities.BANK_TYPE), identity.MspID})
	if err != nil {
		return shim.Error(err.Error())
	}

	resultBytes, err := stub.GetState(bankKey)
	if err != nil {
		return shim.Error(err.Error())
	}

	if resultBytes == nil {
		return shim.Error(fmt.Sprintf("[MembersByBank]: Organization what invoke this does not found in chaincode, org name: %s", identity.MspID))
	}

	var org entities.Organization
	err = json.Unmarshal(resultBytes, &org)
	if err != nil {
		return shim.Error(err.Error())
	}

	if org.Type != entities.BANK_TYPE {
		return shim.Error(fmt.Sprintf("You are not a bank: %s", identity.MspID))
	}

	_, args := stub.GetFunctionAndParameters()

	memberKey, err := stub.CreateCompositeKey(o.key, []string{string(entities.MEMBER_TYPE), args[0]})
	if err != nil {
		return shim.Error(err.Error())
	}

	memberBytes, err := stub.GetState(memberKey)
	if err != nil {
		return shim.Error(err.Error())
	}

	var member entities.Member
	err = json.Unmarshal(memberBytes, &member)
	if err != nil {
		return shim.Error(err.Error())
	}

	if member.BankOrganizationId != identity.MspID {
		return shim.Error(fmt.Sprintf("This member not in your Bank: %s", args[0]))
	}

	if member.ConfirmedByBank == false {
		return shim.Error(fmt.Sprintf("This member not confirmed already: %s", args[0]))
	}

	member.ConfirmedByBank = false
	memberBytesUnconfirmed, err := json.Marshal(member)
	if err != nil {
		return shim.Error(err.Error())
	}

	if err = stub.PutState(memberKey, memberBytesUnconfirmed); err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}
