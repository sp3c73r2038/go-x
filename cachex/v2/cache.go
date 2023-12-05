package cachex

import (
	"fmt"
	"time"
	// "github.com/sp3c73r2038/go-x/common"
)

const ENCODING_JSON = "json"
const ENCODING_MSGPACK = "msgpack"

type Cache interface {
	Set(string, interface{}, ...Opt) error
	Get(string, ...Opt) (*CacheItem, error)
	Del(string, ...Opt) error
}

type CacheOptions struct {
	Encoding      string
	Encryption    bool
	EncryptionKey []byte
}

func DefaultCacheOptions() *CacheOptions {
	return &CacheOptions{
		Encoding:   "",
		Encryption: false,
	}
}

type CacheOption = func(*CacheOptions)

func WithEncoding(enc string) func(*CacheOptions) {
	return func(opts *CacheOptions) {

		switch enc {
		case ENCODING_JSON:
		case ENCODING_MSGPACK:
		default:
			panic(fmt.Errorf("unsupported encoding: %s", enc))
		}

		opts.Encoding = enc
	}
}

func WithEncryption(key []byte) func(*CacheOptions) {
	return func(opts *CacheOptions) {
		if len(key) <= 16 {
			panic(fmt.Errorf("key must be large than 16 bytes"))
		}

		opts.Encryption = true
		opts.EncryptionKey = key
	}
}

type Opts struct {
	Expiry *time.Duration
}

func DefaultOpts() *Opts {
	return &Opts{}
}

type Opt = func(*Opts)

func WithExpiry(ex time.Duration) func(*Opts) {
	return func(opts *Opts) {
		opts.Expiry = &ex
	}
}
