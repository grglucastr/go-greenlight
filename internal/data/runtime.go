package data

import (
	"fmt"
	"strconv"
)

type Runtime int32

func (r Runtime) MarshalJSON() ([]byte, error) {
	jsonValue := fmt.Sprintf("%d mins", r)

	// Wrap in double quotes. This is necessary
	// to be a valid JSON String and it is a must. Because
	// we return a string.
	quotedJSONValue := strconv.Quote(jsonValue)

	return []byte(quotedJSONValue), nil

}
