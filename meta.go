package jsonapi

import (
	"fmt"
	"time"
)

// A Meter represents the interface of an object that has meta information about
// itself.
//
// It is useful for a struct that represents a resource type to implement this
// interface to have a meta property as part of its JSON output.
type Meter interface {
	Meta() Meta
	SetMeta(Meta)
}

type Meta map[string]interface{}

// Has reports whether the Meta map contains or not the given key.
func (m Meta) Has(key string) bool {
	_, ok := m[key]
	return ok
}

// GetString returns the string associated with the given key.
//
// An empty string is returned if the key could not be found or the type is not
// compatible.
func (m Meta) GetString(key string) string {
	return fmt.Sprint(m[key])
}

// GetInt returns the int associated with the given key.
//
// 0 is returned if the key could not be found or the type is not compatible.
func (m Meta) GetInt(key string) int {
	v, _ := m[key].(int)
	return v
}

// GetBool returns the bool associated with the given key.
//
// False is returned if the key could not be found or the type is not
// compatible. The "true" JSON keyword is the only value that will make this
// method return true.
func (m Meta) GetBool(key string) bool {
	b, _ := m[key].(bool)
	return b
}

// GetTime returns the time.Time associated with the given key.
//
// time.Time{} is returned is the value associated with the key could not be
// found or could not be parsed with time.RFC3339Nano.
func (m Meta) GetTime(key string) time.Time {
	t := time.Time{}

	if s, ok := m[key].(string); ok {
		t, _ = time.Parse(time.RFC3339Nano, s)
	}

	return t
}
