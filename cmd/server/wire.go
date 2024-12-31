//go:build wireinject
// +build wireinject

package main

import (
	"fiber-boot/internal/app"
	"fiber-boot/internal/dao"
	"fiber-boot/internal/module/account"

	"github.com/google/wire"
)

var daos = wire.NewSet(
	app.NewConfig,
	app.NewDB,
	dao.NewAccountDAO,
	dao.NewAccessTokenDAO,
)

var handlers = wire.NewSet(
	account.NewAccountHandler,
)

func InitializeServer() *app.Server {
	wire.Build(daos, handlers, app.NewServer)
	return &app.Server{}
}
