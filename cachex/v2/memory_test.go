package cachex

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	// "snippet/common"
)

type testObject struct {
	Name string
	Age  int
}

func TestMemory(t *testing.T) {
	var item *CacheItem
	var err error
	// var mem = NewMemoryCache(WithEncoding(ENCODING_MSGPACK))
	var mem = NewMemoryCache(
		WithEncoding(ENCODING_JSON),
		// WithEncoding(ENCODING_MSGPACK),
		WithEncryption([]byte("12345678901234567890")),
	)

	// empty
	item, err = mem.Get("k1")
	assert.Nil(t, err)
	assert.Nil(t, item)

	// set
	err = mem.Set("k1", 1)
	assert.Nil(t, err)

	// get
	item, err = mem.Get("k1")
	assert.Nil(t, err)
	assert.NotNil(t, item)
	assert.Equal(t, 1, item.Int())

	// expiry
	err = mem.Set("k1", 1, WithExpiry(time.Millisecond*100))
	assert.Nil(t, err)
	item, err = mem.Get("k1")
	assert.Nil(t, err)
	assert.NotNil(t, item)

	time.Sleep(time.Millisecond * 100)

	item, err = mem.Get("k1")
	assert.Nil(t, err)
	assert.Nil(t, item)

	// int64
	err = mem.Set("k2", int64(20))
	assert.Nil(t, err)
	item, err = mem.Get("k2")
	assert.Nil(t, err)
	assert.Equal(t, int64(20), item.Int64())

	// float64
	err = mem.Set("k3", float64(3.2))
	assert.Nil(t, err)
	item, err = mem.Get("k3")
	assert.Nil(t, err)
	assert.Equal(t, float64(3.2), item.Float64())

	// string
	err = mem.Set("k4", "string")
	assert.Nil(t, err)
	item, err = mem.Get("k4")
	assert.Nil(t, err)
	assert.Equal(t, "string", item.String())

	// bool
	err = mem.Set("k5", true)
	assert.Nil(t, err)
	item, err = mem.Get("k5")
	assert.Nil(t, err)
	assert.Equal(t, true, item.Bool())

	// object
	err = mem.Set("k6", &testObject{"name", 10})
	assert.Nil(t, err)
	item, err = mem.Get("k6")
	assert.Nil(t, err)
	var to testObject
	err = item.Object(&to)
	assert.Nil(t, err)
	assert.NotNil(t, to)
	assert.Equal(t, "name", to.Name)
	assert.Equal(t, 10, to.Age)

	// common.Logger.Debug(common.Pretty(mem.store))
}
