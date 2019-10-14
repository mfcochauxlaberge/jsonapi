package jsonapi_test

import (
	"strings"
	"time"
)

func makeOneLineNoSpaces(str string) string {
	str = strings.Replace(str, "\t", "", -1)
	str = strings.Replace(str, "\n", "", -1)
	return strings.Replace(str, " ", "", -1)
}

func ptr(v interface{}) interface{} {
	switch c := v.(type) {
	// String
	case string:
		return &c
	// Integers
	case int:
		return &c
	case int8:
		return &c
	case int16:
		return &c
	case int32:
		return &c
	case int64:
		return &c
	case uint:
		return &c
	case uint8:
		return &c
	case uint16:
		return &c
	case uint32:
		return &c
	case uint64:
		return &c
	// Bool
	case bool:
		return &c
	// time.Time
	case time.Time:
		return &c
	// []byte
	case []byte:
		return &c
	}
	return nil
}

func nilptr(t string) interface{} {
	switch t {
	// String
	case "string":
		var p *string
		return p
	// Integers
	case "int":
		var p *int
		return p
	case "int8":
		var p *int8
		return p
	case "int16":
		var p *int16
		return p
	case "int32":
		var p *int32
		return p
	case "int64":
		var p *int64
		return p
	case "uint":
		var p *uint
		return p
	case "uint8":
		var p *uint8
		return p
	case "uint16":
		var p *uint16
		return p
	case "uint32":
		var p *uint32
		return p
	case "uint64":
		var p *uint64
		return p
	// Bool
	case "bool":
		var p *bool
		return p
	// time.Time
	case "time.Time":
		var p *time.Time
		return p
	// []byte
	case "[]byte":
		var p *[]byte
		return p
	}
	return nil
}
