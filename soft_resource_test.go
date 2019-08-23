package jsonapi_test

import (
	"testing"

	. "github.com/mfcochauxlaberge/jsonapi"

	"github.com/stretchr/testify/assert"
)

var _ Resource = (*SoftResource)(nil)

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
		Name:         "rel1",
		Type:         "type",
		ToOne:        true,
		InverseName:  "rel1",
		InverseType:  "type",
		InverseToOne: true,
	})
	sr = &SoftResource{Type: &typ}
	// TODO assert.Equal(t, &typ, sr.typ)

	// ID and type
	sr.SetID("id")
	typ2 := typ
	typ2.Name = "type2"
	sr.SetType(&typ2)
	assert.Equal(t, "id", sr.GetID())
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
			Name:         "rel1",
			Type:         "type",
			ToOne:        true,
			InverseName:  "rel1",
			InverseType:  "type",
			InverseToOne: true,
		},
		"rel2": {
			Name:         "rel2",
			Type:         "type",
			ToOne:        false,
			InverseName:  "rel1",
			InverseType:  "type",
			InverseToOne: true,
		},
	}
	for _, rel := range rels {
		sr.AddRel(rel)

		assert.Equal(t, rel, sr.Rel(rel.Name))
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
	assert.Equal(t, "", sr.GetToOne("nonexistent"))
	assert.Equal(t, []string{}, sr.GetToMany("nonexistent"))

	// Put the fields back
	for _, attr := range attrs {
		sr.AddAttr(attr)

		assert.Equal(t, attr, sr.Attr(attr.Name))
	}
	for _, rel := range rels {
		sr.AddRel(rel)

		assert.Equal(t, rel, sr.Rel(rel.Name))
	}

	// Set and get some fields
	assert.Equal(t, "", sr.Get("attr1"))
	assert.Equal(t, "", sr.GetToOne("rel1"))
	assert.Equal(t, []string{}, sr.GetToMany("rel2"))
	sr.Set("attr1", "value")
	sr.SetToOne("rel1", "id1")
	sr.SetToMany("rel2", []string{"id1", "id2"})
	assert.Equal(t, "value", sr.Get("attr1"))
	assert.Equal(t, "id1", sr.GetToOne("rel1"))
	assert.Equal(t, []string{"id1", "id2"}, sr.GetToMany("rel2"))

	// Copy
	sr2 := sr.Copy()
	assert.Equal(t, true, Equal(sr, sr2))
}
