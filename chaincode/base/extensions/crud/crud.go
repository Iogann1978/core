package crud

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"s7ab-platform-hyperledger/platform/core/chaincode/base/extensions"
)

var (
	ErrKeyNotPresented = errors.New(`key isn't presented'`)
	ErrAlreadyExists   = errors.New(`state is already exists`)
	ErrNotExists       = errors.New(`state not exists`)
)

type Crudable interface {
	extensions.Keyable
	GetData(stub shim.ChaincodeStubInterface) ([]byte, error)
}

type Crud struct {
	tt Crudable
}

// Get entity by key
// Checks JSON structure of embedded Crudable
func (c Crud) Get(stub shim.ChaincodeStubInterface) pb.Response {
	key, err := c.tt.GetKey(stub)
	if err != nil {
		return c.error(err)
	}
	if data, err := stub.GetState(key); err != nil {
		return c.error(err)
	} else {
		if err = json.Unmarshal(data, c.tt); err != nil {
			return c.error(err)
		}
		return shim.Success(data)
	}
}

// Create Crudable entity
// Returns current TxID if successful
func (c Crud) Create(stub shim.ChaincodeStubInterface) pb.Response {
	if exists, err := c.isExists(stub); err != nil {
		return c.error(err)
	} else if exists == true {
		return c.error(ErrAlreadyExists)
	}

	key, err := c.tt.GetKey(stub)
	if err != nil {
		return c.error(err)
	}

	data, err := c.tt.GetData(stub)
	if err != nil {
		return c.error(err)
	}

	if err = stub.PutState(key, data); err != nil {
		return c.error(err)
	}
	return shim.Success([]byte(stub.GetTxID()))
}

// Update Crudable entity by key
// Returns current TxID if successful
func (c Crud) Update(stub shim.ChaincodeStubInterface) pb.Response {
	if exists, err := c.isExists(stub); err != nil {
		return c.error(err)
	} else if exists == false {
		return c.error(ErrNotExists)
	}

	key, err := c.tt.GetKey(stub)
	if err != nil {
		return c.error(err)
	}

	data, err := c.tt.GetData(stub)
	if err != nil {
		return c.error(err)
	}
	if err = stub.PutState(key, data); err != nil {
		return c.error(err)
	}
	return shim.Success([]byte(stub.GetTxID()))
}

// Delete Crudable entity by key
// Returns current TxID if successful
func (c Crud) Delete(stub shim.ChaincodeStubInterface) pb.Response {
	if exists, err := c.isExists(stub); err != nil {
		return c.error(err)
	} else if exists == false {
		return c.error(ErrNotExists)
	}
	key, err := c.tt.GetKey(stub)
	if err != nil {
		return c.error(err)
	}
	if err := stub.DelState(key); err != nil {
		return c.error(err)
	}
	return shim.Success([]byte(stub.GetTxID()))
}

func (c Crud) List(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Error(`not implemented`)
}

func (c Crud) isExists(stub shim.ChaincodeStubInterface) (bool, error) {
	key, err := c.tt.GetKey(stub)
	if err != nil {
		return false, err
	}
	if d, err := stub.GetState(key); err != nil {
		return false, err
	} else {
		return d == nil, nil
	}
}

func (c Crud) error(err error) pb.Response {
	return shim.Error(fmt.Sprintf("crud extension err: %s", err))
}
