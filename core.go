package core

import (
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/api/apifabca"
	"github.com/hyperledger/fabric-sdk-go/api/apifabclient"
	"github.com/hyperledger/fabric-sdk-go/api/apitxn"
	"github.com/hyperledger/fabric-sdk-go/def/fabapi"
	"github.com/hyperledger/fabric-sdk-go/def/fabapi/opt"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabric-client/events"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabric-client/orderer"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabric-txn/chclient"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/protos/common"
	log "github.com/sirupsen/logrus"
	"os"
	"s7ab-platform-hyperledger/platform/core/configtxlator"
	"s7ab-platform-hyperledger/platform/core/logger"
	"s7ab-platform-hyperledger/platform/core/utils"
	"strings"
)

type SDKCore struct {
	Client        apifabclient.FabricClient
	EventHub      *events.EventHub
	Channel       apifabclient.Channel
	Peer          apifabclient.Peer
	SDK           *fabapi.FabricSDK
	Context       *fabapi.OrgContext
	AdminUser     *apifabca.User
	Orderer       apifabclient.Orderer
	ChannelClient *chclient.ChannelClient
}

const USER_STATE_PATH = "/tmp/enroll_user"

func NewSDK() (*SDKCore, error) {

	var config string

	//if dev := os.Getenv("DEV_MODE"); dev != "" {
	//	config = os.Getenv("GOPATH") + "/src/s7ab-platform-hyperledger/config.local.yaml"
	//} else {
	//	config = os.Getenv("GOPATH") + "/src/s7ab-platform-hyperledger/config.yaml"
	//}

	if config = os.Getenv("CONFIG_PATH"); config == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	options := fabapi.Options{
		ConfigFile: config,
		StateStoreOpts: opt.StateStoreOpts{
			Path: USER_STATE_PATH,
		},
	}

	sdk, err := fabapi.NewSDK(options)

	if err != nil {
		return nil, err
	}

	return &SDKCore{
		SDK: sdk,
	}, nil
}

func Init(org string, channel string, l logger.Logger) (*SDKCore, error) {
	s, err := NewSDK()
	if err != nil {
		return nil, err
	}

	client, err := s.CreatePeerAdmin(org)
	if err != nil {
		return nil, err
	}

	l.Debug(`InitSDK`, logger.KV(`org`, org))

	if err = s.SetUser(&client); err != nil {
		return nil, err
	}

	if err := s.Initialize(channel); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *SDKCore) SetUser(user *apifabclient.FabricClient) error {

	s.Client = *user

	log.Info("[SDK] SetUser MspID: ", s.Client.UserContext().MspID())

	orgName := strings.ToLower(s.Client.UserContext().MspID())

	context, err := s.SDK.NewContext(orgName)

	if err != nil {
		return err
	}

	s.Context = context

	// If CA service not provided in configuration, don't initializing admin user
	//caconfig, err := s.SDK.ConfigProvider().CAConfig(strings.ToLower(s.Client.UserContext().MspID()))
	//if err != nil {
	//	return err
	//}

	//if caconfig.CAName != "" {
	//	// TODO: remove hardcoded bootstrap admin ca credentials
	//	adminCAUser, err := fabapi.NewUser(s.SDK.ConfigProvider(), s.Context.MSPClient(),
	//		caconfig.Registrar.EnrollID, caconfig.Registrar.EnrollSecret, s.Client.UserContext().MspID())
	//
	//	if err != nil {
	//		return err
	//	}
	//
	//	s.AdminUser = &adminCAUser
	//}

	return nil
}

func (s *SDKCore) Initialize(channelName string) error {

	if s.Client == nil {
		return fmt.Errorf("Set User before use this")
	}

	ordererConfig, err := s.SDK.ConfigProvider().OrdererConfig("orderer.example.com")
	if err != nil {
		return err
	}

	localOrderer, err := orderer.NewOrderer(ordererConfig.URL, "", "", s.SDK.ConfigProvider())
	ch, err := s.Client.NewChannel(channelName)
	if err != nil {
		return err
	}

	err = ch.AddOrderer(localOrderer)
	if err != nil {
		return err
	}

	peersConfig, err := s.SDK.ConfigProvider().PeersConfig(strings.ToLower(s.Client.UserContext().MspID()))
	if err != nil {
		return err
	}

	peerConfig := peersConfig[0]

	peer, err := fabapi.NewPeer(peerConfig.URL, "", "", s.SDK.ConfigProvider())

	if err != nil {
		return err
	}

	ch.AddPeer(peer)
	ch.SetPrimaryPeer(peer)

	ev, err := events.NewEventHub(s.Client)
	if err != nil {
		return err
	}

	ev.SetPeerAddr(peerConfig.EventURL, "", "")
	if err := ev.Connect(); err != nil {
		return err
	}

	s.Channel = ch

	channelClient, err := chclient.NewChannelClient(s.Client, s.Channel, nil, nil, ev)
	if err != nil {
		return err
	}

	s.Orderer = localOrderer
	s.Peer = peer
	s.EventHub = ev
	s.ChannelClient = channelClient

	return nil
}

