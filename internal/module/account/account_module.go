package account

import "github.com/google/wire"

var Provider = wire.NewSet(NewAccountHandler)
