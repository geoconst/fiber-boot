// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"fiber-boot/internal/app"
	"fiber-boot/internal/dao"
	"fiber-boot/internal/module/account"
	"github.com/google/wire"
)

// Injectors from wire.go:

func InitializeServer() *app.Server {
	config := app.NewConfig()
	db := app.NewDB(config)
	accountDAO := dao.NewAccountDAO(db)
	accessTokenDAO := dao.NewAccessTokenDAO(db)
	accountHandler := account.NewAccountHandler(accountDAO, accessTokenDAO)
	server := app.NewServer(config, accountHandler, accessTokenDAO)
	return server
}

// wire.go:

var daos = wire.NewSet(app.NewConfig, app.NewDB, dao.NewAccountDAO, dao.NewAccessTokenDAO)

var handlers = wire.NewSet(account.NewAccountHandler)