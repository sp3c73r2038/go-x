package httpx

import (
	"encoding/base64"
	"fmt"
)

type Auth interface {
	AuthHeader() string
}

type BasicAuth struct {
	username string
	password string
}

func NewBasicAuth(username, password string) (rv *BasicAuth) {
	return &BasicAuth{
		username: username,
		password: password,
	}
}

func (this *BasicAuth) AuthHeader() (rv string) {
	var s = fmt.Sprintf("%s:%s", this.username, this.password)
	var v = base64.StdEncoding.EncodeToString([]byte(s))
	return fmt.Sprintf("Basic %s", v)
}

type BearerAuth struct {
	token string
}

func NewBearerAuth(token string) (rv *BearerAuth) {
	return &BearerAuth{
		token: token,
	}
}
func (this *BearerAuth) AuthHeader() (rv string) {
	return fmt.Sprintf("Bearer %s", this.token)
}

type RawAuth struct {
	value string
}

func NewRawAuth(value string) (rv *RawAuth) {
	return &RawAuth{
		value: value,
	}
}
func (this *RawAuth) AuthHeader() (rv string) {
	return this.value
}
