package handlers

import (
	"github.com/labstack/echo"
	"s7ab-platform-hyperledger/platform/core/api/common"
	"s7ab-platform-hyperledger/platform/core/api/member/helpers"
)

func ChaincodeListHandler(c echo.Context) error {
	ctx, ok := c.(helpers.Context)
	if !ok {
		return common.ErrInvalidContext
	}
	res, err := ctx.SDK.Client.QueryInstalledChaincodes(ctx.SDK.Peer)
	if err != nil {
		return ctx.WriteError(err)
	}

	return ctx.WriteSuccess(res.Chaincodes)
}
