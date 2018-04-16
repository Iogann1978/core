package meta

import (
	"errors"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"s7ab-platform-hyperledger/platform/core/chaincode/base/extensions"
)

const DefaultPrefix = `META_`

var (
	ErrMetaKeyNotPresented = errors.New(`meta key isn't presented`)
)

type Metable interface {
	extensions.Keyable
	// GetMetaKey method is used for getting key which will be used in meta state
	GetMetaKey(stub shim.ChaincodeStubInterface) (string, error)
	// GetMetaData method is used for getting data which used for saving
	GetMetaData(stub shim.ChaincodeStubInterface) ([]byte, error)
	// GetStateDataWithMeta method is used for internal setting data in cc state
	GetStateDataWithMeta(stub shim.ChaincodeStubInterface) ([]byte, error)
}

type Meta struct {
	prefix string
	ms     Metable
}

// SetMeta method is used for setting meta in state
// Puts two states:
// 1 - state with prefix_ and result of GetMetaKey() key and state with result of GetMetaData()
// 2 - state with result of GetKey() key and state with result of GetStateDataWithMeta()
func (m *Meta) SetMeta(stub shim.ChaincodeStubInterface) pb.Response {
	if key, err := m.getKey(stub); err != nil {
		return m.error(err)
	} else {
		if data, err := m.ms.GetMetaData(stub); err != nil {
			return m.error(err)
		} else {
			if err = stub.PutState(key, data); err != nil {
				return m.error(err)
			} else {
				if stateData, err := m.ms.GetStateDataWithMeta(stub); err != nil {
					return m.error(err)
				} else {
					if stateKey, err := m.ms.GetKey(stub); err != nil {
						return m.error(err)
					} else {
						if err = stub.PutState(stateKey, stateData); err != nil {
							return m.error(err)
						}
					}
				}
			}
		}
	}
	return shim.Success([]byte(stub.GetTxID()))
}

// GetMeta method is used for getting meta from state
//
func (m *Meta) GetMeta(stub shim.ChaincodeStubInterface) pb.Response {
	if key, err := m.getKey(stub); err != nil {
		return m.error(err)
	} else {
		if data, err := stub.GetState(key); err != nil {
			return m.error(err)
		} else {
			return shim.Success(data)
		}
	}
}

func (m *Meta) getKey(stub shim.ChaincodeStubInterface) (string, error) {
	metaKey, err := m.ms.GetMetaKey(stub)
	if err != nil {
		return ``, err
	}
	return fmt.Sprintf("%s_%s", m.prefix, metaKey), nil
}

func (m *Meta) error(err error) pb.Response {
	return shim.Error(fmt.Sprintf("meta extension err: %s", err))
}

func NewMeta(ms Metable, prefix ...string) Meta {
	if len(prefix) == 1 {
		return Meta{ms: ms, prefix: prefix[0]}
	}
	return Meta{ms: ms, prefix: DefaultPrefix}
}
