package cachex

import (
	"encoding/json"
	"time"

	"github.com/pkg/errors"
	"github.com/vmihailenco/msgpack/v5"
)

type CacheItem struct {
	value    interface{}
	expiry   *time.Time
	encoding string
}

func (this *CacheItem) decode(rv interface{}) (err error) {
	switch this.encoding {
	case "":
		rv = this.value
	default:
		var b []byte
		switch t := this.value.(type) {
		case []byte:
			b = this.value.([]byte)
		default:
			err = errors.Wrapf(err, "expect item value to be []byte, but got %s", t)
			return
		}
		switch this.encoding {
		case ENCODING_MSGPACK:
			err = msgpack.Unmarshal(b, &rv)
			if err != nil {
				return
			}
		case ENCODING_JSON:
			err = json.Unmarshal(b, &rv)
			if err != nil {
				err = errors.Wrap(err, "json decode")
				return
			}
		default:
			err = errors.Wrapf(err, "unsupported encoding: %s", this.encoding)
		}
	}
	return
}

func (this *CacheItem) Object(rv interface{}) (err error) {
	return this.decode(&rv)
}

func (this *CacheItem) Int() int {
	var err error
	var rv int
	err = this.decode(&rv)
	if err != nil {
		panic(err)
	}
	return rv
}

func (this *CacheItem) Int64() int64 {
	var err error
	var rv int64
	err = this.decode(&rv)
	if err != nil {
		panic(err)
	}
	return rv
}

func (this *CacheItem) Float64() float64 {
	var err error
	var rv float64
	err = this.decode(&rv)
	if err != nil {
		panic(err)
	}
	return rv
}

func (this *CacheItem) String() string {
	var err error
	var rv string
	err = this.decode(&rv)
	if err != nil {
		panic(err)
	}
	return rv
}

func (this *CacheItem) Bool() bool {
	var err error
	var rv bool
	err = this.decode(&rv)
	if err != nil {
		panic(err)
	}
	return rv
}
