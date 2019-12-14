package jsonapi

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
	v, _ := m[key].(string)
	return v
}

// GetInt returns the int associated with the given key.
//
// 0 is returned if the key could not be found or the type is not compatible.
func (m Meta) GetInt(key string) int {
	v, _ := m[key].(int)
	return v
}
