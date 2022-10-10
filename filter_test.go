package jsonapi_test

import (
	"encoding/json"
	"fmt"
	"strconv"
	"testing"
	"time"

	. "github.com/mfcochauxlaberge/jsonapi"

	"github.com/stretchr/testify/assert"
)

func TestFilterResource(t *testing.T) {
	assert := assert.New(t)

	now := time.Now()

	// Tests for attributes
	attrTests := []struct {
		rval     any
		op       string
		cval     any
		expected bool
	}{
		// string
		{rval: "bbb", op: "=", cval: "aaa", expected: false},
		{rval: "bbb", op: "=", cval: "bbb", expected: true},
		{rval: "bbb", op: "!=", cval: "aaa", expected: true},
		{rval: "bbb", op: "!=", cval: "bbb", expected: false},
		{rval: "bbb", op: "<", cval: "aaa", expected: false},
		{rval: "bbb", op: "<", cval: "bbb", expected: false},
		{rval: "bbb", op: "<", cval: "ccc", expected: true},
		{rval: "bbb", op: "<=", cval: "aaa", expected: false},
		{rval: "bbb", op: "<=", cval: "bbb", expected: true},
		{rval: "bbb", op: "<=", cval: "ccc", expected: true},
		{rval: "bbb", op: ">", cval: "aaa", expected: true},
		{rval: "bbb", op: ">", cval: "bbb", expected: false},
		{rval: "bbb", op: ">", cval: "ccc", expected: false},
		{rval: "bbb", op: ">=", cval: "aaa", expected: true},
		{rval: "bbb", op: ">=", cval: "bbb", expected: true},
		{rval: "bbb", op: ">=", cval: "ccc", expected: false},
		{rval: "aaa", op: "invalid", cval: "aaa", expected: false},

		// int
		{rval: 1, op: "=", cval: 0, expected: false},
		{rval: 1, op: "=", cval: 1, expected: true},
		{rval: 1, op: "!=", cval: 0, expected: true},
		{rval: 1, op: "!=", cval: 1, expected: false},
		{rval: 1, op: "<", cval: 0, expected: false},
		{rval: 1, op: "<", cval: 1, expected: false},
		{rval: 1, op: "<", cval: 3, expected: true},
		{rval: 1, op: "<=", cval: 0, expected: false},
		{rval: 1, op: "<=", cval: 1, expected: true},
		{rval: 1, op: "<=", cval: 3, expected: true},
		{rval: 1, op: ">", cval: 0, expected: true},
		{rval: 1, op: ">", cval: 1, expected: false},
		{rval: 1, op: ">", cval: 3, expected: false},
		{rval: 1, op: ">=", cval: 0, expected: true},
		{rval: 1, op: ">=", cval: 1, expected: true},
		{rval: 1, op: ">=", cval: 3, expected: false},
		{rval: 0, op: "invalid", cval: 0, expected: false},

		// int8
		{rval: int8(1), op: "=", cval: int8(0), expected: false},
		{rval: int8(1), op: "=", cval: int8(1), expected: true},
		{rval: int8(1), op: "!=", cval: int8(0), expected: true},
		{rval: int8(1), op: "!=", cval: int8(1), expected: false},
		{rval: int8(1), op: "<", cval: int8(0), expected: false},
		{rval: int8(1), op: "<", cval: int8(1), expected: false},
		{rval: int8(1), op: "<", cval: int8(2), expected: true},
		{rval: int8(1), op: "<=", cval: int8(0), expected: false},
		{rval: int8(1), op: "<=", cval: int8(1), expected: true},
		{rval: int8(1), op: "<=", cval: int8(2), expected: true},
		{rval: int8(1), op: ">", cval: int8(0), expected: true},
		{rval: int8(1), op: ">", cval: int8(1), expected: false},
		{rval: int8(1), op: ">", cval: int8(2), expected: false},
		{rval: int8(1), op: ">=", cval: int8(0), expected: true},
		{rval: int8(1), op: ">=", cval: int8(1), expected: true},
		{rval: int8(1), op: ">=", cval: int8(2), expected: false},
		{rval: int8(0), op: "invalid", cval: int8(0), expected: false},

		// int16
		{rval: int16(1), op: "=", cval: int16(0), expected: false},
		{rval: int16(1), op: "=", cval: int16(1), expected: true},
		{rval: int16(1), op: "!=", cval: int16(0), expected: true},
		{rval: int16(1), op: "!=", cval: int16(1), expected: false},
		{rval: int16(1), op: "<", cval: int16(0), expected: false},
		{rval: int16(1), op: "<", cval: int16(1), expected: false},
		{rval: int16(1), op: "<", cval: int16(2), expected: true},
		{rval: int16(1), op: "<=", cval: int16(0), expected: false},
		{rval: int16(1), op: "<=", cval: int16(1), expected: true},
		{rval: int16(1), op: "<=", cval: int16(2), expected: true},
		{rval: int16(1), op: ">", cval: int16(0), expected: true},
		{rval: int16(1), op: ">", cval: int16(1), expected: false},
		{rval: int16(1), op: ">", cval: int16(2), expected: false},
		{rval: int16(1), op: ">=", cval: int16(0), expected: true},
		{rval: int16(1), op: ">=", cval: int16(1), expected: true},
		{rval: int16(1), op: ">=", cval: int16(2), expected: false},
		{rval: int16(0), op: "invalid", cval: int16(0), expected: false},

		// int32
		{rval: int32(1), op: "=", cval: int32(0), expected: false},
		{rval: int32(1), op: "=", cval: int32(1), expected: true},
		{rval: int32(1), op: "!=", cval: int32(0), expected: true},
		{rval: int32(1), op: "!=", cval: int32(1), expected: false},
		{rval: int32(1), op: "<", cval: int32(0), expected: false},
		{rval: int32(1), op: "<", cval: int32(1), expected: false},
		{rval: int32(1), op: "<", cval: int32(2), expected: true},
		{rval: int32(1), op: "<=", cval: int32(0), expected: false},
		{rval: int32(1), op: "<=", cval: int32(1), expected: true},
		{rval: int32(1), op: "<=", cval: int32(2), expected: true},
		{rval: int32(1), op: ">", cval: int32(0), expected: true},
		{rval: int32(1), op: ">", cval: int32(1), expected: false},
		{rval: int32(1), op: ">", cval: int32(2), expected: false},
		{rval: int32(1), op: ">=", cval: int32(0), expected: true},
		{rval: int32(1), op: ">=", cval: int32(1), expected: true},
		{rval: int32(1), op: ">=", cval: int32(2), expected: false},
		{rval: int32(0), op: "invalid", cval: int32(0), expected: false},

		// int64
		{rval: int64(1), op: "=", cval: int64(0), expected: false},
		{rval: int64(1), op: "=", cval: int64(1), expected: true},
		{rval: int64(1), op: "!=", cval: int64(0), expected: true},
		{rval: int64(1), op: "!=", cval: int64(1), expected: false},
		{rval: int64(1), op: "<", cval: int64(0), expected: false},
		{rval: int64(1), op: "<", cval: int64(1), expected: false},
		{rval: int64(1), op: "<", cval: int64(2), expected: true},
		{rval: int64(1), op: "<=", cval: int64(0), expected: false},
		{rval: int64(1), op: "<=", cval: int64(1), expected: true},
		{rval: int64(1), op: "<=", cval: int64(2), expected: true},
		{rval: int64(1), op: ">", cval: int64(0), expected: true},
		{rval: int64(1), op: ">", cval: int64(1), expected: false},
		{rval: int64(1), op: ">", cval: int64(2), expected: false},
		{rval: int64(1), op: ">=", cval: int64(0), expected: true},
		{rval: int64(1), op: ">=", cval: int64(1), expected: true},
		{rval: int64(1), op: ">=", cval: int64(2), expected: false},
		{rval: int64(0), op: "invalid", cval: int64(0), expected: false},

		// uint
		{rval: uint(1), op: "=", cval: uint(0), expected: false},
		{rval: uint(1), op: "=", cval: uint(1), expected: true},
		{rval: uint(1), op: "!=", cval: uint(0), expected: true},
		{rval: uint(1), op: "!=", cval: uint(1), expected: false},
		{rval: uint(1), op: "<", cval: uint(0), expected: false},
		{rval: uint(1), op: "<", cval: uint(1), expected: false},
		{rval: uint(1), op: "<", cval: uint(2), expected: true},
		{rval: uint(1), op: "<=", cval: uint(0), expected: false},
		{rval: uint(1), op: "<=", cval: uint(1), expected: true},
		{rval: uint(1), op: "<=", cval: uint(2), expected: true},
		{rval: uint(1), op: ">", cval: uint(0), expected: true},
		{rval: uint(1), op: ">", cval: uint(1), expected: false},
		{rval: uint(1), op: ">", cval: uint(2), expected: false},
		{rval: uint(1), op: ">=", cval: uint(0), expected: true},
		{rval: uint(1), op: ">=", cval: uint(1), expected: true},
		{rval: uint(1), op: ">=", cval: uint(2), expected: false},
		{rval: uint(0), op: "invalid", cval: uint(0), expected: false},

		// uint8
		{rval: uint8(1), op: "=", cval: uint8(0), expected: false},
		{rval: uint8(1), op: "=", cval: uint8(1), expected: true},
		{rval: uint8(1), op: "!=", cval: uint8(0), expected: true},
		{rval: uint8(1), op: "!=", cval: uint8(1), expected: false},
		{rval: uint8(1), op: "<", cval: uint8(0), expected: false},
		{rval: uint8(1), op: "<", cval: uint8(1), expected: false},
		{rval: uint8(1), op: "<", cval: uint8(2), expected: true},
		{rval: uint8(1), op: "<=", cval: uint8(0), expected: false},
		{rval: uint8(1), op: "<=", cval: uint8(1), expected: true},
		{rval: uint8(1), op: "<=", cval: uint8(2), expected: true},
		{rval: uint8(1), op: ">", cval: uint8(0), expected: true},
		{rval: uint8(1), op: ">", cval: uint8(1), expected: false},
		{rval: uint8(1), op: ">", cval: uint8(2), expected: false},
		{rval: uint8(1), op: ">=", cval: uint8(0), expected: true},
		{rval: uint8(1), op: ">=", cval: uint8(1), expected: true},
		{rval: uint8(1), op: ">=", cval: uint8(2), expected: false},
		{rval: uint8(0), op: "invalid", cval: uint8(0), expected: false},

		// uint16
		{rval: uint16(1), op: "=", cval: uint16(0), expected: false},
		{rval: uint16(1), op: "=", cval: uint16(1), expected: true},
		{rval: uint16(1), op: "!=", cval: uint16(0), expected: true},
		{rval: uint16(1), op: "!=", cval: uint16(1), expected: false},
		{rval: uint16(1), op: "<", cval: uint16(0), expected: false},
		{rval: uint16(1), op: "<", cval: uint16(1), expected: false},
		{rval: uint16(1), op: "<", cval: uint16(2), expected: true},
		{rval: uint16(1), op: "<=", cval: uint16(0), expected: false},
		{rval: uint16(1), op: "<=", cval: uint16(1), expected: true},
		{rval: uint16(1), op: "<=", cval: uint16(2), expected: true},
		{rval: uint16(1), op: ">", cval: uint16(0), expected: true},
		{rval: uint16(1), op: ">", cval: uint16(1), expected: false},
		{rval: uint16(1), op: ">", cval: uint16(2), expected: false},
		{rval: uint16(1), op: ">=", cval: uint16(0), expected: true},
		{rval: uint16(1), op: ">=", cval: uint16(1), expected: true},
		{rval: uint16(1), op: ">=", cval: uint16(2), expected: false},
		{rval: uint16(0), op: "invalid", cval: uint16(0), expected: false},

		// uint32
		{rval: uint32(1), op: "=", cval: uint32(0), expected: false},
		{rval: uint32(1), op: "=", cval: uint32(1), expected: true},
		{rval: uint32(1), op: "!=", cval: uint32(0), expected: true},
		{rval: uint32(1), op: "!=", cval: uint32(1), expected: false},
		{rval: uint32(1), op: "<", cval: uint32(0), expected: false},
		{rval: uint32(1), op: "<", cval: uint32(1), expected: false},
		{rval: uint32(1), op: "<", cval: uint32(2), expected: true},
		{rval: uint32(1), op: "<=", cval: uint32(0), expected: false},
		{rval: uint32(1), op: "<=", cval: uint32(1), expected: true},
		{rval: uint32(1), op: "<=", cval: uint32(2), expected: true},
		{rval: uint32(1), op: ">", cval: uint32(0), expected: true},
		{rval: uint32(1), op: ">", cval: uint32(1), expected: false},
		{rval: uint32(1), op: ">", cval: uint32(2), expected: false},
		{rval: uint32(1), op: ">=", cval: uint32(0), expected: true},
		{rval: uint32(1), op: ">=", cval: uint32(1), expected: true},
		{rval: uint32(1), op: ">=", cval: uint32(2), expected: false},
		{rval: uint32(0), op: "invalid", cval: uint32(0), expected: false},

		// uint64
		{rval: uint64(1), op: "=", cval: uint64(0), expected: false},
		{rval: uint64(1), op: "=", cval: uint64(1), expected: true},
		{rval: uint64(1), op: "!=", cval: uint64(0), expected: true},
		{rval: uint64(1), op: "!=", cval: uint64(1), expected: false},
		{rval: uint64(1), op: "<", cval: uint64(0), expected: false},
		{rval: uint64(1), op: "<", cval: uint64(1), expected: false},
		{rval: uint64(1), op: "<", cval: uint64(2), expected: true},
		{rval: uint64(1), op: "<=", cval: uint64(0), expected: false},
		{rval: uint64(1), op: "<=", cval: uint64(1), expected: true},
		{rval: uint64(1), op: "<=", cval: uint64(2), expected: true},
		{rval: uint64(1), op: ">", cval: uint64(0), expected: true},
		{rval: uint64(1), op: ">", cval: uint64(1), expected: false},
		{rval: uint64(1), op: ">", cval: uint64(2), expected: false},
		{rval: uint64(1), op: ">=", cval: uint64(0), expected: true},
		{rval: uint64(1), op: ">=", cval: uint64(1), expected: true},
		{rval: uint64(1), op: ">=", cval: uint64(2), expected: false},
		{rval: uint64(0), op: "invalid", cval: uint64(0), expected: false},

		// bool
		{rval: true, op: "=", cval: true, expected: true},
		{rval: true, op: "=", cval: false, expected: false},
		{rval: true, op: "!=", cval: true, expected: false},
		{rval: true, op: "!=", cval: false, expected: true},
		{rval: true, op: "invalid", cval: true, expected: false},

		// time.Time
		{rval: now, op: "=", cval: now.Add(-time.Second), expected: false},
		{rval: now, op: "=", cval: now, expected: true},
		{rval: now, op: "!=", cval: now.Add(-time.Second), expected: true},
		{rval: now, op: "!=", cval: now, expected: false},
		{rval: now, op: "<", cval: now.Add(-time.Second), expected: false},
		{rval: now, op: "<", cval: now, expected: false},
		{rval: now, op: "<", cval: now.Add(time.Second), expected: true},
		{rval: now, op: "<=", cval: now.Add(-time.Second), expected: false},
		{rval: now, op: "<=", cval: now, expected: true},
		{rval: now, op: "<=", cval: now.Add(time.Second), expected: true},
		{rval: now, op: ">", cval: now.Add(-time.Second), expected: true},
		{rval: now, op: ">", cval: now, expected: false},
		{rval: now, op: ">", cval: now.Add(time.Second), expected: false},
		{rval: now, op: ">=", cval: now.Add(-time.Second), expected: true},
		{rval: now, op: ">=", cval: now, expected: true},
		{rval: now, op: ">=", cval: now.Add(time.Second), expected: false},
		{rval: now, op: "invalid", cval: now, expected: false},

		// []byte
		{rval: []byte{1}, op: "=", cval: []byte{0}, expected: false},
		{rval: []byte{1}, op: "=", cval: []byte{1}, expected: true},
		{rval: []byte{0, 1}, op: "=", cval: []byte{0, 1}, expected: true},
		{rval: []byte{0, 0}, op: "=", cval: []byte{0, 1}, expected: false},
		{rval: []byte{1}, op: "!=", cval: []byte{0}, expected: true},
		{rval: []byte{1}, op: "!=", cval: []byte{1}, expected: false},
		{rval: []byte{0, 1}, op: "!=", cval: []byte{0, 1}, expected: false},
		{rval: []byte{0, 0}, op: "!=", cval: []byte{0, 1}, expected: true},
		{rval: []byte{1}, op: "<", cval: []byte{0}, expected: false},
		{rval: []byte{1}, op: "<", cval: []byte{1}, expected: false},
		{rval: []byte{1}, op: "<", cval: []byte{2}, expected: true},
		{rval: []byte{0, 1}, op: "<", cval: []byte{0, 0}, expected: false},
		{rval: []byte{0, 1}, op: "<", cval: []byte{0, 1}, expected: false},
		{rval: []byte{0, 1}, op: "<", cval: []byte{0, 2}, expected: true},
		{rval: []byte{1}, op: "<=", cval: []byte{0}, expected: false},
		{rval: []byte{1}, op: "<=", cval: []byte{1}, expected: true},
		{rval: []byte{1}, op: "<=", cval: []byte{2}, expected: true},
		{rval: []byte{0, 1}, op: "<=", cval: []byte{0, 0}, expected: false},
		{rval: []byte{0, 1}, op: "<=", cval: []byte{0, 1}, expected: true},
		{rval: []byte{0, 1}, op: "<=", cval: []byte{0, 2}, expected: true},
		{rval: []byte{1}, op: ">", cval: []byte{0}, expected: true},
		{rval: []byte{1}, op: ">", cval: []byte{1}, expected: false},
		{rval: []byte{1}, op: ">", cval: []byte{2}, expected: false},
		{rval: []byte{0, 1}, op: ">", cval: []byte{0, 0}, expected: true},
		{rval: []byte{0, 1}, op: ">", cval: []byte{0, 1}, expected: false},
		{rval: []byte{0, 1}, op: ">", cval: []byte{0, 2}, expected: false},
		{rval: []byte{1}, op: ">=", cval: []byte{0}, expected: true},
		{rval: []byte{1}, op: ">=", cval: []byte{1}, expected: true},
		{rval: []byte{1}, op: ">=", cval: []byte{2}, expected: false},
		{rval: []byte{0, 1}, op: ">=", cval: []byte{0, 0}, expected: true},
		{rval: []byte{0, 1}, op: ">=", cval: []byte{0, 1}, expected: true},
		{rval: []byte{0, 1}, op: ">=", cval: []byte{0, 2}, expected: false},
		{rval: []byte{0, 0}, op: "invalid", cval: []byte{0, 0}, expected: false},

		// *string
		{rval: nilptr("string"), op: "=", cval: nilptr("string"), expected: true},
		{rval: nilptr("string"), op: "=", cval: ptr("aaa"), expected: false},
		{rval: ptr("bbb"), op: "=", cval: nilptr("string"), expected: false},
		{rval: ptr("bbb"), op: "=", cval: ptr("aaa"), expected: false},
		{rval: ptr("bbb"), op: "=", cval: ptr("bbb"), expected: true},
		{rval: nilptr("string"), op: "!=", cval: nilptr("string"), expected: false},
		{rval: nilptr("string"), op: "!=", cval: ptr("aaa"), expected: true},
		{rval: ptr("bbb"), op: "!=", cval: nilptr("string"), expected: true},
		{rval: ptr("bbb"), op: "!=", cval: ptr("aaa"), expected: true},
		{rval: ptr("bbb"), op: "!=", cval: ptr("bbb"), expected: false},
		{rval: nilptr("string"), op: "<", cval: nilptr("string"), expected: false},
		{rval: nilptr("string"), op: "<", cval: ptr("aaa"), expected: false},
		{rval: ptr("bbb"), op: "<", cval: nilptr("string"), expected: false},
		{rval: ptr("bbb"), op: "<", cval: ptr("aaa"), expected: false},
		{rval: ptr("bbb"), op: "<", cval: ptr("bbb"), expected: false},
		{rval: ptr("bbb"), op: "<", cval: ptr("ccc"), expected: true},
		{rval: nilptr("string"), op: "<=", cval: nilptr("string"), expected: false},
		{rval: nilptr("string"), op: "<=", cval: ptr("aaa"), expected: false},
		{rval: ptr("bbb"), op: "<=", cval: nilptr("string"), expected: false},
		{rval: ptr("bbb"), op: "<=", cval: ptr("aaa"), expected: false},
		{rval: ptr("bbb"), op: "<=", cval: ptr("bbb"), expected: true},
		{rval: ptr("bbb"), op: "<=", cval: ptr("ccc"), expected: true},
		{rval: nilptr("string"), op: ">", cval: nilptr("string"), expected: false},
		{rval: nilptr("string"), op: ">", cval: ptr("aaa"), expected: false},
		{rval: ptr("bbb"), op: ">", cval: nilptr("string"), expected: false},
		{rval: ptr("bbb"), op: ">", cval: ptr("aaa"), expected: true},
		{rval: ptr("bbb"), op: ">", cval: ptr("bbb"), expected: false},
		{rval: ptr("bbb"), op: ">", cval: ptr("ccc"), expected: false},
		{rval: nilptr("string"), op: ">=", cval: nilptr("string"), expected: false},
		{rval: nilptr("string"), op: ">=", cval: ptr("aaa"), expected: false},
		{rval: ptr("bbb"), op: ">=", cval: nilptr("string"), expected: false},
		{rval: ptr("bbb"), op: ">=", cval: ptr("aaa"), expected: true},
		{rval: ptr("bbb"), op: ">=", cval: ptr("bbb"), expected: true},
		{rval: ptr("bbb"), op: ">=", cval: ptr("ccc"), expected: false},

		// *int
		{rval: nilptr("int"), op: "=", cval: nilptr("int"), expected: true},
		{rval: nilptr("int"), op: "=", cval: ptr(-1), expected: false},
		{rval: ptr(0), op: "=", cval: nilptr("int"), expected: false},
		{rval: ptr(0), op: "=", cval: ptr(-1), expected: false},
		{rval: ptr(0), op: "=", cval: ptr(0), expected: true},
		{rval: nilptr("int"), op: "!=", cval: nilptr("int"), expected: false},
		{rval: nilptr("int"), op: "!=", cval: ptr(-1), expected: true},
		{rval: ptr(0), op: "!=", cval: nilptr("int"), expected: true},
		{rval: ptr(0), op: "!=", cval: ptr(-1), expected: true},
		{rval: ptr(0), op: "!=", cval: ptr(0), expected: false},
		{rval: nilptr("int"), op: "<", cval: nilptr("int"), expected: false},
		{rval: nilptr("int"), op: "<", cval: ptr(-1), expected: false},
		{rval: ptr(0), op: "<", cval: nilptr("int"), expected: false},
		{rval: ptr(0), op: "<", cval: ptr(-1), expected: false},
		{rval: ptr(0), op: "<", cval: ptr(0), expected: false},
		{rval: ptr(0), op: "<", cval: ptr(1), expected: true},
		{rval: nilptr("int"), op: "<=", cval: nilptr("int"), expected: false},
		{rval: nilptr("int"), op: "<=", cval: ptr(-1), expected: false},
		{rval: ptr(0), op: "<=", cval: nilptr("int"), expected: false},
		{rval: ptr(0), op: "<=", cval: ptr(-1), expected: false},
		{rval: ptr(0), op: "<=", cval: ptr(0), expected: true},
		{rval: ptr(0), op: "<=", cval: ptr(1), expected: true},
		{rval: nilptr("int"), op: ">", cval: nilptr("int"), expected: false},
		{rval: nilptr("int"), op: ">", cval: ptr(-1), expected: false},
		{rval: ptr(0), op: ">", cval: nilptr("int"), expected: false},
		{rval: ptr(0), op: ">", cval: ptr(-1), expected: true},
		{rval: ptr(0), op: ">", cval: ptr(0), expected: false},
		{rval: ptr(0), op: ">", cval: ptr(1), expected: false},
		{rval: nilptr("int"), op: ">=", cval: nilptr("int"), expected: false},
		{rval: nilptr("int"), op: ">=", cval: ptr(-1), expected: false},
		{rval: ptr(0), op: ">=", cval: nilptr("int"), expected: false},
		{rval: ptr(0), op: ">=", cval: ptr(-1), expected: true},
		{rval: ptr(0), op: ">=", cval: ptr(0), expected: true},
		{rval: ptr(0), op: ">=", cval: ptr(1), expected: false},

		// *int8
		{rval: nilptr("int8"), op: "=", cval: nilptr("int8"), expected: true},
		{rval: nilptr("int8"), op: "=", cval: ptr(int8(-1)), expected: false},
		{rval: ptr(int8(0)), op: "=", cval: nilptr("int8"), expected: false},
		{rval: ptr(int8(0)), op: "=", cval: ptr(int8(-1)), expected: false},
		{rval: ptr(int8(0)), op: "=", cval: ptr(int8(0)), expected: true},
		{rval: nilptr("int8"), op: "!=", cval: nilptr("int8"), expected: false},
		{rval: nilptr("int8"), op: "!=", cval: ptr(int8(-1)), expected: true},
		{rval: ptr(int8(0)), op: "!=", cval: nilptr("int8"), expected: true},
		{rval: ptr(int8(0)), op: "!=", cval: ptr(int8(-1)), expected: true},
		{rval: ptr(int8(0)), op: "!=", cval: ptr(int8(0)), expected: false},
		{rval: nilptr("int8"), op: "<", cval: nilptr("int8"), expected: false},
		{rval: nilptr("int8"), op: "<", cval: ptr(int8(-1)), expected: false},
		{rval: ptr(int8(0)), op: "<", cval: nilptr("int8"), expected: false},
		{rval: ptr(int8(0)), op: "<", cval: ptr(int8(-1)), expected: false},
		{rval: ptr(int8(0)), op: "<", cval: ptr(int8(0)), expected: false},
		{rval: ptr(int8(0)), op: "<", cval: ptr(int8(1)), expected: true},
		{rval: nilptr("int8"), op: "<=", cval: nilptr("int8"), expected: false},
		{rval: nilptr("int8"), op: "<=", cval: ptr(int8(-1)), expected: false},
		{rval: ptr(int8(0)), op: "<=", cval: nilptr("int8"), expected: false},
		{rval: ptr(int8(0)), op: "<=", cval: ptr(int8(-1)), expected: false},
		{rval: ptr(int8(0)), op: "<=", cval: ptr(int8(0)), expected: true},
		{rval: ptr(int8(0)), op: "<=", cval: ptr(int8(1)), expected: true},
		{rval: nilptr("int8"), op: ">", cval: nilptr("int8"), expected: false},
		{rval: nilptr("int8"), op: ">", cval: ptr(int8(-1)), expected: false},
		{rval: ptr(int8(0)), op: ">", cval: nilptr("int8"), expected: false},
		{rval: ptr(int8(0)), op: ">", cval: ptr(int8(-1)), expected: true},
		{rval: ptr(int8(0)), op: ">", cval: ptr(int8(0)), expected: false},
		{rval: ptr(int8(0)), op: ">", cval: ptr(int8(1)), expected: false},
		{rval: nilptr("int8"), op: ">=", cval: nilptr("int8"), expected: false},
		{rval: nilptr("int8"), op: ">=", cval: ptr(int8(-1)), expected: false},
		{rval: ptr(int8(0)), op: ">=", cval: nilptr("int8"), expected: false},
		{rval: ptr(int8(0)), op: ">=", cval: ptr(int8(-1)), expected: true},
		{rval: ptr(int8(0)), op: ">=", cval: ptr(int8(0)), expected: true},
		{rval: ptr(int8(0)), op: ">=", cval: ptr(int8(1)), expected: false},

		// *int16
		{rval: nilptr("int16"), op: "=", cval: nilptr("int16"), expected: true},
		{rval: nilptr("int16"), op: "=", cval: ptr(int16(-1)), expected: false},
		{rval: ptr(int16(0)), op: "=", cval: nilptr("int16"), expected: false},
		{rval: ptr(int16(0)), op: "=", cval: ptr(int16(-1)), expected: false},
		{rval: ptr(int16(0)), op: "=", cval: ptr(int16(0)), expected: true},
		{rval: nilptr("int16"), op: "!=", cval: nilptr("int16"), expected: false},
		{rval: nilptr("int16"), op: "!=", cval: ptr(int16(-1)), expected: true},
		{rval: ptr(int16(0)), op: "!=", cval: nilptr("int16"), expected: true},
		{rval: ptr(int16(0)), op: "!=", cval: ptr(int16(-1)), expected: true},
		{rval: ptr(int16(0)), op: "!=", cval: ptr(int16(0)), expected: false},
		{rval: nilptr("int16"), op: "<", cval: nilptr("int16"), expected: false},
		{rval: nilptr("int16"), op: "<", cval: ptr(int16(-1)), expected: false},
		{rval: ptr(int16(0)), op: "<", cval: nilptr("int16"), expected: false},
		{rval: ptr(int16(0)), op: "<", cval: ptr(int16(-1)), expected: false},
		{rval: ptr(int16(0)), op: "<", cval: ptr(int16(0)), expected: false},
		{rval: ptr(int16(0)), op: "<", cval: ptr(int16(1)), expected: true},
		{rval: nilptr("int16"), op: "<=", cval: nilptr("int16"), expected: false},
		{rval: nilptr("int16"), op: "<=", cval: ptr(int16(-1)), expected: false},
		{rval: ptr(int16(0)), op: "<=", cval: nilptr("int16"), expected: false},
		{rval: ptr(int16(0)), op: "<=", cval: ptr(int16(-1)), expected: false},
		{rval: ptr(int16(0)), op: "<=", cval: ptr(int16(0)), expected: true},
		{rval: ptr(int16(0)), op: "<=", cval: ptr(int16(1)), expected: true},
		{rval: nilptr("int16"), op: ">", cval: nilptr("int16"), expected: false},
		{rval: nilptr("int16"), op: ">", cval: ptr(int16(-1)), expected: false},
		{rval: ptr(int16(0)), op: ">", cval: nilptr("int16"), expected: false},
		{rval: ptr(int16(0)), op: ">", cval: ptr(int16(-1)), expected: true},
		{rval: ptr(int16(0)), op: ">", cval: ptr(int16(0)), expected: false},
		{rval: ptr(int16(0)), op: ">", cval: ptr(int16(1)), expected: false},
		{rval: nilptr("int16"), op: ">=", cval: nilptr("int16"), expected: false},
		{rval: nilptr("int16"), op: ">=", cval: ptr(int16(-1)), expected: false},
		{rval: ptr(int16(0)), op: ">=", cval: nilptr("int16"), expected: false},
		{rval: ptr(int16(0)), op: ">=", cval: ptr(int16(-1)), expected: true},
		{rval: ptr(int16(0)), op: ">=", cval: ptr(int16(0)), expected: true},
		{rval: ptr(int16(0)), op: ">=", cval: ptr(int16(1)), expected: false},

		// *int32
		{rval: nilptr("int32"), op: "=", cval: nilptr("int32"), expected: true},
		{rval: nilptr("int32"), op: "=", cval: ptr(int32(-1)), expected: false},
		{rval: ptr(int32(0)), op: "=", cval: nilptr("int32"), expected: false},
		{rval: ptr(int32(0)), op: "=", cval: ptr(int32(-1)), expected: false},
		{rval: ptr(int32(0)), op: "=", cval: ptr(int32(0)), expected: true},
		{rval: nilptr("int32"), op: "!=", cval: nilptr("int32"), expected: false},
		{rval: nilptr("int32"), op: "!=", cval: ptr(int32(-1)), expected: true},
		{rval: ptr(int32(0)), op: "!=", cval: nilptr("int32"), expected: true},
		{rval: ptr(int32(0)), op: "!=", cval: ptr(int32(-1)), expected: true},
		{rval: ptr(int32(0)), op: "!=", cval: ptr(int32(0)), expected: false},
		{rval: nilptr("int32"), op: "<", cval: nilptr("int32"), expected: false},
		{rval: nilptr("int32"), op: "<", cval: ptr(int32(-1)), expected: false},
		{rval: ptr(int32(0)), op: "<", cval: nilptr("int32"), expected: false},
		{rval: ptr(int32(0)), op: "<", cval: ptr(int32(-1)), expected: false},
		{rval: ptr(int32(0)), op: "<", cval: ptr(int32(0)), expected: false},
		{rval: ptr(int32(0)), op: "<", cval: ptr(int32(1)), expected: true},
		{rval: nilptr("int32"), op: "<=", cval: nilptr("int32"), expected: false},
		{rval: nilptr("int32"), op: "<=", cval: ptr(int32(-1)), expected: false},
		{rval: ptr(int32(0)), op: "<=", cval: nilptr("int32"), expected: false},
		{rval: ptr(int32(0)), op: "<=", cval: ptr(int32(-1)), expected: false},
		{rval: ptr(int32(0)), op: "<=", cval: ptr(int32(0)), expected: true},
		{rval: ptr(int32(0)), op: "<=", cval: ptr(int32(1)), expected: true},
		{rval: nilptr("int32"), op: ">", cval: nilptr("int32"), expected: false},
		{rval: nilptr("int32"), op: ">", cval: ptr(int32(-1)), expected: false},
		{rval: ptr(int32(0)), op: ">", cval: nilptr("int32"), expected: false},
		{rval: ptr(int32(0)), op: ">", cval: ptr(int32(-1)), expected: true},
		{rval: ptr(int32(0)), op: ">", cval: ptr(int32(0)), expected: false},
		{rval: ptr(int32(0)), op: ">", cval: ptr(int32(1)), expected: false},
		{rval: nilptr("int32"), op: ">=", cval: nilptr("int32"), expected: false},
		{rval: nilptr("int32"), op: ">=", cval: ptr(int32(-1)), expected: false},
		{rval: ptr(int32(0)), op: ">=", cval: nilptr("int32"), expected: false},
		{rval: ptr(int32(0)), op: ">=", cval: ptr(int32(-1)), expected: true},
		{rval: ptr(int32(0)), op: ">=", cval: ptr(int32(0)), expected: true},
		{rval: ptr(int32(0)), op: ">=", cval: ptr(int32(1)), expected: false},

		// *int64
		{rval: nilptr("int64"), op: "=", cval: nilptr("int64"), expected: true},
		{rval: nilptr("int64"), op: "=", cval: ptr(int64(-1)), expected: false},
		{rval: ptr(int64(0)), op: "=", cval: nilptr("int64"), expected: false},
		{rval: ptr(int64(0)), op: "=", cval: ptr(int64(-1)), expected: false},
		{rval: ptr(int64(0)), op: "=", cval: ptr(int64(0)), expected: true},
		{rval: nilptr("int64"), op: "!=", cval: nilptr("int64"), expected: false},
		{rval: nilptr("int64"), op: "!=", cval: ptr(int64(-1)), expected: true},
		{rval: ptr(int64(0)), op: "!=", cval: nilptr("int64"), expected: true},
		{rval: ptr(int64(0)), op: "!=", cval: ptr(int64(-1)), expected: true},
		{rval: ptr(int64(0)), op: "!=", cval: ptr(int64(0)), expected: false},
		{rval: nilptr("int64"), op: "<", cval: nilptr("int64"), expected: false},
		{rval: nilptr("int64"), op: "<", cval: ptr(int64(-1)), expected: false},
		{rval: ptr(int64(0)), op: "<", cval: nilptr("int64"), expected: false},
		{rval: ptr(int64(0)), op: "<", cval: ptr(int64(-1)), expected: false},
		{rval: ptr(int64(0)), op: "<", cval: ptr(int64(0)), expected: false},
		{rval: ptr(int64(0)), op: "<", cval: ptr(int64(1)), expected: true},
		{rval: nilptr("int64"), op: "<=", cval: nilptr("int64"), expected: false},
		{rval: nilptr("int64"), op: "<=", cval: ptr(int64(-1)), expected: false},
		{rval: ptr(int64(0)), op: "<=", cval: nilptr("int64"), expected: false},
		{rval: ptr(int64(0)), op: "<=", cval: ptr(int64(-1)), expected: false},
		{rval: ptr(int64(0)), op: "<=", cval: ptr(int64(0)), expected: true},
		{rval: ptr(int64(0)), op: "<=", cval: ptr(int64(1)), expected: true},
		{rval: nilptr("int64"), op: ">", cval: nilptr("int64"), expected: false},
		{rval: nilptr("int64"), op: ">", cval: ptr(int64(-1)), expected: false},
		{rval: ptr(int64(0)), op: ">", cval: nilptr("int64"), expected: false},
		{rval: ptr(int64(0)), op: ">", cval: ptr(int64(-1)), expected: true},
		{rval: ptr(int64(0)), op: ">", cval: ptr(int64(0)), expected: false},
		{rval: ptr(int64(0)), op: ">", cval: ptr(int64(1)), expected: false},
		{rval: nilptr("int64"), op: ">=", cval: nilptr("int64"), expected: false},
		{rval: nilptr("int64"), op: ">=", cval: ptr(int64(-1)), expected: false},
		{rval: ptr(int64(0)), op: ">=", cval: nilptr("int64"), expected: false},
		{rval: ptr(int64(0)), op: ">=", cval: ptr(int64(-1)), expected: true},
		{rval: ptr(int64(0)), op: ">=", cval: ptr(int64(0)), expected: true},
		{rval: ptr(int64(0)), op: ">=", cval: ptr(int64(1)), expected: false},

		// *uint
		{rval: nilptr("uint"), op: "=", cval: nilptr("uint"), expected: true},
		{rval: nilptr("uint"), op: "=", cval: ptr(uint(0)), expected: false},
		{rval: ptr(uint(1)), op: "=", cval: nilptr("uint"), expected: false},
		{rval: ptr(uint(1)), op: "=", cval: ptr(uint(0)), expected: false},
		{rval: ptr(uint(1)), op: "=", cval: ptr(uint(1)), expected: true},
		{rval: nilptr("uint"), op: "!=", cval: nilptr("uint"), expected: false},
		{rval: nilptr("uint"), op: "!=", cval: ptr(uint(0)), expected: true},
		{rval: ptr(uint(1)), op: "!=", cval: nilptr("uint"), expected: true},
		{rval: ptr(uint(1)), op: "!=", cval: ptr(uint(0)), expected: true},
		{rval: ptr(uint(1)), op: "!=", cval: ptr(uint(1)), expected: false},
		{rval: nilptr("uint"), op: "<", cval: nilptr("uint"), expected: false},
		{rval: nilptr("uint"), op: "<", cval: ptr(uint(0)), expected: false},
		{rval: ptr(uint(1)), op: "<", cval: nilptr("uint"), expected: false},
		{rval: ptr(uint(1)), op: "<", cval: ptr(uint(0)), expected: false},
		{rval: ptr(uint(1)), op: "<", cval: ptr(uint(1)), expected: false},
		{rval: ptr(uint(1)), op: "<", cval: ptr(uint(2)), expected: true},
		{rval: nilptr("uint"), op: "<=", cval: nilptr("uint"), expected: false},
		{rval: nilptr("uint"), op: "<=", cval: ptr(uint(0)), expected: false},
		{rval: ptr(uint(1)), op: "<=", cval: nilptr("uint"), expected: false},
		{rval: ptr(uint(1)), op: "<=", cval: ptr(uint(0)), expected: false},
		{rval: ptr(uint(1)), op: "<=", cval: ptr(uint(1)), expected: true},
		{rval: ptr(uint(1)), op: "<=", cval: ptr(uint(2)), expected: true},
		{rval: nilptr("uint"), op: ">", cval: nilptr("uint"), expected: false},
		{rval: nilptr("uint"), op: ">", cval: ptr(uint(0)), expected: false},
		{rval: ptr(uint(1)), op: ">", cval: nilptr("uint"), expected: false},
		{rval: ptr(uint(1)), op: ">", cval: ptr(uint(0)), expected: true},
		{rval: ptr(uint(1)), op: ">", cval: ptr(uint(1)), expected: false},
		{rval: ptr(uint(1)), op: ">", cval: ptr(uint(2)), expected: false},
		{rval: nilptr("uint"), op: ">=", cval: nilptr("uint"), expected: false},
		{rval: nilptr("uint"), op: ">=", cval: ptr(uint(0)), expected: false},
		{rval: ptr(uint(1)), op: ">=", cval: nilptr("uint"), expected: false},
		{rval: ptr(uint(1)), op: ">=", cval: ptr(uint(0)), expected: true},
		{rval: ptr(uint(1)), op: ">=", cval: ptr(uint(1)), expected: true},
		{rval: ptr(uint(1)), op: ">=", cval: ptr(uint(2)), expected: false},

		// *uint8
		{rval: nilptr("uint8"), op: "=", cval: nilptr("uint8"), expected: true},
		{rval: nilptr("uint8"), op: "=", cval: ptr(uint8(0)), expected: false},
		{rval: ptr(uint8(1)), op: "=", cval: nilptr("uint8"), expected: false},
		{rval: ptr(uint8(1)), op: "=", cval: ptr(uint8(0)), expected: false},
		{rval: ptr(uint8(1)), op: "=", cval: ptr(uint8(1)), expected: true},
		{rval: nilptr("uint8"), op: "!=", cval: nilptr("uint8"), expected: false},
		{rval: nilptr("uint8"), op: "!=", cval: ptr(uint8(0)), expected: true},
		{rval: ptr(uint8(1)), op: "!=", cval: nilptr("uint8"), expected: true},
		{rval: ptr(uint8(1)), op: "!=", cval: ptr(uint8(0)), expected: true},
		{rval: ptr(uint8(1)), op: "!=", cval: ptr(uint8(1)), expected: false},
		{rval: nilptr("uint8"), op: "<", cval: nilptr("uint8"), expected: false},
		{rval: nilptr("uint8"), op: "<", cval: ptr(uint8(0)), expected: false},
		{rval: ptr(uint8(1)), op: "<", cval: nilptr("uint8"), expected: false},
		{rval: ptr(uint8(1)), op: "<", cval: ptr(uint8(0)), expected: false},
		{rval: ptr(uint8(1)), op: "<", cval: ptr(uint8(1)), expected: false},
		{rval: ptr(uint8(1)), op: "<", cval: ptr(uint8(2)), expected: true},
		{rval: nilptr("uint8"), op: "<=", cval: nilptr("uint8"), expected: false},
		{rval: nilptr("uint8"), op: "<=", cval: ptr(uint8(0)), expected: false},
		{rval: ptr(uint8(1)), op: "<=", cval: nilptr("uint8"), expected: false},
		{rval: ptr(uint8(1)), op: "<=", cval: ptr(uint8(0)), expected: false},
		{rval: ptr(uint8(1)), op: "<=", cval: ptr(uint8(1)), expected: true},
		{rval: ptr(uint8(1)), op: "<=", cval: ptr(uint8(2)), expected: true},
		{rval: nilptr("uint8"), op: ">", cval: nilptr("uint8"), expected: false},
		{rval: nilptr("uint8"), op: ">", cval: ptr(uint8(0)), expected: false},
		{rval: ptr(uint8(1)), op: ">", cval: nilptr("uint8"), expected: false},
		{rval: ptr(uint8(1)), op: ">", cval: ptr(uint8(0)), expected: true},
		{rval: ptr(uint8(1)), op: ">", cval: ptr(uint8(1)), expected: false},
		{rval: ptr(uint8(1)), op: ">", cval: ptr(uint8(2)), expected: false},
		{rval: nilptr("uint8"), op: ">=", cval: nilptr("uint8"), expected: false},
		{rval: nilptr("uint8"), op: ">=", cval: ptr(uint8(0)), expected: false},
		{rval: ptr(uint8(1)), op: ">=", cval: nilptr("uint8"), expected: false},
		{rval: ptr(uint8(1)), op: ">=", cval: ptr(uint8(0)), expected: true},
		{rval: ptr(uint8(1)), op: ">=", cval: ptr(uint8(1)), expected: true},
		{rval: ptr(uint8(1)), op: ">=", cval: ptr(uint8(2)), expected: false},

		// *uint16
		{rval: nilptr("uint16"), op: "=", cval: nilptr("uint16"), expected: true},
		{rval: nilptr("uint16"), op: "=", cval: ptr(uint16(0)), expected: false},
		{rval: ptr(uint16(1)), op: "=", cval: nilptr("uint16"), expected: false},
		{rval: ptr(uint16(1)), op: "=", cval: ptr(uint16(0)), expected: false},
		{rval: ptr(uint16(1)), op: "=", cval: ptr(uint16(1)), expected: true},
		{rval: nilptr("uint16"), op: "!=", cval: nilptr("uint16"), expected: false},
		{rval: nilptr("uint16"), op: "!=", cval: ptr(uint16(0)), expected: true},
		{rval: ptr(uint16(1)), op: "!=", cval: nilptr("uint16"), expected: true},
		{rval: ptr(uint16(1)), op: "!=", cval: ptr(uint16(0)), expected: true},
		{rval: ptr(uint16(1)), op: "!=", cval: ptr(uint16(1)), expected: false},
		{rval: nilptr("uint16"), op: "<", cval: nilptr("uint16"), expected: false},
		{rval: nilptr("uint16"), op: "<", cval: ptr(uint16(0)), expected: false},
		{rval: ptr(uint16(1)), op: "<", cval: nilptr("uint16"), expected: false},
		{rval: ptr(uint16(1)), op: "<", cval: ptr(uint16(0)), expected: false},
		{rval: ptr(uint16(1)), op: "<", cval: ptr(uint16(1)), expected: false},
		{rval: ptr(uint16(1)), op: "<", cval: ptr(uint16(2)), expected: true},
		{rval: nilptr("uint16"), op: "<=", cval: nilptr("uint16"), expected: false},
		{rval: nilptr("uint16"), op: "<=", cval: ptr(uint16(0)), expected: false},
		{rval: ptr(uint16(1)), op: "<=", cval: nilptr("uint16"), expected: false},
		{rval: ptr(uint16(1)), op: "<=", cval: ptr(uint16(0)), expected: false},
		{rval: ptr(uint16(1)), op: "<=", cval: ptr(uint16(1)), expected: true},
		{rval: ptr(uint16(1)), op: "<=", cval: ptr(uint16(2)), expected: true},
		{rval: nilptr("uint16"), op: ">", cval: nilptr("uint16"), expected: false},
		{rval: nilptr("uint16"), op: ">", cval: ptr(uint16(0)), expected: false},
		{rval: ptr(uint16(1)), op: ">", cval: nilptr("uint16"), expected: false},
		{rval: ptr(uint16(1)), op: ">", cval: ptr(uint16(0)), expected: true},
		{rval: ptr(uint16(1)), op: ">", cval: ptr(uint16(1)), expected: false},
		{rval: ptr(uint16(1)), op: ">", cval: ptr(uint16(2)), expected: false},
		{rval: nilptr("uint16"), op: ">=", cval: nilptr("uint16"), expected: false},
		{rval: nilptr("uint16"), op: ">=", cval: ptr(uint16(0)), expected: false},
		{rval: ptr(uint16(1)), op: ">=", cval: nilptr("uint16"), expected: false},
		{rval: ptr(uint16(1)), op: ">=", cval: ptr(uint16(0)), expected: true},
		{rval: ptr(uint16(1)), op: ">=", cval: ptr(uint16(1)), expected: true},
		{rval: ptr(uint16(1)), op: ">=", cval: ptr(uint16(2)), expected: false},

		// *uint32
		{rval: nilptr("uint32"), op: "=", cval: nilptr("uint32"), expected: true},
		{rval: nilptr("uint32"), op: "=", cval: ptr(uint32(0)), expected: false},
		{rval: ptr(uint32(1)), op: "=", cval: nilptr("uint32"), expected: false},
		{rval: ptr(uint32(1)), op: "=", cval: ptr(uint32(0)), expected: false},
		{rval: ptr(uint32(1)), op: "=", cval: ptr(uint32(1)), expected: true},
		{rval: nilptr("uint32"), op: "!=", cval: nilptr("uint32"), expected: false},
		{rval: nilptr("uint32"), op: "!=", cval: ptr(uint32(0)), expected: true},
		{rval: ptr(uint32(1)), op: "!=", cval: nilptr("uint32"), expected: true},
		{rval: ptr(uint32(1)), op: "!=", cval: ptr(uint32(0)), expected: true},
		{rval: ptr(uint32(1)), op: "!=", cval: ptr(uint32(1)), expected: false},
		{rval: nilptr("uint32"), op: "<", cval: nilptr("uint32"), expected: false},
		{rval: nilptr("uint32"), op: "<", cval: ptr(uint32(0)), expected: false},
		{rval: ptr(uint32(1)), op: "<", cval: nilptr("uint32"), expected: false},
		{rval: ptr(uint32(1)), op: "<", cval: ptr(uint32(0)), expected: false},
		{rval: ptr(uint32(1)), op: "<", cval: ptr(uint32(1)), expected: false},
		{rval: ptr(uint32(1)), op: "<", cval: ptr(uint32(2)), expected: true},
		{rval: nilptr("uint32"), op: "<=", cval: nilptr("uint32"), expected: false},
		{rval: nilptr("uint32"), op: "<=", cval: ptr(uint32(0)), expected: false},
		{rval: ptr(uint32(1)), op: "<=", cval: nilptr("uint32"), expected: false},
		{rval: ptr(uint32(1)), op: "<=", cval: ptr(uint32(0)), expected: false},
		{rval: ptr(uint32(1)), op: "<=", cval: ptr(uint32(1)), expected: true},
		{rval: ptr(uint32(1)), op: "<=", cval: ptr(uint32(2)), expected: true},
		{rval: nilptr("uint32"), op: ">", cval: nilptr("uint32"), expected: false},
		{rval: nilptr("uint32"), op: ">", cval: ptr(uint32(0)), expected: false},
		{rval: ptr(uint32(1)), op: ">", cval: nilptr("uint32"), expected: false},
		{rval: ptr(uint32(1)), op: ">", cval: ptr(uint32(0)), expected: true},
		{rval: ptr(uint32(1)), op: ">", cval: ptr(uint32(1)), expected: false},
		{rval: ptr(uint32(1)), op: ">", cval: ptr(uint32(2)), expected: false},
		{rval: nilptr("uint32"), op: ">=", cval: nilptr("uint32"), expected: false},
		{rval: nilptr("uint32"), op: ">=", cval: ptr(uint32(0)), expected: false},
		{rval: ptr(uint32(1)), op: ">=", cval: nilptr("uint32"), expected: false},
		{rval: ptr(uint32(1)), op: ">=", cval: ptr(uint32(0)), expected: true},
		{rval: ptr(uint32(1)), op: ">=", cval: ptr(uint32(1)), expected: true},
		{rval: ptr(uint32(1)), op: ">=", cval: ptr(uint32(2)), expected: false},

		// *uint64
		{rval: nilptr("uint64"), op: "=", cval: nilptr("uint64"), expected: true},
		{rval: nilptr("uint64"), op: "=", cval: ptr(uint64(0)), expected: false},
		{rval: ptr(uint64(1)), op: "=", cval: nilptr("uint64"), expected: false},
		{rval: ptr(uint64(1)), op: "=", cval: ptr(uint64(0)), expected: false},
		{rval: ptr(uint64(1)), op: "=", cval: ptr(uint64(1)), expected: true},
		{rval: nilptr("uint64"), op: "!=", cval: nilptr("uint64"), expected: false},
		{rval: nilptr("uint64"), op: "!=", cval: ptr(uint64(0)), expected: true},
		{rval: ptr(uint64(1)), op: "!=", cval: nilptr("uint64"), expected: true},
		{rval: ptr(uint64(1)), op: "!=", cval: ptr(uint64(0)), expected: true},
		{rval: ptr(uint64(1)), op: "!=", cval: ptr(uint64(1)), expected: false},
		{rval: nilptr("uint64"), op: "<", cval: nilptr("uint64"), expected: false},
		{rval: nilptr("uint64"), op: "<", cval: ptr(uint64(0)), expected: false},
		{rval: ptr(uint64(1)), op: "<", cval: nilptr("uint64"), expected: false},
		{rval: ptr(uint64(1)), op: "<", cval: ptr(uint64(0)), expected: false},
		{rval: ptr(uint64(1)), op: "<", cval: ptr(uint64(1)), expected: false},
		{rval: ptr(uint64(1)), op: "<", cval: ptr(uint64(2)), expected: true},
		{rval: nilptr("uint64"), op: "<=", cval: nilptr("uint64"), expected: false},
		{rval: nilptr("uint64"), op: "<=", cval: ptr(uint64(0)), expected: false},
		{rval: ptr(uint64(1)), op: "<=", cval: nilptr("uint64"), expected: false},
		{rval: ptr(uint64(1)), op: "<=", cval: ptr(uint64(0)), expected: false},
		{rval: ptr(uint64(1)), op: "<=", cval: ptr(uint64(1)), expected: true},
		{rval: ptr(uint64(1)), op: "<=", cval: ptr(uint64(2)), expected: true},
		{rval: nilptr("uint64"), op: ">", cval: nilptr("uint64"), expected: false},
		{rval: nilptr("uint64"), op: ">", cval: ptr(uint64(0)), expected: false},
		{rval: ptr(uint64(1)), op: ">", cval: nilptr("uint64"), expected: false},
		{rval: ptr(uint64(1)), op: ">", cval: ptr(uint64(0)), expected: true},
		{rval: ptr(uint64(1)), op: ">", cval: ptr(uint64(1)), expected: false},
		{rval: ptr(uint64(1)), op: ">", cval: ptr(uint64(2)), expected: false},
		{rval: nilptr("uint64"), op: ">=", cval: nilptr("uint64"), expected: false},
		{rval: nilptr("uint64"), op: ">=", cval: ptr(uint64(0)), expected: false},
		{rval: ptr(uint64(1)), op: ">=", cval: nilptr("uint64"), expected: false},
		{rval: ptr(uint64(1)), op: ">=", cval: ptr(uint64(0)), expected: true},
		{rval: ptr(uint64(1)), op: ">=", cval: ptr(uint64(1)), expected: true},
		{rval: ptr(uint64(1)), op: ">=", cval: ptr(uint64(2)), expected: false},

		// *bool
		{rval: nilptr("bool"), op: "=", cval: nilptr("bool"), expected: true},
		{rval: nilptr("bool"), op: "=", cval: ptr(false), expected: false},
		{rval: ptr(true), op: "=", cval: nilptr("bool"), expected: false},
		{rval: ptr(true), op: "=", cval: ptr(true), expected: true},
		{rval: ptr(true), op: "=", cval: ptr(false), expected: false},
		{rval: nilptr("bool"), op: "!=", cval: nilptr("bool"), expected: false},
		{rval: nilptr("bool"), op: "!=", cval: ptr(false), expected: true},
		{rval: ptr(true), op: "!=", cval: nilptr("bool"), expected: true},
		{rval: ptr(true), op: "!=", cval: ptr(true), expected: false},
		{rval: ptr(true), op: "!=", cval: ptr(false), expected: true},
		{rval: nilptr("bool"), op: "bad", cval: nilptr("bool"), expected: false},

		// *time.Time
		{rval: nilptr("time.Time"), op: "=", cval: nilptr("time.Time"), expected: true},
		{rval: nilptr("time.Time"), op: "=", cval: ptr(now), expected: false},
		{rval: ptr(now), op: "=", cval: nilptr("time.Time"), expected: false},
		{rval: ptr(now), op: "=", cval: ptr(now.Add(-time.Second)), expected: false},
		{rval: ptr(now), op: "=", cval: ptr(now), expected: true},
		{rval: nilptr("time.Time"), op: "!=", cval: nilptr("time.Time"), expected: false},
		{rval: nilptr("time.Time"), op: "!=", cval: ptr(now), expected: true},
		{rval: ptr(now), op: "!=", cval: nilptr("time.Time"), expected: true},
		{rval: ptr(now), op: "!=", cval: ptr(now.Add(-time.Second)), expected: true},
		{rval: ptr(now), op: "!=", cval: ptr(now), expected: false},
		{rval: nilptr("time.Time"), op: "<", cval: nilptr("time.Time"), expected: false},
		{rval: nilptr("time.Time"), op: "<", cval: ptr(now), expected: false},
		{rval: ptr(now), op: "<", cval: nilptr("time.Time"), expected: false},
		{rval: ptr(now), op: "<", cval: ptr(now.Add(-time.Second)), expected: false},
		{rval: ptr(now), op: "<", cval: ptr(now), expected: false},
		{rval: ptr(now), op: "<", cval: ptr(now.Add(time.Second)), expected: true},
		{rval: nilptr("time.Time"), op: "<=", cval: nilptr("time.Time"), expected: false},
		{rval: nilptr("time.Time"), op: "<=", cval: ptr(now), expected: false},
		{rval: ptr(now), op: "<=", cval: nilptr("time.Time"), expected: false},
		{rval: ptr(now), op: "<=", cval: ptr(now.Add(-time.Second)), expected: false},
		{rval: ptr(now), op: "<=", cval: ptr(now), expected: true},
		{rval: ptr(now), op: "<=", cval: ptr(now.Add(time.Second)), expected: true},
		{rval: nilptr("time.Time"), op: ">", cval: nilptr("time.Time"), expected: false},
		{rval: nilptr("time.Time"), op: ">", cval: ptr(now), expected: false},
		{rval: ptr(now), op: ">", cval: nilptr("time.Time"), expected: false},
		{rval: ptr(now), op: ">", cval: ptr(now.Add(-time.Second)), expected: true},
		{rval: ptr(now), op: ">", cval: ptr(now), expected: false},
		{rval: ptr(now), op: ">", cval: ptr(now.Add(time.Second)), expected: false},
		{rval: nilptr("time.Time"), op: ">=", cval: nilptr("time.Time"), expected: false},
		{rval: nilptr("time.Time"), op: ">=", cval: ptr(now), expected: false},
		{rval: ptr(now), op: ">=", cval: nilptr("time.Time"), expected: false},
		{rval: ptr(now), op: ">=", cval: ptr(now.Add(-time.Second)), expected: true},
		{rval: ptr(now), op: ">=", cval: ptr(now), expected: true},
		{rval: ptr(now), op: ">=", cval: ptr(now.Add(time.Second)), expected: false},

		// *[]byte
		{rval: nilptr("[]byte"), op: "=", cval: nilptr("[]byte"), expected: true},
		{rval: nilptr("[]byte"), op: "=", cval: ptr([]byte{0}), expected: false},
		{rval: ptr([]byte{1}), op: "=", cval: nilptr("[]byte"), expected: false},
		{rval: ptr([]byte{1}), op: "=", cval: ptr([]byte{0}), expected: false},
		{rval: ptr([]byte{1}), op: "=", cval: ptr([]byte{1}), expected: true},
		{rval: nilptr("[]byte"), op: "!=", cval: nilptr("[]byte"), expected: false},
		{rval: nilptr("[]byte"), op: "!=", cval: ptr([]byte{0}), expected: true},
		{rval: ptr([]byte{1}), op: "!=", cval: nilptr("[]byte"), expected: true},
		{rval: ptr([]byte{1}), op: "!=", cval: ptr([]byte{0}), expected: true},
		{rval: ptr([]byte{1}), op: "!=", cval: ptr([]byte{1}), expected: false},
		{rval: nilptr("[]byte"), op: "<", cval: nilptr("[]byte"), expected: false},
		{rval: nilptr("[]byte"), op: "<", cval: ptr([]byte{0}), expected: false},
		{rval: ptr([]byte{1}), op: "<", cval: nilptr("[]byte"), expected: false},
		{rval: ptr([]byte{1}), op: "<", cval: ptr([]byte{0}), expected: false},
		{rval: ptr([]byte{1}), op: "<", cval: ptr([]byte{1}), expected: false},
		{rval: ptr([]byte{1}), op: "<", cval: ptr([]byte{2}), expected: true},
		{rval: nilptr("[]byte"), op: "<=", cval: nilptr("[]byte"), expected: false},
		{rval: nilptr("[]byte"), op: "<=", cval: ptr([]byte{0}), expected: false},
		{rval: ptr([]byte{1}), op: "<=", cval: nilptr("[]byte"), expected: false},
		{rval: ptr([]byte{1}), op: "<=", cval: ptr([]byte{0}), expected: false},
		{rval: ptr([]byte{1}), op: "<=", cval: ptr([]byte{1}), expected: true},
		{rval: ptr([]byte{1}), op: "<=", cval: ptr([]byte{2}), expected: true},
		{rval: nilptr("[]byte"), op: ">", cval: nilptr("[]byte"), expected: false},
		{rval: nilptr("[]byte"), op: ">", cval: ptr([]byte{0}), expected: false},
		{rval: ptr([]byte{1}), op: ">", cval: nilptr("[]byte"), expected: false},
		{rval: ptr([]byte{1}), op: ">", cval: ptr([]byte{0}), expected: true},
		{rval: ptr([]byte{1}), op: ">", cval: ptr([]byte{1}), expected: false},
		{rval: ptr([]byte{1}), op: ">", cval: ptr([]byte{2}), expected: false},
		{rval: nilptr("[]byte"), op: ">=", cval: nilptr("[]byte"), expected: false},
		{rval: nilptr("[]byte"), op: ">=", cval: ptr([]byte{0}), expected: false},
		{rval: ptr([]byte{1}), op: ">=", cval: nilptr("[]byte"), expected: false},
		{rval: ptr([]byte{1}), op: ">=", cval: ptr([]byte{0}), expected: true},
		{rval: ptr([]byte{1}), op: ">=", cval: ptr([]byte{1}), expected: true},
		{rval: ptr([]byte{1}), op: ">=", cval: ptr([]byte{2}), expected: false},

		// Invalid type
		{rval: func() {}, op: "=", cval: func() {}, expected: false},
	}

	for _, test := range attrTests {
		typ := &Type{Name: "type"}
		ty, n := GetAttrType(fmt.Sprintf("%T", test.rval))
		typ.Attrs = map[string]Attr{
			"attr": {
				Name:     "attr",
				Type:     ty,
				Nullable: n,
			},
		}

		res := &SoftResource{}
		res.SetType(typ)
		res.Set("attr", test.rval)

		filter := &Filter{
			Field: "attr",
			Op:    test.op,
			Val:   test.cval,
		}

		assert.Equal(
			test.expected,
			filter.IsAllowed(res),
			fmt.Sprintf("%v %s %v should be %v", test.rval, test.op, test.cval, test.expected),
		)
	}

	// Tests for relationships
	relTests := []struct {
		rval     any
		op       string
		cval     any
		expected bool
	}{
		// to-one
		{rval: "id1", op: "=", cval: "id1", expected: true},
		{rval: "id1", op: "=", cval: "id2", expected: false},
		{rval: "id1", op: "!=", cval: "id1", expected: false},
		{rval: "id1", op: "!=", cval: "id2", expected: true},
		{rval: "id1", op: "in", cval: []string{"id1"}, expected: true},
		{rval: "id1", op: "in", cval: []string{"id2"}, expected: false},
		{rval: "id1", op: "in", cval: []string{"id1", "id2"}, expected: true},
		{rval: "id1", op: "in", cval: []string{"id2", "id3"}, expected: false},

		// to-many
		{rval: []string{"id1"}, op: "=", cval: []string{"id1"}, expected: true},
		{rval: []string{"id1"}, op: "=", cval: []string{"id2"}, expected: false},
		{rval: []string{"id1"}, op: "=", cval: []string{"id1, id2"}, expected: false},
		{rval: []string{"id1", "id2"}, op: "=", cval: []string{"id1", "id2"}, expected: true},
		{rval: []string{"id1", "id2"}, op: "=", cval: []string{"id1", "id3"}, expected: false},
		{rval: []string{"id1"}, op: "!=", cval: []string{"id1"}, expected: false},
		{rval: []string{"id1"}, op: "!=", cval: []string{"id2"}, expected: true},
		{rval: []string{"id1"}, op: "!=", cval: []string{"id1, id2"}, expected: true},
		{rval: []string{"id1", "id2"}, op: "!=", cval: []string{"id1", "id2"}, expected: false},
		{rval: []string{"id1", "id2"}, op: "!=", cval: []string{"id1", "id3"}, expected: true},
		{rval: []string{"id1"}, op: "has", cval: "id1", expected: true},
		{rval: []string{"id2"}, op: "has", cval: "id1", expected: false},
		{rval: []string{"id1", "id2"}, op: "has", cval: "id1", expected: true},
		{rval: []string{"id2", "id3"}, op: "has", cval: "id1", expected: false},
		{rval: []string{"id1"}, op: "invalid", cval: []string{}, expected: false},
	}

	for _, test := range relTests {
		typ := &Type{Name: "type"}
		toOne := true

		if _, ok := test.rval.([]string); ok {
			toOne = false
		}
		// ty, n := GetAttrType(fmt.Sprintf("%T", test.rval))
		typ.Rels = map[string]Rel{
			"rel": {
				FromName: "rel",
				ToOne:    toOne,
				ToType:   "type",
			},
		}

		res := &SoftResource{}
		res.SetType(typ)

		if toOne {
			res.Set("rel", test.rval.(string))
		} else {
			res.Set("rel", test.rval.([]string))
		}

		filter := &Filter{
			Field: "rel",
			Op:    test.op,
			Val:   test.cval,
		}

		assert.Equal(
			test.expected,
			filter.IsAllowed(res),
			fmt.Sprintf("%v %s %v should be %v", test.cval, test.op, test.rval, test.expected),
		)
	}

	// Tests for "and" and "or"
	andOrTests := []struct {
		rvals       []any
		ops         []string
		cvals       []any
		expectedAnd bool
		expectedOr  bool
	}{
		{
			rvals:       []any{"abc", 1, true, now},
			ops:         []string{"=", "=", "=", "="},
			cvals:       []any{"abc", 1, true, now},
			expectedAnd: true,
			expectedOr:  true,
		}, {
			rvals:       []any{"abc", 1, false, now},
			ops:         []string{"=", "=", "=", "="},
			cvals:       []any{"abc", 1, true, now},
			expectedAnd: false,
			expectedOr:  true,
		}, {
			rvals:       []any{"abc", 1, false, now},
			ops:         []string{"=", "!=", "!=", "="},
			cvals:       []any{"abc", 2, true, now},
			expectedAnd: true,
			expectedOr:  true,
		}, {
			rvals:       []any{"abc", 1, false, now},
			ops:         []string{"=", "!=", "=", "!="},
			cvals:       []any{"def", 1, true, now},
			expectedAnd: false,
			expectedOr:  false,
		},
	}

	for i, test := range andOrTests {
		typ := &Type{Name: "type"}
		res := &SoftResource{}
		res.SetType(typ)

		filters := []*Filter{}

		for j := range test.rvals {
			attrName := "attr" + strconv.Itoa(j)
			ty, n := GetAttrType(fmt.Sprintf("%T", test.rvals[j]))
			_ = typ.AddAttr(
				Attr{
					Name:     attrName,
					Type:     ty,
					Nullable: n,
				},
			)

			res.Set(attrName, test.rvals[j])

			filters = append(filters, &Filter{
				Field: attrName,
				Op:    test.ops[j],
				Val:   test.cvals[j],
			})
		}

		filter := &Filter{
			Val: filters,
		}

		filter.Op = "and"
		// filter = marshalUnmarshalFilter(t, filter)
		result := filter.IsAllowed(res)
		assert.Equal(
			test.expectedAnd,
			result,
			fmt.Sprintf("'and' test %d is %t instead of %t", i, result, test.expectedAnd),
		)

		filter.Op = "or"
		// filter = marshalUnmarshalFilter(t, filter)
		result = filter.IsAllowed(res)
		assert.Equal(
			test.expectedOr,
			result,
			fmt.Sprintf("'or' test %d is %t instead of %t", i, result, test.expectedOr),
		)
	}
}

