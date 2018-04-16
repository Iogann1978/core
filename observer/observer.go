package observer

import (
	"errors"
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/api/apifabclient"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/protos/common"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/protos/peer"
	log "github.com/sirupsen/logrus"
	"s7ab-platform-hyperledger/platform/core"
	"s7ab-platform-hyperledger/platform/core/observer/util"
)

type ChaincodeEvent struct {
	TxID        string
	Name        string
	ChaincodeID string
	ChannelID   string
	Payload     []byte
}

type BlockEvent struct {
	TxID    string
	Channel string
}

type Observer struct {
	s *core.SDKCore
}

func (o *Observer) GetChaincodeEventsByRegex(chainCode string, eventName string) <-chan ChaincodeEvent {

	ev, err := o.s.NewEventHub()
	if err != nil {
		log.Fatal(err)
	}

	ch := make(chan ChaincodeEvent)
	ev.RegisterChaincodeEvent(chainCode, eventName, func(event *apifabclient.ChaincodeEvent) {
		ch <- ChaincodeEvent{TxID: event.TxID, Name: eventName, ChaincodeID: chainCode, ChannelID: event.ChannelID, Payload: event.Payload}
	})
	return ch
}

func (o *Observer) GetChainCodeEvents() <-chan ChaincodeEvent {
	ch := make(chan ChaincodeEvent)
	o.s.EventHub.RegisterBlockEvent(func(block *common.Block) {
		event := ChaincodeEvent{}
		for i, r := range block.Data.Data {
			tx, _ := o.getTxPayload(r)
			txsFltr := util.TxValidationFlags(block.Metadata.Metadata[common.BlockMetadataIndex_TRANSACTIONS_FILTER])
			if tx != nil {
				chdr, err := util.UnmarshalChannelHeader(tx.Header.ChannelHeader)
				if err != nil {
					log.Errorln("Error extracting channel header\n")
					return
				}
				event.TxID = chdr.TxId
				event.ChannelID = chdr.ChannelId
				if !txsFltr.IsInvalid(i) {
					if e, err := o.getChainCodeEvents(r); err == nil && tx.Data != nil {
						event.Name = e.EventName
						event.ChaincodeID = e.ChaincodeId
						event.Payload = e.Payload
						ch <- event
					} else {
						log.Infof("Error get chaincode events: %s\n", err)
					}
				}
			}
		}
	})
	return ch
}

//func (o *Observer) GetBlockEvents() (<-chan BlockEvent) {
//
//}

func (o *Observer) getTxPayload(tdata []byte) (*common.Payload, error) {
	if tdata == nil {
		return nil, errors.New("Cannot extract payload from nil transaction")
	}

	if env, err := util.GetEnvelopeFromBlock(tdata); err != nil {
		return nil, fmt.Errorf("Error getting tx from block(%s)", err)
	} else if env != nil {
		// get the payload from the envelope
		payload, err := util.GetPayload(env)
		if err != nil {
			return nil, fmt.Errorf("Could not extract payload from envelope, err %s", err)
		}
		return payload, nil
	}
	return nil, nil
}

func (o *Observer) getChainCodeEvents(tdata []byte) (*peer.ChaincodeEvent, error) {
	if tdata == nil {
		return nil, errors.New("Cannot extract payload from nil transaction")
	}

	if env, err := util.GetEnvelopeFromBlock(tdata); err != nil {
		return nil, fmt.Errorf("Error getting tx from block(%s)", err)
	} else if env != nil {
		// get the payload from the envelope
		payload, err := util.GetPayload(env)
		if err != nil {
			return nil, fmt.Errorf("Could not extract payload from envelope, err %s", err)
		}

		chdr, err := util.UnmarshalChannelHeader(payload.Header.ChannelHeader)
		if err != nil {
			return nil, fmt.Errorf("Could not extract channel header from envelope, err %s", err)
		}

		if common.HeaderType(chdr.Type) == common.HeaderType_ENDORSER_TRANSACTION {
			tx, err := util.GetTransaction(payload.Data)
			if err != nil {
				return nil, fmt.Errorf("Error unmarshalling transaction payload for block event: %s", err)
			}
			chaincodeActionPayload, err := util.GetChaincodeActionPayload(tx.Actions[0].Payload)
			if err != nil {
				return nil, fmt.Errorf("Error unmarshalling transaction action payload for block event: %s", err)
			}
			propRespPayload, err := util.GetProposalResponsePayload(chaincodeActionPayload.Action.ProposalResponsePayload)
			if err != nil {
				return nil, fmt.Errorf("Error unmarshalling proposal response payload for block event: %s", err)
			}
			caPayload, err := util.GetChaincodeAction(propRespPayload.Extension)
			if err != nil {
				return nil, fmt.Errorf("Error unmarshalling chaincode action for block event: %s", err)
			}
			ccEvent, err := util.GetChaincodeEvents(caPayload.Events)

			if ccEvent != nil {
				return ccEvent, nil
			}
		}
	}
	return nil, errors.New("No events found")
}

func NewObserver(s *core.SDKCore) *Observer {
	return &Observer{s: s}
}
