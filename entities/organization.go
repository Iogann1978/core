package entities

import (
	"fmt"
)

const ORG_KEY = "ORGANIZATIONS"

type OrganizationType string

const MEMBER_TYPE OrganizationType = "MEMBER"
const BANK_TYPE OrganizationType = "BANK"

type OrganizationCredentials struct {
	OrganizationId     string `json:"organization_id"`
	OrganizationCert   string `json:"organization_cert"`
	OrganizationCACert string `json:"organization_ca_cert"`
}

type Organization struct {
	OrganizationCredentials
	Type       OrganizationType `json:"type"`
	Requisites Requisites       `json:"requisites"`
}

type Bank struct {
	Organization
}

type Member struct {
	Organization
	BankOrganizationId string `json:"bank_organization_id"`
	ConfirmedByBank    bool   `json:"confirmed_by_bank"`
}

type Requisites struct {
	// Название организации
	Name string `json:"name"`
	// Юридический адрес
	LegalAddress string `json:"legal_address"`
	// Фактический адрес
	ActualAddress string `json:"actual_address"`
	// ИНН
	ITN string `json:"itn"`
	// КПП
	IEC string `json:"iec"`
	// БИК
	BIC string `json:"bic"`
	// Расчетный счет
	SettlementAccount string `json:"settlement_account"`
	// Корреспондентский счет
	CorrespondentAccount string `json:"correspondent_account"`
}

func (o *Organization) Key() string {
	return GetOrganizationKey(o.OrganizationId, o.Type)
}

type OrganizationRequest struct {
	OrganizationCredentials
	Approves map[string]interface{} `json:"approves"`
}

func GetOrganizationKey(id interface{}, _type interface{}) string {
	return fmt.Sprintf("%s_%s_%s", ORG_KEY, _type, id)
}
