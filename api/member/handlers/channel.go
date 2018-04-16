package handlers

import (
	"encoding/json"
	"github.com/labstack/echo"
	"io/ioutil"
	"s7ab-platform-hyperledger/platform/core/api/common"
	me "s7ab-platform-hyperledger/platform/core/api/member/entities"
	"s7ab-platform-hyperledger/platform/core/api/member/helpers"
	"s7ab-platform-hyperledger/platform/core/entities"
	"s7ab-platform-hyperledger/platform/core/logger"
	"sync"
	"time"
)

func ChannelListHandler(c echo.Context) error {
	ctx, ok := c.(helpers.Context)
	if !ok {
		return common.ErrInvalidContext
	}

	response := make([]entities.ChannelListResponse, 0)

	res, err := ctx.SDK.Client.QueryChannels(ctx.SDK.Peer)
	if err != nil {
		return ctx.WriteError(err)
	}
	var wg, secondWg sync.WaitGroup
	outChannel := make(chan entities.ChannelListResponse, len(res.Channels))

	secondWg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		for clResponse := range outChannel {
			response = append(response, clResponse)
		}
	}(&secondWg)

	for _, c := range res.Channels {
		wg.Add(1)
		go retrieveChannelConfigBlock(&wg, ctx.SDK, ctx.Log, c.ChannelId, outChannel)
	}

	wg.Wait()
	close(outChannel)

	secondWg.Wait()

	return ctx.WriteSuccess(response)
}

func retrieveChannelConfigBlock(wg *sync.WaitGroup, sdk *helpers.MemberSDK, l logger.Logger, ch string, outChannel chan<- entities.ChannelListResponse) {
	defer wg.Done()
	if chConfig, err := sdk.GetChannelConfigBlock(ch); err != nil {
		l.Warn(`GetChannelConfigBlock`, logger.KV(`channel`, err.Error()))
		return
	} else {
		outChannel <- entities.ChannelListResponse{
			ChannelName:   ch,
			ChannelConfig: chConfig,
		}
	}
}

func ChannelJoinHandler(c echo.Context) error {
	ctx, ok := c.(helpers.Context)
	if !ok {
		return common.ErrInvalidContext
	}
	body, err := ioutil.ReadAll(ctx.Request().Body)
	if err != nil {
		return ctx.WriteError(err)
	}

	defer ctx.Request().Body.Close()

	var req me.JoinChannelRequest

	if err = json.Unmarshal(body, &req); err != nil {
		return ctx.WriteError(err)
	}

	if err = ctx.SDK.JoinChannel(req.ChannelName); err != nil {
		return ctx.WriteError(err)
	}

	time.Sleep(3 * time.Second)

	return ctx.WriteSuccess(true)
}
