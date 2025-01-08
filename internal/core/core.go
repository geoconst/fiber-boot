package core

import "github.com/google/wire"

var Provider = wire.NewSet(NewConfig, NewDB)
