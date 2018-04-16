package helpers

import (
	"github.com/labstack/echo"
	"s7ab-platform-hyperledger/platform/core/api/common"
	"s7ab-platform-hyperledger/platform/core/logger"
)

type Context struct {
	common.Context
	SDK *MemberSDK
	BankSDK *BankSDK
}

type BankContext struct {
	common.Context
	SDK *BankSDK
}

func NewContext(e echo.Context, s *MemberSDK, l logger.Logger) Context {
	c := Context{}
	c.SDK = s
	c.Context = common.NewContext(e, s.SDKCore, l)
	return c
}

func NewBankContext(e echo.Context, s *BankSDK, l logger.Logger) BankContext {
	c := BankContext{}
	c.SDK = s
	c.Context = common.NewContext(e, s.SDKCore, l)
	return c
}
