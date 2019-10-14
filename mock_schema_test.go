package jsonapi_test

import (
	"time"

	. "github.com/mfcochauxlaberge/jsonapi"
)

// newMockSchema ...
func newMockSchema() *Schema {
	schema := &Schema{}

	typ := MustBuildType(mockType1{})
	_ = schema.AddType(typ)
	typ = MustBuildType(mockType2{})
	_ = schema.AddType(typ)
	typ = MustBuildType(mockType3{})
	_ = schema.AddType(typ)

	for t, typ := range schema.Types {
		for r, rel := range typ.Rels {
			invType := schema.GetType(rel.ToType)
			rel := schema.Types[t].Rels[r]
			rel.FromOne = invType.Rels[rel.ToName].ToOne
			schema.Types[t].Rels[r] = rel
		}
	}

	errs := schema.Check()
	if len(errs) > 0 {
		panic(errs[0])
	}

	return schema
}

// mockType1 ...
type mockType1 struct {
	ID string `json:"id" api:"mocktypes1"`

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

	// Relationships
	ToOne          string   `json:"to-one" api:"rel,mocktypes2"`
	ToOneFromOne   string   `json:"to-one-from-one" api:"rel,mocktypes2,to-one-from-one"`
	ToOneFromMany  string   `json:"to-one-from-many" api:"rel,mocktypes2,to-many-from-one"`
	ToMany         []string `json:"to-many" api:"rel,mocktypes2"`
	ToManyFromOne  []string `json:"to-many-from-one" api:"rel,mocktypes2,to-one-from-many"`
	ToManyFromMany []string `json:"to-many-from-many" api:"rel,mocktypes2,to-many-from-many"`
}

// mockType2 ...
type mockType2 struct {
	ID string `json:"id" api:"mocktypes2"`

	// Attributes
	StrPtr    *string    `json:"strptr" api:"attr"`
	IntPtr    *int       `json:"intptr" api:"attr"`
	Int8Ptr   *int8      `json:"int8ptr" api:"attr"`
	Int16Ptr  *int16     `json:"int16ptr" api:"attr"`
	Int32Ptr  *int32     `json:"int32ptr" api:"attr"`
	Int64Ptr  *int64     `json:"int64ptr" api:"attr"`
	UintPtr   *uint      `json:"uintptr" api:"attr"`
	Uint8Ptr  *uint8     `json:"uint8ptr" api:"attr"`
	Uint16Ptr *uint16    `json:"uint16ptr" api:"attr"`
	Uint32Ptr *uint32    `json:"uint32ptr" api:"attr"`
	Uint64Ptr *uint64    `json:"uint64ptr" api:"attr"`
	BoolPtr   *bool      `json:"boolptr" api:"attr"`
	TimePtr   *time.Time `json:"timeptr" api:"attr"`

	// Relationships
	ToOneFromOne   string   `json:"to-one-from-one" api:"rel,mocktypes1,to-one-from-one"`
	ToOneFromMany  string   `json:"to-one-from-many" api:"rel,mocktypes1,to-many-from-one"`
	ToManyFromOne  []string `json:"to-many-from-one" api:"rel,mocktypes1,to-one-from-many"`
	ToManyFromMany []string `json:"to-many-from-many" api:"rel,mocktypes1,to-many-from-many"`
}

// mockType3 ...
type mockType3 struct {
	ID string `json:"id" api:"mocktypes3"`

	// Attributes
	Attr1 string `json:"attr1" api:"attr"`
	Attr2 int    `json:"attr2" api:"attr"`

	// Relationships
	Rel1 string   `json:"rel1" api:"rel,mocktypes1"`
	Rel2 []string `json:"rel2" api:"rel,mocktypes1"`
}
