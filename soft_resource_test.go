package jsonapi_test

import (
	"testing"
	"time"

	. "github.com/mfcochauxlaberge/jsonapi"

	"github.com/stretchr/testify/assert"
)

var _ Resource = (*SoftResource)(nil)
var _ Copier = (*SoftResource)(nil)

func TestSoftResource(t *testing.T) {
	sr := &SoftResource{}

	assert.Equal(t, map[string]Attr{}, sr.Attrs())
	assert.Equal(t, map[string]Rel{}, sr.Rels())

	// NewSoftResource
	typ := Type{Name: "type"}
	_ = typ.AddAttr(Attr{
		Name:     "attr1",
		Type:     AttrTypeString,
		Nullable: false,
	})
	_ = typ.AddRel(Rel{
		FromName: "rel1",
		FromType: "type",
		ToOne:    true,
		ToName:   "rel1",
		ToType:   "type",
		FromOne:  true,
	})
	sr = &SoftResource{Type: &typ}
	// TODO assert.Equal(t, &typ, sr.typ)

	// ID and type
	sr.SetID("id")

	typ2 := typ
	typ2.Name = "type2"
	sr.SetType(&typ2)
	assert.Equal(t, "id", sr.Get("id").(string))
	assert.Equal(t, "type2", sr.GetType().Name)

	// Attributes
	attrs := map[string]Attr{
		"attr1": {
			Name:     "attr1",
			Type:     AttrTypeString,
			Nullable: false,
		},
		"attr2": {
			Name:     "attr2",
			Type:     AttrTypeString,
			Nullable: true,
		},
	}
	for _, attr := range attrs {
		sr.AddAttr(attr)

		assert.Equal(t, attr, sr.Attr(attr.Name))
	}

	assert.Equal(t, attrs, sr.Attrs())

	// Relationships
	rels := map[string]Rel{
		"rel1": {
			FromName: "rel1",
			FromType: "type",
			ToOne:    true,
			ToName:   "rel1",
			ToType:   "type",
			FromOne:  true,
		},
		"rel2": {
			FromName: "rel2",
			FromType: "type",
			ToOne:    false,
			ToName:   "rel1",
			ToType:   "type",
			FromOne:  true,
		},
	}
	for _, rel := range rels {
		sr.AddRel(rel)

		assert.Equal(t, rel, sr.Rel(rel.FromName))
	}

	assert.Equal(t, rels, sr.Rels())

	sr.RemoveField("attr1")
	assert.Equal(t, Attr{}, sr.Attr("attr1"))
	sr.RemoveField("attr2")
	assert.Equal(t, map[string]Attr{}, sr.Attrs())

	sr.RemoveField("rel1")
	assert.Equal(t, Rel{}, sr.Rel("rel1"))
	sr.RemoveField("rel2")
	assert.Equal(t, map[string]Rel{}, sr.Rels())

	assert.Equal(t, nil, sr.Get("nonexistent"))

	// Put the fields back
	for _, attr := range attrs {
		sr.AddAttr(attr)

		assert.Equal(t, attr, sr.Attr(attr.Name))
	}

	for _, rel := range rels {
		sr.AddRel(rel)

		assert.Equal(t, rel, sr.Rel(rel.FromName))
	}

	// Set and get some fields
	assert.Equal(t, "", sr.Get("attr1"))
	assert.Equal(t, "", sr.Get("rel1").(string))
	assert.Equal(t, []string{}, sr.Get("rel2").([]string))
	sr.Set("attr1", "value")
	sr.Set("rel1", "id1")
	sr.Set("rel2", []string{"id1", "id2"})
	assert.Equal(t, "value", sr.Get("attr1"))
	assert.Equal(t, "id1", sr.Get("rel1").(string))
	assert.Equal(t, []string{"id1", "id2"}, sr.Get("rel2").([]string))

	// Set a nullable attribute to nil
	_ = sr.Type.AddAttr(Attr{
		Name:     "nullable-str",
		Type:     AttrTypeString,
		Nullable: true,
	})

	assert.Nil(t, sr.Get("nullable-str"))

	str := "abc"
	sr.Set("nullable-str", &str)
	assert.Equal(t, &str, sr.Get("nullable-str"))
	sr.Set("nullable-str", nil)
	assert.Nil(t, sr.Get("nullable-str"))
	assert.Equal(t, (*string)(nil), sr.Get("nullable-str"))

	// Getting the value of an unset field returns
	// the zero value of the type.
	sr = &SoftResource{}

	sr.AddAttr(Attr{
		Name:     "zero-str",
		Type:     AttrTypeString,
		Nullable: false,
	})
	assert.Equal(t, "", sr.Get("zero-str"))

	sr.AddAttr(Attr{
		Name:     "zero-str-null",
		Type:     AttrTypeString,
		Nullable: true,
	})
	assert.Equal(t, (*string)(nil), sr.Get("zero-str-null"))

	sr.AddRel(Rel{
		FromName: "zero-to-one",
		ToOne:    true,
	})
	assert.Equal(t, "", sr.Get("zero-to-one"))

	sr.AddRel(Rel{
		FromName: "zero-to-many",
		ToOne:    false,
	})
	assert.Equal(t, []string{}, sr.Get("zero-to-many"))
}

