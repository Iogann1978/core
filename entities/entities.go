package entities

const SYSTEM_CHANNEL_NAME = "mychannel"

const GOVERNER_ORG_NAME = "Org1MSP"

type Response struct {
	Success bool        `json:"success"`
	Error   string      `json:"error"`
	Result  interface{} `json:"result"`
}

type AddOrganizationToConfigPayload struct {
	OrganizationCredentials
	ChannelConfig string `json:"channel_config"`
	GenesisConfig string `json:"genesis_config"`
	ChannelName   string `json:"channel_name"`
}

type AddOrganizationToConfigResponse struct {
	Config    string `json:"config"`
	OldConfig string `json:"old_config"`
	Error     string `json:"error"`
}

type CreateChannelConfigPayload struct {
	OrgIds    []string `json:"organization_ids"`
	ChannelId string   `json:"channel_id"`
}

type ChannelListResponse struct {
	ChannelName   string `json:"channel_name"`
	ChannelConfig string `json:"channel_config"`
}

type KeyModification struct {
	TxID          string      `json:"tx_id"`
	Payload       []byte      `json:"payload"`
	PayloadParsed interface{} `json:"payload_parsed"`
	Time          int64       `json:"time"`
	IsDelete      bool        `json:"is_delete"`
}
