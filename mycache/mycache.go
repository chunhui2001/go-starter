package mycache

import (
	"time"

	"github.com/patrickmn/go-cache"
)

var ShortIdCached *cache.Cache = cache.New(5*time.Minute, 10*time.Minute)

func ShortIdPut(shortid string) string {
	ShortIdCached.Set(shortid, nil, cache.NoExpiration)
	return shortid
}

func ShortIdExists(shortid string) bool {
	_, found := ShortIdCached.Get(shortid)
	return found
}
