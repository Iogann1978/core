package configtxlator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"s7ab-platform-hyperledger/platform/core/entities"
)

var CONFIGTXLATOR_HOST = "127.0.0.1:7059"

func init() {
	if host := os.Getenv("CONFIGTXLATOR_HOST"); host != "" {
		CONFIGTXLATOR_HOST = host
	}
}

// ConfigUpdateEnvelope

func EncodeConfigUpdateEnvelope(blob []byte) ([]byte, error) {
	res, err := http.Post(
		"http://"+CONFIGTXLATOR_HOST+"/protolator/encode/common.ConfigUpdateEnvelope",
		"",
		bytes.NewBuffer(blob))

	if err != nil {
		return []byte{}, err
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return []byte{}, err
	}

	var js json.RawMessage
	err = json.Unmarshal(body, &js)
	if err != nil {
		return []byte{}, fmt.Errorf("configtxlator return not json result, error: %s, response: %s", err, body)
	}

	return body, nil
}

func DecodeBlock(blob []byte) (string, error) {
	//proto.Unmarshal(blob)
	//protolator.DeepMarshalJSON()
	res, err := http.Post(
		"http://"+CONFIGTXLATOR_HOST+"/protolator/decode/common.Block",
		"",
		bytes.NewBuffer(blob))

	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	var js json.RawMessage
	err = json.Unmarshal(body, &js)
	if err != nil {
		return "", fmt.Errorf("configtxlator return not json result, error: %s, response: %s", err, body)
	}

	return string(body), nil
}

func EncodeBlock(block string) ([]byte, error) {
	res, err := http.Post(
		"http://"+CONFIGTXLATOR_HOST+"/protolator/encode/common.Block",
		"",
		bytes.NewBuffer([]byte(block)))

	if err != nil {
		return []byte{}, err
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return []byte{}, err
	}

	_, err = DecodeBlock(body)
	if err != nil {
		return []byte{}, err
	}

	return body, err
}

func DecodeConfigUpdateBlock(blob []byte) (string, error) {
	res, err := http.Post(
		"http://"+CONFIGTXLATOR_HOST+"/protolator/decode/common.ConfigUpdate",
		"",
		bytes.NewBuffer(blob))

	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func EncodeConfigUpdateBlock(blob []byte) ([]byte, error) {
	res, err := http.Post(
		"http://"+CONFIGTXLATOR_HOST+"/protolator/encode/common.ConfigUpdate",
		"",
		bytes.NewBuffer(blob))

	if err != nil {
		return []byte{}, err
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return []byte{}, err
	}

	return body, nil
}

func EncodeConfigBlock(blob []byte) ([]byte, error) {

	res, err := http.Post(
		"http://"+CONFIGTXLATOR_HOST+"/protolator/encode/common.Config",
		"",
		bytes.NewBuffer(blob))

	if err != nil {
		return []byte{}, err
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return []byte{}, err
	}

	return body, nil
}

func GetDiffBeetwenChannelConfig(data *entities.AddOrganizationToConfigResponse, channel string) ([]byte, error) {

	logrus.Info("[GetDiffBeetwenChannelConfig] OldConfig: ", data.OldConfig)
	logrus.Info("[GetDiffBeetwenChannelConfig] Config: ", data.Config)

	oldBlock, err := EncodeConfigBlock([]byte(data.OldConfig))
	if err != nil {
		return []byte{}, err
	}

	newBlock, err := EncodeConfigBlock([]byte(data.Config))
	if err != nil {
		return []byte{}, err
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	wr1, err := writer.CreateFormFile("original", "config.pb")
	if err != nil {
		return []byte{}, err
	}

	wr1.Write(oldBlock)

	wr2, err := writer.CreateFormFile("updated", "updated_config.pb")
	if err != nil {
		return []byte{}, err
	}

	wr2.Write(newBlock)

	writer.WriteField("channel", channel)

	err = writer.Close()
	if err != nil {
		return []byte{}, err
	}

	req, err := http.NewRequest(
		"POST", "http://"+CONFIGTXLATOR_HOST+"/configtxlator/compute/update-from-configs",
		body)
	if err != nil {
		return []byte{}, err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return []byte{}, err
	}

	defer resp.Body.Close()

	bodyResp, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}

	bodyRespAsString, err := DecodeConfigUpdateBlock(bodyResp)
	if err != nil {
		return []byte{}, err
	}

	logrus.Info("[GetDiffBeetwenChannelConfig] UpdateBlock: ", bodyRespAsString)

	return bodyResp, nil
}
