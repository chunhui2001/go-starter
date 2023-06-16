package cache

import (
	"time"

	"github.com/dgraph-io/ristretto"
)

var (
	MY_CACHE *ristretto.Cache
)

func init() {

	go_cache, error := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,     // Num keys to track frequency of (10M). (key 跟踪频率)
		MaxCost:     1 << 30, // Maximum cost of cache (1GB). (缓存的最大成本)
		BufferItems: 64,      // Number of keys per Get buffer. (每个 Get buffer的 key 数)
	})

	if error != nil {
		panic(error)
	}

	MY_CACHE = go_cache

}

func PutCache(key string, val any) {
	MY_CACHE.Set(key, val, 1)
}

func PutCacheWithTTL(key string, val any, ttl time.Duration) {
	MY_CACHE.SetWithTTL(key, val, 1, ttl)
}

func CacheGet(key string) any {
	value, found := MY_CACHE.Get(key)
	if !found {
		return nil
	}
	return value
}

func CacheTTl(key string) time.Duration {
	timeDuration, _ := MY_CACHE.GetTTL(key)
	return timeDuration
}
