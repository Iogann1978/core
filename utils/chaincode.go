package utils

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric/protos/msp"
	"strings"
)

var ErrInvalidPEMStructure = errors.New(`invalid PEM structure`)

type Creator struct {
	MspID string
	User  string
	Cert  *x509.Certificate
}

func NewCreator(b []byte) (*Creator, error) {
	identity := &msp.SerializedIdentity{}
	var err error
	if err = proto.Unmarshal(b, identity); err != nil {
		return nil, err
	}

	creator := &Creator{
		MspID: identity.Mspid,
	}

	var pb *pem.Block

	if pb, _ = pem.Decode(identity.IdBytes); pb == nil {
		return nil, ErrInvalidPEMStructure
	}

	var cert *x509.Certificate

	if cert, err = x509.ParseCertificate(pb.Bytes); err != nil {
		return nil, err
	}
	creator.Cert = cert
	possibleEmail := strings.Split(cert.Subject.CommonName, `@`)
	if len(possibleEmail) == 2 {
		creator.User = possibleEmail[0]
	} else {
		creator.User = cert.Subject.CommonName
	}

	return creator, nil
}
