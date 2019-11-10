package jsonapi_test

import (
	"flag"
	"time"
)

var update = flag.Bool("update-golden-files", false, "update the golden files")

func getTime() time.Time {
	now, _ := time.Parse(time.RFC3339Nano, "2013-06-24T22:03:34.8276Z")
	return now
}

// mocktype is a fake struct that defines a JSON:API type for test purposes.
type mocktype struct {
	ID string `json:"id" api:"mocktype"`

	// Attributes
	Str    string    `json:"str" api:"attr"`
	Int    int       `json:"int" api:"attr"`
	Int8   int8      `json:"int8" api:"attr"`
	Int16  int16     `json:"int16" api:"attr"`
	Int32  int32     `json:"int32" api:"attr"`
	Int64  int64     `json:"int64" api:"attr"`
	Uint   uint      `json:"uint" api:"attr"`
	Uint8  uint8     `json:"uint8" api:"attr"`
	Uint16 uint16    `json:"uint16" api:"attr"`
	Uint32 uint32    `json:"uint32" api:"attr"`
	Uint64 uint64    `json:"uint64" api:"attr"`
	Bool   bool      `json:"bool" api:"attr"`
	Time   time.Time `json:"time" api:"attr"`
	Bytes  []byte    `json:"bytes" api:"attr"`

	// Relationships
	To1      string   `json:"to-1" api:"rel,mocktype"`
	To1From1 string   `json:"to-1-from-1" api:"rel,mocktype,to-1-from-1"`
	To1FromX string   `json:"to-1-from-x" api:"rel,mocktype,to-x-from-1"`
	ToX      []string `json:"to-x" api:"rel,mocktype"`
	ToXFrom1 []string `json:"to-x-from-1" api:"rel,mocktype,to-1-from-x"`
	ToXFromX []string `json:"to-x-from-x" api:"rel,mocktype,to-x-from-x"`
}
