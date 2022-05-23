package aes

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
)

func GCMEncrypt(key, input, additional []byte) (rv []byte, err error) {
	var block cipher.Block
	block, err = aes.NewCipher(key)
	if err != nil {
		err = fmt.Errorf("create cipher: %w", err)
		return
	}

	var nonce = RandBytes(12)

	var c cipher.AEAD
	c, err = cipher.NewGCM(block)
	if err != nil {
		err = fmt.Errorf("create gcm: %w", err)
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
		err = fmt.Errorf("create cipher: %w", err)
		return
	}

	var c cipher.AEAD
	c, err = cipher.NewGCM(block)
	if err != nil {
		err = fmt.Errorf("create gcm: %w", err)
		return
	}

	rv, err = c.Open(nil, nonce, encrypted, additional)
	return
}
