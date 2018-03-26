package helpers

import "encoding/json"

func JsonEncode(v interface{}) []byte {
	marshaled, err := json.Marshal(v)

	if err != nil {
		panic(err)
	}

	return marshaled
}