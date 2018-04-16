package member

import (
	"github.com/labstack/echo"
	"s7ab-platform-hyperledger/platform/core/api/member/handlers"
)

const (
	DefaultUrlPath = `/member`
)

func NewModule(e *echo.Echo, urlPath string, m ...echo.MiddlewareFunc) {
	if urlPath == `` {
		urlPath = DefaultUrlPath
	}
	g := e.Group(urlPath, m...)
	setRouter(g)
}

func setRouter(g *echo.Group) {
	g.GET(`/chaincode`, handlers.ChaincodeListHandler)

	bankGroup := g.Group(`/bank`)
	bankMemberGroup := bankGroup.Group(`/members`)
	// Получение списка участников сети у которых в качестве банка проставлен id участника которая вызывает этот метод
	bankMemberGroup.GET(`/`, handlers.BankMemberListHandler)
	// Подтверждение участника сети банком
	bankMemberGroup.POST(`/:id/confirm`, handlers.BankMemberConfirmHandler)
	// Снятие подтверждения
	bankMemberGroup.POST(`/:id/unconfirm`, handlers.BankMemberUnconfirmHandler)

	channelGroup := g.Group(`/channel`)
	// Список каналов к которым присоеден участник
	channelGroup.GET(`/`, handlers.ChannelListHandler)
	// Присоедениться к каналу по его имени
	channelGroup.POST(`/join`, handlers.ChannelJoinHandler)

	systemGroup := g.Group(`/system`)
	systemGroup.GET(`/info`, handlers.SystemInfoHandler)
	// Получить генезис блок
	systemGroup.GET(`/genesis`, handlers.SystemGenesisHandler)
}
