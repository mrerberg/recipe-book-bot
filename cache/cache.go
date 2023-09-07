package cache

import "github.com/bradfitz/gomemcache/memcache"

type Cache interface {
	Get(key string) (item *memcache.Item, err error)
	Set(item *memcache.Item) error
}
