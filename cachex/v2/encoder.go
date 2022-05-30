package cachex

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
	"github.com/vmihailenco/msgpack/v5"
)

func Encode(enc string, input interface{}) (rv interface{}, err error) {
	switch enc {
	case "":
		rv = input
	case ENCODING_MSGPACK:
		rv, err = msgpack.Marshal(input)
		if err != nil {
			err = errors.Wrap(err, "msgpack encode")
			return
		}
	case ENCODING_JSON:
		rv, err = json.Marshal(input)
		if err != nil {
			err = errors.Wrap(err, "json encode")
			return
		}
	default:
		err = fmt.Errorf("unsupported encoding: %s", enc)
	}
	return
}
