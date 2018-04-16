package util

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	cb "github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/protos/common"
)

func UnmarshalChannelHeader(bytes []byte) (*cb.ChannelHeader, error) {
	chdr := &cb.ChannelHeader{}
	err := proto.Unmarshal(bytes, chdr)
	if err != nil {
		return nil, fmt.Errorf("UnmarshalChannelHeader failed, err %s", err)
	}

	return chdr, nil
}
