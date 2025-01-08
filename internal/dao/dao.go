package dao

import "github.com/google/wire"

var Provider = wire.NewSet(NewAccessTokenDAO, NewAccountDAO)
