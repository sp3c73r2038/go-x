package cachex

import (
	"testing"
	"time"

	// "github.com/sp3c73r2038/go-x/common"
	"github.com/stretchr/testify/assert"
)

func TestFile(t *testing.T) {
	var err error
	var item *CacheItem

	var cache Cache
	var file = NewFileCache(
		"_cache", time.Hour, time.Millisecond*100,
		WithEncryption([]byte("12345678901234567890")),
	)
	cache = file
	_, _ = cache.Get("k0")

	// empty
	item, err = file.Get("k1")
	assert.Nil(t, err)
	assert.Nil(t, item)

	// set & get
	err = file.Set("k2", 1)
	assert.Nil(t, err)

	item, err = file.Get("k2")
	assert.Nil(t, err)
	assert.NotNil(t, item)
	assert.Equal(t, 1, item.Int())

	// expiry
	err = file.Set("k3", 1, WithExpiry(time.Millisecond*100))
	assert.Nil(t, err)

	item, err = file.Get("k3")
	assert.Nil(t, err)
	assert.NotNil(t, item)
	assert.Equal(t, 1, item.Int())

	time.Sleep(time.Millisecond * 100)

	item, err = file.Get("k3")
	assert.Nil(t, err)
	assert.Nil(t, item)

	// int64
	err = file.Set("k4", int64(20))
	assert.Nil(t, err)
	item, err = file.Get("k4")
	assert.Nil(t, err)
	assert.Equal(t, int64(20), item.Int64())

	// float64
	err = file.Set("k5", float64(3.2))
	assert.Nil(t, err)
	item, err = file.Get("k5")
	assert.Nil(t, err)
	assert.Equal(t, float64(3.2), item.Float64())

	// string
	err = file.Set("k6", "string")
	assert.Nil(t, err)
	item, err = file.Get("k6")
	assert.Nil(t, err)
	assert.Equal(t, "string", item.String())

	// bool
	err = file.Set("k7", true)
	assert.Nil(t, err)
	item, err = file.Get("k7")
	assert.Nil(t, err)
	assert.Equal(t, true, item.Bool())

	// object
	var obj = testObject{"name", 10}
	err = file.Set("k8", &obj)
	assert.Nil(t, err)
	item = nil
	item, err = file.Get("k8")
	assert.Nil(t, err)
	assert.NotNil(t, item)
	// common.Logger.Debug(common.Pretty(item))
	var to testObject
	err = item.Object(&to)
	assert.Nil(t, err)
	assert.NotNil(t, to)
	assert.Equal(t, "name", to.Name)
	assert.Equal(t, 10, to.Age)

}