func TestSoftResourceNew(t *testing.T) {
	assert := assert.New(t)

	typ, _ := BuildType(mocktype{})
	sr := &SoftResource{}
	sr.Type = &typ

	// Modify the SoftResource object
	sr.SetID("id")
	sr.Set("str", "abc123")
	sr.Set("int", 42)

	nsr := sr.New()

	// The new
	assert.Equal("", nsr.Get("id").(string))
	assert.Equal("", nsr.Get("str"))
	assert.Equal(0, nsr.Get("int"))
}

func TestSoftResourceCopy(t *testing.T) {
	assert := assert.New(t)

	now, _ := time.Parse(time.RFC3339, "2019-11-19T23:17:01-05:00")

	sr := &SoftResource{}

	// Attributes
	attrs := map[string]any{
		"string":     "abc",
		"int":        42,
		"int8":       8,
		"int16":      16,
		"int32":      32,
		"int64":      64,
		"uint":       42,
		"uint8":      8,
		"uint16":     16,
		"uint32":     32,
		"uint64":     64,
		"bool":       true,
		"time.Time":  now,
		"[]uint8":    []byte{'a', 'b', 'c'},
		"*string":    ptr("abc"),
		"*int":       ptr(42),
		"*int8":      ptr(8),
		"*int16":     ptr(16),
		"*int32":     ptr(32),
		"*int64":     ptr(64),
		"*uint":      ptr(42),
		"*uint8":     ptr(8),
		"*uint16":    ptr(16),
		"*uint32":    ptr(32),
		"*uint64":    ptr(64),
		"*bool":      ptr(true),
		"*time.Time": ptr(now),
		"*[]uint8":   ptr([]byte{'a', 'b', 'c'}),
	}

	for t, v := range attrs {
		typ, null := GetAttrType(t)

		sr.AddAttr(Attr{
			Name:     t,
			Type:     typ,
			Nullable: null,
		})

		sr.Set(t, v)
	}

	// Special cases
	sr.AddAttr(Attr{
		Name:     "nil-*[]byte",
		Type:     AttrTypeBytes,
		Nullable: true,
	})

	sr.Set("nil-*[]byte", (*[]byte)(nil))

	// Relationships
	sr.AddRel(Rel{
		FromName: "to-one",
		ToOne:    true,
	})
	sr.Set("to-one", "id1")

	sr.AddRel(Rel{
		FromName: "to-many",
		ToOne:    false,
	})
	sr.Set("to-many", []string{"id2", "id3"})

	// Copy
	sr2 := sr.Copy()
	assert.Equal(true, Equal(sr, sr2))
}

func TestSoftResourceMeta(t *testing.T) {
	assert := assert.New(t)

	typ, _ := BuildType(mocktype{})
	sr := &SoftResource{}
	sr.Type = &typ
	sr.SetID("id")

	meta := Meta(map[string]any{
		"key1": "a string",
		"key2": 200,
		"key3": false,
		"key4": getTime(),
	})

	// Add some meta values
	sr.SetMeta(meta)

	// The new
	assert.Equal(meta, sr.Meta())
}

func TestSoftResourceGetSetID(t *testing.T) {
	assert := assert.New(t)

	sr := &SoftResource{}
	sr.Set("id", "abc123")

	assert.Equal("abc123", sr.Get("id"))
}
