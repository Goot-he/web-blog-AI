package infras

import (
	"gateway/internal/infras/cache"
	"gateway/internal/infras/clog"

	"github.com/google/wire"
)

var CacheProvide = wire.NewSet(
	cache.NewCacheClient,
	cache.NewBlacklist,
)

var LoggerProvide = wire.NewSet(
	clog.NewLog,
)
