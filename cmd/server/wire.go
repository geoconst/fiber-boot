//go:build wireinject
// +build wireinject

package main

import (
	"fiber-boot/internal/app"
	"fiber-boot/internal/dao"
	"fiber-boot/internal/module/account"

	"github.com/google/wire"
)

func InitializeServer() *app.Server {
	wire.Build(dao.Provider, account.Provider, app.Provider)
	return &app.Server{}
}
