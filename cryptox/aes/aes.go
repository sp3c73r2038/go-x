package aes

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha1"
	"errors"
	"io"

	"golang.org/x/crypto/pbkdf2"

	"github.com/sp3c73r2038/go-x/common"
)

type AESKey struct {
	Key    []byte
	Salt   []byte
	Iter   int
	Keylen int
}

func (this *AESKey) Bytes() (rv []byte) {
	return AESKeyBytes(this.Key, this.Salt, this.Iter, this.Keylen)
}

func (this *AESKey) Encrypt(pl []byte) (rv []byte, err error) {
	return Encrypt(this.Bytes(), pl)
}

func (this *AESKey) Decrypt(pl []byte) (rv []byte, err error) {
	return Decrypt(this.Bytes(), pl)
}

func CreateAESKey(key, salt []byte, keylen int) *AESKey {
	return &AESKey{
		Key:    key,
		Salt:   salt,
		Iter:   4096,
		Keylen: keylen,
	}
}

func GenAESKey(keylen int) *AESKey {
	return CreateAESKey(RandBytes(16), RandBytes(16), keylen)
}

var (
	// ErrInvalidBlockSize indicates hash blocksize <= 0.
	ErrInvalidBlockSize = errors.New("invalid blocksize")

	// ErrInvalidPKCS7Data indicates bad input to PKCS7 pad or unpad.
	ErrInvalidPKCS7Data = errors.New("invalid PKCS7 data (empty or not padded)")

	// ErrInvalidPKCS7Padding indicates PKCS7 unpad fails to bad input.
	ErrInvalidPKCS7Padding = errors.New("invalid padding on input")
)

func Encrypt(key []byte, input []byte) (rv []byte, err error) {

	// Create the AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	var padded, _ = PKCS7Pad(input, block.BlockSize())
	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	rv = make([]byte, aes.BlockSize+len(padded))
	iv := rv[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}
	common.Logger.Debug("encrypt iv: ", iv)
	bm := cipher.NewCBCEncrypter(block, iv)
	bm.CryptBlocks(rv[aes.BlockSize:], padded)
	return
}

func Decrypt(key []byte, input []byte) (rv []byte, err error) {
	// Create the AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	iv := input[:aes.BlockSize]
	common.Logger.Debug("decrypt iv: ", iv)
	input = input[aes.BlockSize:]
	bm := cipher.NewCBCDecrypter(block, iv)
	bm.CryptBlocks(input, input)
	rv, _ = PKCS7Unpad(input, aes.BlockSize)
	return
}

// func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
// 	padding := blockSize - len(ciphertext)%blockSize
// 	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
// 	return append(ciphertext, padtext...)
// }
//
// func PKCS5Trimming(encrypt []byte) []byte {
// 	padding := encrypt[len(encrypt)-1]
// 	return encrypt[:len(encrypt)-int(padding)]
// }

// pkcs7Pad right-pads the given byte slice with 1 to n bytes, where
// n is the block size. The size of the result is x times n, where x
// is at least 1.
func PKCS7Pad(b []byte, blocksize int) ([]byte, error) {
	if blocksize <= 0 {
		return nil, ErrInvalidBlockSize
	}
	if b == nil || len(b) == 0 {
		return nil, ErrInvalidPKCS7Data
	}
	n := blocksize - (len(b) % blocksize)
	pb := make([]byte, len(b)+n)
	copy(pb, b)
	copy(pb[len(b):], bytes.Repeat([]byte{byte(n)}, n))
	return pb, nil
}

// pkcs7Unpad validates and unpads data from the given bytes slice.
// The returned value will be 1 to n bytes smaller depending on the
// amount of padding, where n is the block size.
func PKCS7Unpad(b []byte, blocksize int) ([]byte, error) {
	if blocksize <= 0 {
		return nil, ErrInvalidBlockSize
	}
	if b == nil || len(b) == 0 {
		return nil, ErrInvalidPKCS7Data
	}
	if len(b)%blocksize != 0 {
		return nil, ErrInvalidPKCS7Padding
	}
	c := b[len(b)-1]
	n := int(c)
	if n == 0 || n > len(b) {
		return nil, ErrInvalidPKCS7Padding
	}
	for i := 0; i < n; i++ {
		if b[len(b)-n+i] != c {
			return nil, ErrInvalidPKCS7Padding
		}
	}
	return b[:len(b)-n], nil
}

func RandBytes(l int) []byte {
	var rv = make([]byte, l)
	rand.Read(rv)
	return rv
}

// padding key to aes qualified key length (16/32/48)
func AESKeyBytes(key, salt []byte, iter, keylen int) (rv []byte) {
	rv = pbkdf2.Key(key, salt, iter, keylen, sha1.New)
	return
}
