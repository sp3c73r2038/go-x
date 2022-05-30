package aes

import (
	"crypto/aes"
	"crypto/cipher"

	"github.com/pkg/errors"
)

func GCMEncrypt(key, input, additional []byte) (rv []byte, err error) {
	var block cipher.Block
	block, err = aes.NewCipher(key)
	if err != nil {
		err = errors.Wrap(err, "create cipher")
		return
	}

	var nonce []byte
	nonce, err = RandBytes(12)
	if err != nil {
		return
	}

	var c cipher.AEAD
	c, err = cipher.NewGCM(block)
	if err != nil {
		err = errors.Wrap(err, "create gcm")
		return
	}

	rv = c.Seal(nil, nonce, input, additional)
	rv = append(nonce, rv...)
	return
}

func GCMDecrypt(key, input, additional []byte) (rv []byte, err error) {
	var nonce = input[:12]
	var encrypted = input[12:]

	var block cipher.Block
	block, err = aes.NewCipher(key)
	if err != nil {
		err = errors.Wrap(err, "create cipher")
		return
	}

	var c cipher.AEAD
	c, err = cipher.NewGCM(block)
	if err != nil {
		err = errors.Wrap(err, "create gcm")
		return
	}

	rv, err = c.Open(nil, nonce, encrypted, additional)
	return
}
