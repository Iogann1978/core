package multisig

import (
	"errors"
)

var (
	ErrAlreadySigned = errors.New(`contact already signed`)
)

type Multisig struct{}

//func (m *Multisig ) Init(APIstub shim.ChaincodeStubInterface) pb.Response {
//	return shim.Success(nil)
//}
//
//func (m *Multisig ) Invoke(APIstub shim.ChaincodeStubInterface) pb.Response {
//	function, args := APIstub.GetFunctionAndParameters()
//	switch function {
//	//case `Add`:
//	//	return s.add(APIstub, args)
//	case `All`:
//		return m.all(APIstub)
//	case `Get`:
//		return m.get(APIstub, args)
//	case `Approve`:
//		return m.approve(APIstub, args)
//
//	}
//	return shim.Error(`unknown method`)
//}
//
//// ====================================
////	Add new organization to chaincode
////
//// ====================================
////func (s *SmartContract) add(APIstub shim.ChaincodeStubInterface, args []string) pb.Response {
////	if len(args) != 1 {
////		return shim.Error(fmt.Sprintf("expected args length: %d, got %d", 1, len(args)))
////	}
////	var org entities.OrganizationCredentials
////	if err := json.Unmarshal([]byte(args[0]), &org); err != nil {
////		return shim.Error(err.Error())
////	}
////
////	if err := APIstub.PutState(org.Key(), []byte(args[0])); err != nil {
////		return shim.Error(err.Error())
////	}
////
////	APIstub.SetEvent(`test`, []byte(`asdasdasdasd`))
////
////	return shim.Success(nil)
////}
//
////	===========================================
////	Get all organizations
////
////	===========================================
//func (s *Multisig) all(APIstub shim.ChaincodeStubInterface) pb.Response {
//	iter, err := APIstub.GetStateByRange(entities.ORG_KEY, ``)
//	if err != nil {
//		return shim.Error(err.Error())
//	}
//	defer iter.Close()
//
//	c, _ := APIstub.GetCreator()
//	cr, _ := chaincodes.NewCreator(c)
//	fmt.Printf("%v", *cr)
//
//	var out entities.OrganizationRequest
//	var requests []entities.OrganizationRequest
//	for iter.HasNext() {
//		res, err := iter.Next()
//		if err != nil {
//			return shim.Error(err.Error())
//		}
//		v := res.GetValue()
//		if err := json.Unmarshal(v, &out); err != nil {
//			return shim.Error(err.Error())
//		} else {
//			requests = append(requests, out)
//		}
//	}
//
//	reqOut, err := json.Marshal(requests)
//	if err != nil {
//		return shim.Error(err.Error())
//	}
//	return shim.Success(reqOut)
//}
//
////========================================
//// Get sign contract
////
////========================================
//func (s *Multisig) get(APIstub shim.ChaincodeStubInterface, args []string) pb.Response {
//	// Check current organization existence in organizations chaincode
//	c, _ := APIstub.GetCreator()
//	if creator, err := chaincodes.NewCreator(c); err != nil {
//		return shim.Error(err.Error())
//	} else {
//		b, _ := json.Marshal(creator)
//		fmt.Println(string(b))
//	}
//
//	//resp := APIstub.InvokeChaincode(`organizations`, util.ToChaincodeArgs(`Get`, args[0]), `organizations`)
//	//if resp.Status == shim.ERROR {
//	//	return shim.Error(resp.Message)
//	//}
//	return shim.Success(nil)
//}
//
//func (s *Multisig) approve(APIStub shim.ChaincodeStubInterface, args []string) pb.Response {
//
//	if len(args) != 1 {
//		return shim.Error(fmt.Sprintf("expected args length: %d, got %d", 1, len(args)))
//	}
//
//	key := entities.GetOrganizationKey(args[0], args[1])
//
//	state, err := APIStub.GetState(key)
//	if err != nil {
//		return shim.Error(err.Error())
//	}
//
//	var request entities.OrganizationRequest
//	if err := json.Unmarshal(state, &request); err != nil {
//		return shim.Error(err.Error())
//	}
//
//	cb, err := APIStub.GetCreator()
//	if err != nil {
//		return shim.Error(err.Error())
//	}
//	creator, err := chaincodes.NewCreator(cb)
//	if err != nil {
//		return shim.Error(err.Error())
//	}
//	if _, ok := request.Approves[creator.MspID]; ok {
//		return shim.Error(ErrAlreadySigned.Error())
//	} else {
//		request.Approves[creator.MspID] = creator.User
//		if b, err := json.Marshal(request); err != nil {
//			return shim.Error(err.Error())
//		} else {
//			APIStub.PutState(key, b)
//		}
//	}
//
//	return shim.Success(nil)
//
//}

