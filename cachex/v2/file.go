package cachex

import (
	"time"

	"github.com/gadelkareem/cachita"
	"github.com/pkg/errors"
	// "snippet/common"
)

type FileCache struct {
	base  *BaseCache
	store cachita.Cache
}

func NewFileCache(
	dir string, ttl, ticker time.Duration, opts ...CacheOption) (rv *FileCache) {
	var err error

	var options = DefaultCacheOptions()
	for _, opt := range opts {
		opt(options)
	}

	var store cachita.Cache
	store, err = cachita.NewFileCache(dir, ttl, ticker)
	if err != nil {
		panic(errors.Wrap(err, "create filecache"))
	}

	rv = &FileCache{
		base:  NewBaseCache(opts...),
		store: store,
	}

	return
}

func (this *FileCache) Set(k string, v interface{}, opts ...Opt) (err error) {
	var ttl time.Duration = -1

	var options = DefaultOpts()
	for _, opt := range opts {
		opt(options)
	}

	if options.Expiry != nil {
		ttl = *options.Expiry
	}

	v, err = this.base.Encode(v)
	if err != nil {
		return
	}

	err = this.store.Put(k, v, ttl)
	return
}

func (this *FileCache) Get(k string, opts ...Opt) (rv *CacheItem, err error) {
	var b []byte
	err = this.store.Get(k, &b)
	if err != nil {
		if errors.Is(err, cachita.ErrNotFound) {
			err = nil
			return
		} else if errors.Is(err, cachita.ErrExpired) {
			err = nil
			return
		} else {
			err = errors.Wrap(err, "filecache get error")
			return
		}
	}

	//	common.Logger.Debug(b)
	var v interface{}
	v, err = this.base.Decode(b)
	if err != nil {
		return
	}
	rv = &CacheItem{
		value:    v,
		encoding: this.base.Encoding,
	}
	return
}

func (this *FileCache) Del(k string, opts ...Opt) (err error) {
	err = this.store.Invalidate(k)
	return
}
