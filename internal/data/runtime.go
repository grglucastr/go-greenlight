package data

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var ErrInvalidRuntimeFormat = errors.New("invalid runtime format")

type Runtime int32

func (r Runtime) MarshalJSON() ([]byte, error) {
	jsonValue := fmt.Sprintf("%d mins", r)

	// Wrap in double quotes. This is necessary
	// to be a valid JSON String and it is a must. Because
	// we return a string.
	quotedJSONValue := strconv.Quote(jsonValue)

	return []byte(quotedJSONValue), nil

}

func (r *Runtime) UnmarshalJSON(jsonValue []byte) error {

	// The value will come as "<runtime> mins".
	// We need to remove the double quotes from it.
	// Otherwise, we'll get an undesired error.

	unquotedJSONValue, err := strconv.Unquote(string(jsonValue))
	if err != nil {
		return ErrInvalidRuntimeFormat
	}

	// Split the string to isolate the number
	parts := strings.Split(unquotedJSONValue, " ")

	if len(parts) != 2 || parts[1] != "mins" {
		return ErrInvalidRuntimeFormat
	}

	i, err := strconv.ParseInt(parts[0], 10, 32)
	if err != nil {
		return ErrInvalidRuntimeFormat
	}

	*r = Runtime(i)
	return nil
}
