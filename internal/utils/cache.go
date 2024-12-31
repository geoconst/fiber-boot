package utils

import (
	"time"

	"github.com/patrickmn/go-cache"
)

var accessTokenCache = cache.New(2*time.Hour, 5*time.Minute)

func AccessTokenCache() *cache.Cache {
	return accessTokenCache
}
