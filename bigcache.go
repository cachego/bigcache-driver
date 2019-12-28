package bigcache

import (
	"encoding/json"
	"time"

	"github.com/allegro/bigcache"
	"github.com/easy-cache/cache"
)

type bigCacheDriver struct {
	bigcache *bigcache.BigCache
}

func (bcd bigCacheDriver) Get(key string) ([]byte, bool, error) {
	bs, err := bcd.bigcache.Get(key)
	if err != nil {
		if err == bigcache.ErrEntryNotFound {
			return nil, false, nil
		}
		return nil, false, err
	}
	var item cache.Item
	if err = json.Unmarshal(bs, &item); err != nil {
		return nil, false, err
	}
	val, ok := item.GetValue()
	return val, ok, err
}

func (bcd bigCacheDriver) Set(key string, val []byte, ttl time.Duration) error {
	item := cache.NewItem(val, ttl)
	bs, err := json.Marshal(item)
	if err == nil {
		err = bcd.bigcache.Set(key, bs)
	}
	return err
}

func (bcd bigCacheDriver) Del(key string) error {
	return bcd.bigcache.Delete(key)
}

func (bcd bigCacheDriver) Has(key string) (bool, error) {
	_, ok, err := bcd.Get(key)
	return ok, err
}

func NewDriver(bc *bigcache.BigCache) cache.DriverInterface {
	return bigCacheDriver{bigcache: bc}
}

func NewCache(bc *bigcache.BigCache, args ...interface{}) *cache.Cache {
	return cache.New(append(args, NewDriver(bc))...)
}