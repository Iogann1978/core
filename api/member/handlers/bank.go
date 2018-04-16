package handlers

import (
	"errors"
	"github.com/labstack/echo"
	"s7ab-platform-hyperledger/platform/core/api/common"
	"s7ab-platform-hyperledger/platform/core/api/member/helpers"
)

var ErrIdParamRequired = errors.New(`id param required`)

func BankMemberListHandler(c echo.Context) error {
	ctx, ok := c.(helpers.Context)
	if !ok {
		return common.ErrInvalidContext
	}
	if members, err := ctx.SDK.GetMembersByBank(); err != nil {
		return ctx.WriteError(err)
	} else {
		return ctx.WriteSuccess(members)
	}
}

func BankMemberConfirmHandler(c echo.Context) error {
	ctx, ok := c.(helpers.Context)
	if !ok {
		return common.ErrInvalidContext
	}
	if member := ctx.Param(`id`); member != `` {
		if err := ctx.SDK.ConfirmMemberByBank(member); err != nil {
			return ctx.WriteError(err)
		} else {
			return ctx.WriteSuccess(true)
		}
	} else {
		return ctx.WriteError(ErrIdParamRequired)
	}
}

func BankMemberUnconfirmHandler(c echo.Context) error {
	ctx, ok := c.(helpers.Context)
	if !ok {
		return common.ErrInvalidContext
	}
	if member := ctx.Param(`id`); member != `` {
		if err := ctx.SDK.UnconfirmMemberByBank(member); err != nil {
			return ctx.WriteError(err)
		} else {
			return ctx.WriteSuccess(true)
		}
	} else {
		return ctx.WriteError(ErrIdParamRequired)
	}
}
