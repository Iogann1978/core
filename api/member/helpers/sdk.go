package helpers

import (
	"encoding/json"
	"s7ab-platform-hyperledger/platform/core"
	"s7ab-platform-hyperledger/platform/core/entities"
	"s7ab-platform-hyperledger/platform/core/logger"
	entities2 "s7ab-platform-hyperledger/platform/s7ticket/entities"
	entities3 "s7ab-platform-hyperledger/platform/s7ticket/api/tickets/entities"
)

const (
	ORGANIZATIONS_CHAINCODE = `organizations`
)

type MemberSDK struct {
	*core.SDKCore
}

func (s *MemberSDK) GetMembersByBank() ([]entities.Member, error) {
	res, err := s.Query(ORGANIZATIONS_CHAINCODE, "/bank/member/list", []string{})

	if err != nil {
		return nil, err
	}

	if res == nil {
		return []entities.Member{}, nil
	}

	var result []entities.Member
	err = json.Unmarshal(res, &result)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *MemberSDK) ConfirmMemberByBank(memberId string) error {
	_, err := s.Invoke(ORGANIZATIONS_CHAINCODE, "/bank/member/confirm", []string{memberId})

	return err
}

func (s *MemberSDK) UnconfirmMemberByBank(memberId string) error {
	_, err := s.Invoke(ORGANIZATIONS_CHAINCODE, "/bank/member/unconfirm", []string{memberId})

	return err
}

func InitMemberSDK(org string, channel string, l logger.Logger) (*MemberSDK, error) {
	if s, err := core.Init(org, channel, l); err != nil {
		return nil, err
	} else {
		return &MemberSDK{s}, nil
	}
}

type BankSDK struct {
	*MemberSDK
}

func (s *BankSDK) SetPaymentCheckFundsInProgress(paymentId string) error {

	payload := &entities3.RequestUpdateState{
		PaymentId: paymentId,
		State: entities2.CheckFundsInProgress,
	}

	payloadBytes, err := json.Marshal(payload)

	_, err = s.Invoke("tickets", "/updateState", []string{string(payloadBytes)})
	if err != nil {
		return err
	}

	return nil
}

func (s *BankSDK) SetPaymentCheckFundsSuccess(paymentId string) error {

	payload := &entities3.RequestUpdateState{
		PaymentId: paymentId,
		State: entities2.CheckFundsSuccess,
	}

	payloadBytes, err := json.Marshal(payload)

	_, err = s.Invoke("tickets", "/updateState", []string{string(payloadBytes)})
	if err != nil {
		return err
	}

	return nil
}

func (s *BankSDK) SetPaymentPaymentCheckFundsFail(paymentId string) error {

	payload := &entities3.RequestUpdateState{
		PaymentId: paymentId,
		State: entities2.CheckFundsFail,
	}

	payloadBytes, err := json.Marshal(payload)

	_, err = s.Invoke("tickets", "/updateState", []string{string(payloadBytes)})
	if err != nil {
		return err
	}

	return nil
}

func (s *BankSDK) SetPaymentPaymentDebitInProgress(paymentId string) error {

	payload := &entities3.RequestUpdateState{
		PaymentId: paymentId,
		State: entities2.DebitInProgress,
	}

	payloadBytes, err := json.Marshal(payload)

	_, err = s.Invoke("tickets", "/updateState", []string{string(payloadBytes)})
	if err != nil {
		return err
	}

	return nil
}

func (s *BankSDK) SetPaymentPaymentDebitSuccess(paymentId string) error {

	payload := &entities3.RequestUpdateState{
		PaymentId: paymentId,
		State: entities2.DebitSuccess,
	}

	payloadBytes, err := json.Marshal(payload)

	_, err = s.Invoke("tickets", "/updateState", []string{string(payloadBytes)})
	if err != nil {
		return err
	}

	return nil
}

func (s *BankSDK) SetPaymentPaymentDebitFail(paymentId string) error {

	payload := &entities3.RequestUpdateState{
		PaymentId: paymentId,
		State: entities2.DebitFail,
	}

	payloadBytes, err := json.Marshal(payload)

	_, err = s.Invoke("tickets", "/updateState", []string{string(payloadBytes)})
	if err != nil {
		return err
	}

	return nil
}

func InitBankSDK(org string, channel string, l logger.Logger) (*BankSDK, error) {
	msdk, err := InitMemberSDK(org, channel, l)
	if err != nil {
		return nil, err
	}

	return &BankSDK{
		msdk,
	}, nil
}