func (s *SDKCore) NewEventHub() (*events.EventHub, error) {
	ev, err := events.NewEventHub(s.Client)
	if err != nil {
		return nil, err
	}

	peersConfig, err := s.SDK.ConfigProvider().PeersConfig(strings.ToLower(s.Client.UserContext().MspID()))
	if err != nil {
		return nil, err
	}

	peerConfig := peersConfig[0]

	ev.SetPeerAddr(peerConfig.EventURL, "", "")
	if err := ev.Connect(); err != nil {
		return nil, err
	}

	return ev, nil
}

func (s *SDKCore) Query(chaincodeName string, chaincodeFunction string, args []string) ([]byte, error) {

	queryResult, err := s.ChannelClient.QueryWithOpts(apitxn.QueryRequest{
		ChaincodeID: chaincodeName,
		Fcn:         chaincodeFunction,
		Args:        utils.ArrayToChaincodeArgs(args),
	}, apitxn.QueryOpts{
		ProposalProcessors: []apitxn.ProposalProcessor{s.Peer},
	})

	if err != nil {
		log.WithFields(log.Fields{
			"action":    "Query",
			"chaincode": chaincodeName,
			"function":  chaincodeFunction,
			"args":      args,
		}).Error(err.Error())

		return nil, err
	}

	log.WithFields(log.Fields{
		"action":    "Query",
		"chaincode": chaincodeName,
		"function":  chaincodeFunction,
		"args":      args,
	}).Info(fmt.Sprintf("%s", string(queryResult)))

	return queryResult, nil
}

func (s *SDKCore) Invoke(chaincodeName string, chaincodeFunction string, args []string) (string, error) {

	txResult, err := s.ChannelClient.ExecuteTxWithOpts(apitxn.ExecuteTxRequest{
		ChaincodeID: chaincodeName,
		Fcn:         chaincodeFunction,
		Args:        utils.ArrayToChaincodeArgs(args),
	}, apitxn.ExecuteTxOpts{
		ProposalProcessors: []apitxn.ProposalProcessor{s.Peer},
	})

	if err != nil {

		log.WithFields(log.Fields{
			"action":    "Invoke",
			"chaincode": chaincodeName,
			"function":  chaincodeFunction,
			"args":      args,
		}).Error(err.Error())

		return "", err
	}

	log.WithFields(log.Fields{
		"action":    "Invoke",
		"chaincode": chaincodeName,
		"function":  chaincodeFunction,
		"args":      args,
	}).Info(fmt.Sprintf("%v", txResult.ID))

	return txResult.ID, nil
}

func (s *SDKCore) CreatePeerAdmin(org string) (apifabclient.FabricClient, error) {

	user, err := s.SDK.NewPreEnrolledUser(org, "Admin")
	if err != nil {
		return nil, err
	}

	return fabapi.NewClient(user, false, USER_STATE_PATH, s.SDK.CryptoSuiteProvider(), s.SDK.ConfigProvider())
}

func (s *SDKCore) GetChannelGenesisBlock(channel string) (*common.Block, error) {
	tx, err := s.Client.NewTxnID()
	if err != nil {
		return nil, err
	}

	req := &apifabclient.GenesisBlockRequest{
		TxnID: tx,
	}

	var localChannel apifabclient.Channel

	if localChannel = s.Client.Channel(channel); localChannel == nil {
		localChannel, err = s.Client.NewChannel(channel)
		if err != nil {
			return nil, err
		}

		localChannel.AddOrderer(s.Orderer)
		localChannel.AddPeer(s.Peer)
		localChannel.SetPrimaryPeer(s.Peer)
	}

	return localChannel.GenesisBlock(req)
}

func (s *SDKCore) GetChannelConfigBlock(channel string) (string, error) {

	req := apitxn.ChaincodeInvokeRequest{
		Args: [][]byte{
			[]byte(channel),
		},
		Fcn:         "GetConfigBlock",
		ChaincodeID: "cscc",
	}

	result, err := s.Channel.QueryBySystemChaincode(req)
	if err != nil {
		return "", err
	}

	result1, err := configtxlator.DecodeBlock(result[0])

	return result1, nil
}

func (s *SDKCore) JoinChannel(channel_id string) error {

	block, err := s.GetChannelGenesisBlock(channel_id)
	if err != nil {
		return err
	}

	tx, err := s.Client.NewTxnID()
	if err != nil {
		return err
	}

	var ch apifabclient.Channel

	ch = s.Client.Channel(channel_id)
	if err != nil {
		return err
	}

	if ch == nil {
		ch, err = s.Client.NewChannel(channel_id)
		if err != nil {
			return err
		}
	}

	err = ch.JoinChannel(&apifabclient.JoinChannelRequest{
		TxnID:        tx,
		GenesisBlock: block,
		Targets:      []apifabclient.Peer{s.Peer},
	})

	if err != nil {
		return err
	}

	return nil
}
