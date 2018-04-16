package chaincodes

import (
	"encoding/json"
	"errors"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"reflect"
	"s7ab-platform-hyperledger/platform/core/logger"
	"strconv"
)

var (
	ErrMethodNotFound   = errors.New(`method not found`)
	ErrGetStateInvalid  = errors.New(`invalid state arguments`)
	ErrTooManyArguments = errors.New(`too many arguments to return`)
	ErrInvalidType      = errors.New(`invalid type`)
)

const (
	TestSmartContractResponse = `it works!`
	Version                   = `0.1`
)

type BaseSmartContract struct {
	shim.Chaincode
	cc  shim.Chaincode
	Log logger.Logger
}

func (bs BaseSmartContract) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success([]byte(Version))
}

func (bs BaseSmartContract) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	function, _ := stub.GetFunctionAndParameters()
	if method := reflect.ValueOf(bs.cc).MethodByName(function); method.IsValid() {
		result := method.Call([]reflect.Value{reflect.ValueOf(stub)})
		if len(result) != 1 {
			return shim.Error(ErrTooManyArguments.Error())
		}
		if response, ok := result[0].Interface().(peer.Response); ok {
			return response
		}
		return shim.Error(ErrInvalidType.Error())
	}
	return shim.Error(ErrMethodNotFound.Error())
}

func (bs *BaseSmartContract) CallMethodByStubParameters(cc shim.Chaincode, APIstub shim.ChaincodeStubInterface) peer.Response {
	functionName, _ := APIstub.GetFunctionAndParameters()

	method := bs.GetMethod(cc, functionName)

	if !method.IsValid() {
		return shim.Error(ErrMethodNotFound.Error())
	}

	result := method.Call([]reflect.Value{reflect.ValueOf(APIstub)})
	if len(result) != 1 {
		return shim.Error(ErrTooManyArguments.Error())
	}
	if response, ok := result[0].Interface().(peer.Response); ok {
		return response
	}

	return shim.Error(ErrInvalidType.Error())
}

func (bs *BaseSmartContract) GetMethod(cc shim.Chaincode, functionName string) reflect.Value {
	return reflect.ValueOf(cc).MethodByName(functionName)
}

func (bs *BaseSmartContract) GetByKey(stub shim.ChaincodeStubInterface) peer.Response {
	_, args := stub.GetFunctionAndParameters()
	if len(args) != 1 {
		return shim.Error(ErrGetStateInvalid.Error())
	}
	if state, err := stub.GetState(args[0]); err != nil {
		return shim.Error(err.Error())
	} else {
		return shim.Success(state)
	}
}

func (bs *BaseSmartContract) List(stub shim.ChaincodeStubInterface) peer.Response {
	_, args := stub.GetFunctionAndParameters()
	if len(args) != 2 {
		return shim.Error(ErrTooManyArguments.Error())
	}
	limit, err := strconv.Atoi(args[0])
	if err != nil {
		bs.Log.Debug(`invalid type:`, logger.KV(`limit`, err))
	}
	offset, err := strconv.Atoi(args[1])
	if err != nil {
		bs.Log.Debug(`invalid type:`, logger.KV(`offset`, err))
	}
	iter, err := stub.GetQueryResult(Query{Limit: &limit, Offset: &offset}.String())
	if err != nil {
		return shim.Error(err.Error())
	}

	defer iter.Close()

	rawData := make([]map[string]interface{}, 0)
	jsonData := make(map[string]interface{})

	for iter.HasNext() {
		data, err := iter.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		if err := json.Unmarshal(data.Value, &jsonData); err != nil {
			return shim.Error(err.Error())
		}
		rawData = append(rawData, jsonData)
	}

	data, err := json.Marshal(jsonData)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(data)
}

func (bs *BaseSmartContract) TestInvoke(stub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success([]byte(TestSmartContractResponse))
}

func NewChainCode(cc shim.Chaincode, l logger.Logger) shim.Chaincode {
	bs := new(BaseSmartContract)
	bs.cc = cc
	bs.Log = l
	return bs
}
