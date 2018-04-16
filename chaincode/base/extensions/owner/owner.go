package owner

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"s7ab-platform-hyperledger/platform/core/chaincodes"
	"s7ab-platform-hyperledger/platform/core/logger"
)

const DefaultKey = `OWNER`

var ErrToMuchArguments = errors.New(`too much arguments`)

type Owner struct {
	key string
	log logger.Logger
}

// Get returns current owner
func (o *Owner) Get(stub shim.ChaincodeStubInterface) ([]byte, error) {
	return stub.GetState(o.key)
}

// Set current chaincode owner
// Sets current MspID from stub creator if owner isn't presented
// Returns fixed owner and error if exists
func (o *Owner) Set(stub shim.ChaincodeStubInterface, owner ...[]byte) ([]byte, error) {
	switch len(owner) {
	case 0:

		payload, err := stub.GetCreator()
		if err != nil {
			return nil, err
		}

		creator, err := chaincodes.NewCreator(payload)

		if err != nil {
			return nil, err
		} else {
			o.log.Info(`owner`, logger.KV(`Set from MspID`, creator.MspID))
			return []byte(creator.MspID), stub.PutState(o.key, []byte(creator.MspID))
		}

		return []byte(creator.MspID), stub.PutState(o.key, []byte(creator.MspID))
	case 1:
		o.log.Info(`owner`, logger.KV(`Set from cli`, string(owner[0])))
		return owner[0], stub.PutState(o.key, owner[0])
	default:
		return nil, ErrToMuchArguments
	}
}

func (o *Owner) SetFromFirstArgOrCreator(stub shim.ChaincodeStubInterface) pb.Response {
	_, args := stub.GetFunctionAndParameters()
	if len(args) == 1 {

		if _, err := o.Set(stub, []byte(args[0])); err != nil {
			return shim.Error(fmt.Sprintf("%s", err))
		} else {
			return shim.Success([]byte(args[0]))
		}
	} else {
		if o, err := o.Set(stub); err != nil {
			return shim.Error(fmt.Sprintf("%s", err))
		} else {
			return shim.Success([]byte(o))
		}
	}

}

// IsOwner checks chaincode owner
// Uses current MspID from stub creator if owner isn't presented
func (o *Owner) IsOwner(stub shim.ChaincodeStubInterface, owner ...[]byte) (bool, error) {
	switch len(owner) {
	case 0:

		payload, err := stub.GetCreator()
		if err != nil {
			return false, err
		}

		if creator, err := chaincodes.NewCreator(payload); err != nil {
			return false, err
		} else {
			if cb, err := stub.GetState(o.key); err != nil {
				return false, err
			} else {
				return bytes.Equal([]byte(creator.MspID), cb), nil
			}
		}
	case 1:
		if co, err := stub.GetState(o.key); err != nil {
			return false, err
		} else {
			return bytes.Equal(co, owner[0]), nil
		}
	default:
		return false, ErrToMuchArguments
	}
}

// NewOwner returns new owner instance for using in chaincode
func NewOwner(l logger.Logger, key ...string) *Owner {
	if len(key) == 1 {
		return &Owner{key: key[0], log: l}
	}
	return &Owner{key: DefaultKey, log: l}
}
