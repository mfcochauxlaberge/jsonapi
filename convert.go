package jsonapi

import (
	"strconv"
)

// Convert ...
func Convert(src interface{}, dst interface{}, typ int) interface{} {
	// The job is to convert to a string and then back to the type
	// represented by the attribute type (typ).

	var str string
	switch s := src.(type) {
	case string:
		str = s
	case int:
		str = strconv.FormatInt(int64(s), 10)
	case int8:
		str = strconv.FormatInt(int64(s), 10)
	case int16:
		str = strconv.FormatInt(int64(s), 10)
	case int32:
		str = strconv.FormatInt(int64(s), 10)
	case int64:
		str = strconv.FormatInt(s, 10)
	case uint:
		str = strconv.FormatUint(uint64(s), 10)
	case uint8:
		str = strconv.FormatUint(uint64(s), 10)
	case uint16:
		str = strconv.FormatUint(uint64(s), 10)
	case uint32:
		str = strconv.FormatUint(uint64(s), 10)
	case uint64:
		str = strconv.FormatUint(s, 10)
	}

	return str
}
