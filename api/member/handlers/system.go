package handlers

import (
	"github.com/labstack/echo"
	"s7ab-platform-hyperledger/platform/core/api/common"
	"s7ab-platform-hyperledger/platform/core/api/member/helpers"
)

func SystemInfoHandler(c echo.Context) error {
	ctx, ok := c.(helpers.Context)
	if !ok {
		return common.ErrInvalidContext
	}
	if info, err := ctx.SDK.Channel.QueryInfo(); err != nil {
		return ctx.WriteError(err)
	} else {
		return ctx.WriteSuccess(info)
	}
}

func SystemGenesisHandler(c echo.Context) error {
	ctx, ok := c.(helpers.Context)
	if !ok {
		return common.ErrInvalidContext
	}
	if ce, err := ctx.SDK.Channel.ChannelConfig(); err != nil {
		return ctx.WriteError(err)
	} else {
		return ctx.WriteSuccess(ce)
	}
}
