package cachex

import (
	"fmt"
	"time"
	// "snippet/common"
)

type MemoryCache struct {
	base  *BaseCache
	store map[string]CacheItem
}

func NewMemoryCache(opts ...CacheOption) *MemoryCache {

	var options = DefaultCacheOptions()

	for _, opt := range opts {
		opt(options)
	}

	var rv = &MemoryCache{
		base:  NewBaseCache(opts...),
		store: make(map[string]CacheItem),
	}

	return rv
}

func (this *MemoryCache) Set(k string, v interface{}, opts ...Opt) (err error) {

	var options = DefaultOpts()

	for _, opt := range opts {
		opt(options)
	}

	var ex *time.Time
	if options.Expiry != nil {
		var t = time.Now().Add(*options.Expiry)
		ex = &t
	}

	v, err = this.base.Encode(v)
	if err != nil {
		err = fmt.Errorf("encode error: %w", err)
		return
	}

	this.store[k] = CacheItem{value: v, encoding: this.base.Encoding, expiry: ex}
	return
}

func (this *MemoryCache) Get(k string, opts ...Opt) (rv *CacheItem, err error) {
	var item CacheItem
	item, ok := this.store[k]
	if !ok {
		return
	}
	if item.expiry != nil && time.Since(*item.expiry) > 0 {
		return
	}

	var v interface{}
	v, err = this.base.Decode(item.value)
	if err != nil {
		err = fmt.Errorf("decode error: %w", err)
		return
	}
	item.value = v
	rv = &item
	return
}

func (this *MemoryCache) Delete(k string, opts ...Opt) (err error) {
	delete(this.store, k)
	return
}