func TestFilterUnmarshaling(t *testing.T) {
	assert := assert.New(t)

	// time1, _ := time.Parse(time.RFC3339Nano, "2012-05-16T17:45:28.2539Z")
	// time2, _ := time.Parse(time.RFC3339Nano, "2013-06-24T22:03:34.8276Z")

	tests := []struct {
		name           string
		query          string
		expectedFilter Filter
		expectedError  bool
	}{
		{
			name:          "empty",
			query:         ``,
			expectedError: true,
		}, {
			name:          "null value",
			query:         `{"v":null}`,
			expectedError: false, // TODO Is this okay?
		}, {
			name:           "invalid JSON data",
			query:          `{"thisis:invalid}`,
			expectedFilter: Filter{},
			expectedError:  true,
		}, {
			name:           "invalid data",
			query:          `{"f":42}`,
			expectedFilter: Filter{},
			expectedError:  true,
		}, {
			name: "standard values",
			query: `{
				"c": "col",
				"f": "field",
				"o": "=",
				"v": "string"
			}`,
			expectedFilter: Filter{
				Field: "field",
				Op:    "=",
				Val:   "string",
				Col:   "col",
			},
			expectedError: false,
		}, {
			name: "and & or",
			query: `{
				"f": "",
				"o": "or",
				"v": [
					{
						"f": "field1",
						"o": "=",
						"v": 10
					},
					{
						"f": "",
						"o": "and",
						"v": [
							{
								"f": "field2",
								"o": ">=",
								"v": 20
							},
							{
								"f": "field3",
								"o": "<=",
								"v": 30
							}
						]
					}
				]
			}`,
			expectedFilter: Filter{
				Field: "",
				Op:    "or",
				Val: []*Filter{
					{
						Field: "field1",
						Op:    "=",
						Val:   float64(10),
					}, {
						Field: "",
						Op:    "and",
						Val: []*Filter{
							{
								Field: "field2",
								Op:    ">=",
								Val:   float64(20),
							}, {
								Field: "field3",
								Op:    "<=",
								Val:   float64(30),
							},
						},
					},
				},
			},
			expectedError: false,
		}, {
			name: "invalid or",
			query: `{
				"f": "",
				"o": "or",
				"v": "should not be a string"
			}`,
			expectedFilter: Filter{},
			expectedError:  true,
		},
	}

	for _, test := range tests {
		filter := Filter{}
		err := json.Unmarshal([]byte(test.query), &filter)

		assert.Equal(test.expectedError, err != nil, test.name)

		if !test.expectedError {
			assert.Equal(test.expectedFilter, filter, test.name)
		}
	}
}
