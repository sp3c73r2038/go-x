package aes

import (
	// "bufio"
	// "crypto/rand"
	"encoding/base64"
	// "fmt"
	// "io"
	"os"
	"testing"

	"github.com/sp3c73r2038/go-x/common"
)

func TestGCM(t *testing.T) {
	var err error
	var pass = os.Getenv("PASSWORD")
	var key = os.Getenv("KEY")
	if len(pass) <= 0 && len(key) <= 0 {
		t.Errorf("please set PASSWORD or KEY env var")
		return
	}

	var salt []byte
	if len(key) <= 0 {
		salt, _ = RandBytes(16)
		var kb = append(salt, []byte(pass)...)
		common.Logger.Debug("key: ", base64.StdEncoding.EncodeToString(kb))
	} else {
		var b []byte
		b, err = base64.StdEncoding.DecodeString(key)
		if err != nil {
			t.Errorf("base64 decode: %v", err)
			return
		}

		salt = b[:16]
		pass = string(b[16:])
	}

	var akey = CreateAESKey([]byte(pass), salt, 32)
	common.Logger.Debug(
		"aes key: ", base64.StdEncoding.EncodeToString(akey.Bytes()))

	var payload []byte
	payload, _ = RandBytes(16)
	var ps = base64.StdEncoding.EncodeToString(payload)

	common.Logger.Debug("payload: ", ps)

	var encrypted []byte
	encrypted, err = GCMEncrypt(akey.Bytes(), []byte(ps), nil)
	if err != nil {
		t.Errorf("gcm encrypt: %v", err)
		return
	}

	common.Logger.Debug(
		"encrypted: ", base64.StdEncoding.EncodeToString(encrypted))

	var decrypted []byte
	decrypted, err = GCMDecrypt(akey.Bytes(), encrypted, nil)
	if err != nil {
		t.Errorf("gcm decrypt: %v", err)
		return
	}

	common.Logger.Debug("decrypted: ", string(decrypted))
}
