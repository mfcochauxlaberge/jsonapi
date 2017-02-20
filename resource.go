package jsonapi

// Resource ...
type Resource interface {
	// Structure
	IDAndType() (string, string)
	Attrs() []Attr
	Rels() map[string]Rel
	Attr(key string) Attr
	Rel(key string) Rel
	// AttrPtrs(attrs []Attr) ([]string, []interface{})
	// ApplyAttrPtrs(keys []string, vals []interface{})
	New() Resource

	// Read
	// GetID() string
	Get(key string) interface{}
	// GetString(key string) string
	// GetStringPtr(key string) *string
	// GetInt(key string) int
	// GetIntPtr(key string) *int
	// GetInt8(key string) int8
	// GetInt8Ptr(key string) *int8
	// GetInt16(key string) int16
	// GetInt16Ptr(key string) *int16
	// GetInt32(key string) int32
	// GetInt32Ptr(key string) *int32
	// GetInt64(key string) int64
	// GetInt64Ptr(key string) *int64
	// GetUint(key string) uint
	// GetUintPtr(key string) *uint
	// GetUint8(key string) uint8
	// GetUint8Ptr(key string) *uint8
	// GetUint16(key string) uint16
	// GetUint16Ptr(key string) *uint16
	// GetUint32(key string) uint32
	// GetUint32Ptr(key string) *uint32
	// GetBool(key string) bool
	// GetBoolPtr(key string) *bool
	// GetTime(key string) time.Time
	// GetTimePtr(key string) *time.Time

	// Update
	SetID(id string)
	Set(key string, val interface{})
	// SetString(key string, val string)
	// SetStringPtr(key string, val *string)
	// SetInt(key string, val int)
	// SetIntPtr(key string, val *int)
	// SetInt8(key string, val int8)
	// SetInt8Ptr(key string, val *int8)
	// SetInt16(key string, val int16)
	// SetInt16Ptr(key string, val *int16)
	// SetInt32(key string, val int32)
	// SetInt32Ptr(key string, val *int32)
	// SetInt64(key string, val int64)
	// SetInt64Ptr(key string, val *int64)
	// SetUint(key string, val uint)
	// SetUintPtr(key string, val *uint)
	// SetUint8(key string, val uint8)
	// SetUint8Ptr(key string, val *uint8)
	// SetUint16(key string, val uint16)
	// SetUint16Ptr(key string, val *uint16)
	// SetUint32(key string, val uint32)
	// SetUint32Ptr(key string, val *uint32)
	// SetBool(key string, val bool)
	// SetBoolPtr(key string, val *bool)
	// SetTime(key string, val time.Time)
	// SetTimePtr(key string, val *time.Time)

	// Read relationship
	GetToOne(key string) string
	GetToMany(key string) []string

	// Update relationship
	SetToOne(key string, rel string)
	SetToMany(key string, rels []string)

	// Validate
	Validate(keys []string) []error

	// JSON
	MarshalJSONOptions(opts *Options) ([]byte, error)
	UnmarshalJSON(payload []byte) error
}
