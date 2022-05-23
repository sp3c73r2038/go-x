package cachex

import (
	"encoding/json"
	"fmt"

	"github.com/vmihailenco/msgpack/v5"
)

func Encode(enc string, input interface{}) (rv interface{}, err error) {
	switch enc {
	case "":
		rv = input
	case ENCODING_MSGPACK:
		rv, err = msgpack.Marshal(input)
		if err != nil {
			err = fmt.Errorf("msgpack encode : %w", err)
			return
		}
	case ENCODING_JSON:
		rv, err = json.Marshal(input)
		if err != nil {
			err = fmt.Errorf("json encode: %w", err)
			return
		}
	default:
		err = fmt.Errorf("unsupported encoding: %s", enc)
	}
	return
}
