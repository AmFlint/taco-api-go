package helpers

import (
	"encoding/json"
	"io"
)

// JsonEncode - Encode interface to []byte
func JsonEncode(v interface{}) []byte {
	marshaled, err := json.Marshal(v)

	if err != nil {
		panic(err)
	}

	return marshaled
}

// JsonDecode - Decode
func DecodeBody(body io.ReadCloser, result *map[string]interface{}) error {
	decoder := json.NewDecoder(body)
	return decoder.Decode(result)
}
