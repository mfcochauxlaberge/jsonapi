package jsonapi

import (
	"testing"

	"github.com/mfcochauxlaberge/tchek"
)

func TestSoftResource(t *testing.T) {
	sr := &SoftResource{}

	tchek.AreEqual(t, "no attributes", map[string]Attr{}, sr.Attrs())
	tchek.AreEqual(t, "no relationships", map[string]Rel{}, sr.Rels())

	// ID and type
	sr.SetID("id")
	sr.SetType("type")
	tchek.AreEqual(t, "get id", "id", sr.GetID())
	tchek.AreEqual(t, "get type", "type", sr.GetType())

	// Attributes
	attrs := map[string]Attr{
		"attr1": Attr{
			Name: "attr1",
			Type: AttrTypeString,
			Null: false,
		},
		"attr2": Attr{
			Name: "attr2",
			Type: AttrTypeString,
			Null: true,
		},
	}
	for _, attr := range attrs {
		sr.AddAttr(attr)

		tchek.AreEqual(t, "get an attribute", attr, sr.Attr(attr.Name))
	}
	tchek.AreEqual(t, "list all attributes", attrs, sr.Attrs())

	// Relationships
	rels := map[string]Rel{
		"rel1": Rel{
			Name:         "rel1",
			Type:         "type",
			ToOne:        true,
			InverseName:  "rel1",
			InverseType:  "type",
			InverseToOne: true,
		},
		"rel2": Rel{
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

		tchek.AreEqual(t, "get an attribute", rel, sr.Rel(rel.Name))
	}
	tchek.AreEqual(t, "list all attributes", rels, sr.Rels())

	sr.RemoveField("attr1")
	tchek.AreEqual(t, "can't get removed attribute", Attr{}, sr.Attr("attr1"))
	sr.RemoveField("attr2")
	tchek.AreEqual(t, "all attributes are removed", map[string]Attr{}, sr.Attrs())

	sr.RemoveField("rel1")
	tchek.AreEqual(t, "can't get removed relationship", Rel{}, sr.Rel("rel1"))
	sr.RemoveField("rel2")
	tchek.AreEqual(t, "all relationships are removed", map[string]Rel{}, sr.Rels())

	tchek.AreEqual(t, "get an nonexistent value", nil, sr.Get("nonexistent"))
	tchek.AreEqual(t, "get an nonexistent to-one rel", "", sr.GetToOne("nonexistent"))
	tchek.AreEqual(t, "get an nonexistent to-many rel", []string{}, sr.GetToMany("nonexistent"))

	// Put the fields back
	for _, attr := range attrs {
		sr.AddAttr(attr)

		tchek.AreEqual(t, "get an attribute", attr, sr.Attr(attr.Name))
	}
	for _, rel := range rels {
		sr.AddRel(rel)

		tchek.AreEqual(t, "get an attribute", rel, sr.Rel(rel.Name))
	}

	// Set and get some fields
	tchek.AreEqual(t, "get a zero value 1", "", sr.Get("attr1"))
	tchek.AreEqual(t, "get a zero value 2", "", sr.GetToOne("rel1"))
	tchek.AreEqual(t, "get a zero value 3", []string{}, sr.GetToMany("rel2"))
	sr.Set("attr1", "value")
	sr.SetToOne("rel1", "id1")
	sr.SetToMany("rel2", []string{"id1", "id2"})
	tchek.AreEqual(t, "get a value 1", "value", sr.Get("attr1"))
	tchek.AreEqual(t, "get a value 2", "id1", sr.GetToOne("rel1"))
	tchek.AreEqual(t, "get a value 3", []string{"id1", "id2"}, sr.GetToMany("rel2"))

	// Copy
	sr2 := sr.Copy()
	tchek.AreEqual(t, "copy is equal", true, Equal(sr, sr2))
}
