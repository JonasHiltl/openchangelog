package xcache

import (
	"github.com/gregjones/httpcache"
	"github.com/gregjones/httpcache/diskcache"
	"github.com/jonashiltl/openchangelog/internal/config"
	"github.com/peterbourgon/diskv"
	"github.com/sourcegraph/s3cache"
)

type Cache = httpcache.Cache

func NewS3Cache(bucket string) httpcache.Cache {
	return s3cache.New(bucket)
}

func NewDiskCache(cfg config.Config) httpcache.Cache {
	return diskcache.NewWithDiskv(diskv.New(diskv.Options{
		BasePath:     cfg.Cache.Disk.Location,
		CacheSizeMax: cfg.Cache.Disk.MaxSize, // bytes
	}))
}

func NewMemoryCache() httpcache.Cache {
	return httpcache.NewMemoryCache()
}
