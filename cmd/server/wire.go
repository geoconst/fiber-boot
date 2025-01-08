//go:build wireinject
// +build wireinject

package main

import (
	"fiber-boot/internal/app"
	"fiber-boot/internal/core"
	"fiber-boot/internal/dao"
	"fiber-boot/internal/module/account"

	"github.com/google/wire"
)

func InitApp() *app.App {
	wire.Build(core.Provider, dao.Provider, account.Provider, app.Provider)
	return &app.App{}
}
