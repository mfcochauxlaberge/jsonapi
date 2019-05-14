package jsonapi

import (
	"testing"
	"time"

	"github.com/mitchellh/copystructure"
	"github.com/stretchr/testify/assert"
)

func TestSoftCollection(t *testing.T) {
	sc := &SoftCollection{}

	// Add type
	typ := Type{Name: "thistype"}
	typ.AddAttr(Attr{
		Name: "attr1",
		Type: AttrTypeInt,
		Null: false,
	})
	typ.AddAttr(Attr{
		Name: "attr2",
		Type: AttrTypeString,
		Null: true,
	})
	typ.AddRel(Rel{
		Name:         "rel1",
		Type:         "othertype",
		ToOne:        true,
		InverseName:  "rel2",
		InverseType:  "thistype",
		InverseToOne: true,
	})
	typ.AddRel(Rel{
		Name:         "rel3",
		Type:         "othertype",
		ToOne:        false,
		InverseName:  "rel4",
		InverseType:  "thistype",
		InverseToOne: true,
	})

	// Make a copy so that modifying the original typ
	// does not modify the SoftCollection's type.
	typcopy := copystructure.Must(copystructure.Copy(typ)).(Type)
	sc.SetType(&typcopy)

	assert.Equal(t, sc.Type(), typ)

	// Modify the SoftCollection's type and the local type
	// at the same time and check whether they still are
	// the same.
	attr3 := Attr{
		Name: "attr3",
		Type: AttrTypeBool,
		Null: false,
	}
	rel5 := Rel{
		Name:         "rel5",
		Type:         "othertype",
		ToOne:        true,
		InverseName:  "rel6",
		InverseType:  "thistype",
		InverseToOne: false,
	}
	typ.AddAttr(attr3)
	sc.AddAttr(attr3)
	typ.AddRel(rel5)
	sc.AddRel(rel5)

	assert.Equal(t, sc.Type(), typ)

	// Add a SoftResource with more fields than those
	// specified in the SoftCollection.
	sr := NewSoftResource(Type{Name: "thirdtype"}, nil)
	attr4 := Attr{
		Name: "attr4",
		Type: AttrTypeUint16,
		Null: true,
	}
	sr.AddAttr(attr4)
	typ.AddAttr(attr4)
	rel7 := Rel{
		Name:         "rel7",
		Type:         "othertype",
		ToOne:        true,
		InverseName:  "rel8",
		InverseType:  "thirdtype",
		InverseToOne: true,
	}
	sr.AddRel(rel7)
	typ.AddRel(rel7)

	sc.Add(sr)

	assert.Equal(t, sc.Type(), typ)

	// Add more elements to the SoftCollection.
	sc.Add(&SoftResource{id: "res1"})
	sc.Add(&SoftResource{id: "res2"})

	assert.Equal(t, 3, sc.Len())
}

func TestSoftCollectionSort(t *testing.T) {
	now := time.Now()
	sc := &SoftCollection{}

	// Add type with some attributes
	typ := Type{Name: "thistype"}
	typ.AddAttr(Attr{
		Name: "attr1",
		Type: AttrTypeInt,
		Null: false,
	})
	typ.AddAttr(Attr{
		Name: "attr2",
		Type: AttrTypeString,
		Null: true,
	})
	typ.AddAttr(Attr{
		Name: "attr3",
		Type: AttrTypeBool,
		Null: true,
	})
	typ.AddAttr(Attr{
		Name: "attr4",
		Type: AttrTypeTime,
		Null: false,
	})
	sc.SetType(&typ)

	// Add some resources
	sr := NewSoftResource(typ, nil)
	sr.SetID("res1")
	sr.Set("attr1", 0)
	sr.Set("attr2", nil)
	sr.Set("attr3", (*bool)(nil))
	sr.Set("attr4", now)
	sc.Add(sr)

	sr = NewSoftResource(typ, nil)
	sr.SetID("res2")
	sr.Set("attr1", 0)
	sr.Set("attr2", nil)
	b1 := false
	sr.Set("attr3", &b1)
	sr.Set("attr4", now)
	sc.Add(sr)

	sr = NewSoftResource(typ, nil)
	sr.SetID("res3")
	sr.Set("attr1", 1)
	sr.Set("attr2", "")
	sr.Set("attr3", (*bool)(nil))
	sr.Set("attr4", now.Add(-time.Second))
	sc.Add(sr)

	sr = NewSoftResource(typ, nil)
	sr.SetID("res4")
	sr.Set("attr1", -1)
	sr.Set("attr2", "abc")
	sr.Set("attr3", (*bool)(nil))
	sr.Set("attr4", now.Add(time.Second))
	sc.Add(sr)

	sr = NewSoftResource(typ, nil)
	sr.SetID("res5")
	sr.Set("attr1", -1)
	sr.Set("attr2", "abc")
	b2 := true
	sr.Set("attr3", &b2)
	sr.Set("attr4", now.Add(time.Second))
	sc.Add(sr)

	sr = NewSoftResource(typ, nil)
	sr.SetID("res6")
	sr.Set("attr1", 2)
	sr.Set("attr2", "")
	sr.Set("attr3", (*bool)(nil))
	sr.Set("attr4", now.Add(time.Second))
	sc.Add(sr)

	sr = NewSoftResource(typ, nil)
	sr.SetID("res7")
	sr.Set("attr1", 2)
	sr.Set("attr2", "abc")
	b3 := true
	sr.Set("attr3", &b3)
	sr.Set("attr4", now.Add(-time.Second))
	sc.Add(sr)

	sr = NewSoftResource(typ, nil)
	sr.SetID("res8")
	sr.Set("attr1", 4)
	sr.Set("attr2", "")
	sr.Set("attr3", (*bool)(nil))
	sr.Set("attr4", now.Add(time.Second))
	sc.Add(sr)

	sr = NewSoftResource(typ, nil)
	sr.SetID("res9")
	sr.Set("attr1", -1)
	sr.Set("attr2", "def")
	sr.Set("attr3", (*bool)(nil))
	sr.Set("attr4", now.Add(time.Second))
	sc.Add(sr)

	sr = NewSoftResource(typ, nil)
	sr.SetID("res10")
	sr.Set("attr1", 4)
	sr.Set("attr2", "")
	sr.Set("attr3", (*bool)(nil))
	sr.Set("attr4", now.Add(time.Second))
	sc.Add(sr)

	// Sort the collection
	rules := []string{"-attr3", "-attr4", "attr1", "-attr2", "id"}
	sc.Sort(rules)

	assert.Equal(t, rules, sc.sort)

	// Make a ordered list of IDs
	ids := []string{}
	for i := 0; i < sc.Len(); i++ {
		ids = append(ids, sc.Elem(i).GetID())
	}

	expectedIDs := []string{
		"res5", "res7", "res2", "res9", "res4", "res6", "res10", "res8", "res1", "res3",
	}
	assert.Equal(t, expectedIDs, ids)
}
