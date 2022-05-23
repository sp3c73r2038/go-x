package cachex

import (
	"fmt"
	// "snippet/common"

	"github.com/sp3c73r2038/go-x/cryptox/aes"
)

// base implementation for mixin
type BaseCache struct {
	Encoding   string
	Encryption bool
	Key        []byte
}

func NewBaseCache(opts ...CacheOption) *BaseCache {

	var options = DefaultCacheOptions()

	for _, opt := range opts {
		opt(options)
	}

	var rv = &BaseCache{
		Encoding:   options.Encoding,
		Encryption: options.Encryption,
	}

	// refactor: base properties like *encryptor
	// encryption
	if options.Encryption {
		var salt = options.EncryptionKey[:16]
		var k = options.EncryptionKey[16:]
		var key = aes.CreateAESKey(k, salt, 32)
		rv.Key = key.Bytes()

		// required encryption, so we must set an encoding
		if rv.Encoding == "" {
			rv.Encoding = ENCODING_MSGPACK
		}
	}

	return rv
}

func (this *BaseCache) Encode(v interface{}) (rv interface{}, err error) {
	v, err = Encode(this.Encoding, v)
	if err != nil {
		err = fmt.Errorf("encode error: %w", err)
		return
	}

	if this.Encryption {
		v, err = aes.GCMEncrypt(this.Key, v.([]byte), nil)
		if err != nil {
			err = fmt.Errorf("encrypt error: %w", err)
			return
		}
	}

	rv = v
	return
}

func (this *BaseCache) Decode(v interface{}) (rv interface{}, err error) {
	rv = v
	if this.Encryption {
		var b []byte
		b, err = aes.GCMDecrypt(this.Key, v.([]byte), nil)
		if err != nil {
			err = fmt.Errorf("decrypt error: %w", err)
			return
		}
		rv = b
	}
	return
}
