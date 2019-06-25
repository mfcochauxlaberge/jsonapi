package jsonapi_test

import (
	"testing"
	"time"

	. "github.com/mfcochauxlaberge/jsonapi"

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

	assert.Equal(t, sc.GetType(), &typ)

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

	assert.Equal(t, sc.GetType(), &typ)

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

	assert.Equal(t, sc.GetType(), &typ)

	// Add more elements to the SoftCollection.
	sr = &SoftResource{}
	sr.SetID("res1")
	sc.Add(sr)
	sr = &SoftResource{}
	sr.SetID("res2")
	sc.Add(sr)

	assert.Equal(t, 3, sc.Len())

	// Remove an element.
	sc.Remove("res1")
	sc.Remove("res99")

	assert.Equal(t, 2, sc.Len())
}

func TestSoftCollectionResource(t *testing.T) {
	sc := &SoftCollection{}
	sc.SetType(&Type{})

	sc.GetType().Name = "type1"
	sc.GetType().AddAttr(Attr{
		Name: "attr1",
		Type: AttrTypeString,
		Null: false,
	})
	sc.GetType().AddAttr(Attr{
		Name: "attr2",
		Type: AttrTypeInt,
		Null: true,
	})
	sc.GetType().AddRel(Rel{
		Name:  "rel1",
		Type:  "type2",
		ToOne: true,
	})

	sr := &SoftResource{}
	sr.SetType(sc.GetType())
	sr.SetID("res1")
	sr.Set("attr", "value1")
	sc.Add(sr)

	// Resource with all fields
	assert.Equal(t, sr, sc.Resource("res1", nil))

	// Resource with some fields
	// TODO Fix this test. It seems like defining any set of
	// fields will make the assert pass.
	assert.Equal(t, sr, sc.Resource("res1", []string{"attr2", "rel1"}))

	// Resource not found
	assert.Equal(t, nil, sc.Resource("notfound", nil))
}

func TestSoftCollectionSort(t *testing.T) {
	now := time.Now()
	sc := &SoftCollection{}

	// Add type with some attributes.
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

	// Add some resources.
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
	sr.Set("attr2", ptr(""))
	sr.Set("attr3", (*bool)(nil))
	sr.Set("attr4", now.Add(-time.Second))
	sc.Add(sr)

	sr = NewSoftResource(typ, nil)
	sr.SetID("res4")
	sr.Set("attr1", -1)
	sr.Set("attr2", ptr("abc"))
	sr.Set("attr3", (*bool)(nil))
	sr.Set("attr4", now.Add(time.Second))
	sc.Add(sr)

	sr = NewSoftResource(typ, nil)
	sr.SetID("res5")
	sr.Set("attr1", -1)
	sr.Set("attr2", ptr("abc"))
	b2 := true
	sr.Set("attr3", &b2)
	sr.Set("attr4", now.Add(time.Second))
	sc.Add(sr)

	sr = NewSoftResource(typ, nil)
	sr.SetID("res6")
	sr.Set("attr1", 2)
	sr.Set("attr2", ptr(""))
	sr.Set("attr3", (*bool)(nil))
	sr.Set("attr4", now.Add(time.Second))
	sc.Add(sr)

	sr = NewSoftResource(typ, nil)
	sr.SetID("res7")
	sr.Set("attr1", 2)
	sr.Set("attr2", ptr("abc"))
	b3 := true
	sr.Set("attr3", &b3)
	sr.Set("attr4", now.Add(-time.Second))
	sc.Add(sr)

	sr = NewSoftResource(typ, nil)
	sr.SetID("res8")
	sr.Set("attr1", 4)
	sr.Set("attr2", ptr(""))
	sr.Set("attr3", (*bool)(nil))
	sr.Set("attr4", now.Add(time.Second))
	sc.Add(sr)

	sr = NewSoftResource(typ, nil)
	sr.SetID("res9")
	sr.Set("attr1", -1)
	sr.Set("attr2", ptr("def"))
	sr.Set("attr3", (*bool)(nil))
	sr.Set("attr4", now.Add(time.Second))
	sc.Add(sr)

	sr = NewSoftResource(typ, nil)
	sr.SetID("res10")
	sr.Set("attr1", 4)
	sr.Set("attr2", ptr(""))
	sr.Set("attr3", (*bool)(nil))
	sr.Set("attr4", now.Add(time.Second))
	sc.Add(sr)

	// Sort the collection.
	rules := []string{"-attr3", "-attr4", "attr1", "-attr2", "id"}
	sc.Sort(rules)

	// Make an ordered list of IDs.
	ids := []string{}
	for i := 0; i < sc.Len(); i++ {
		ids = append(ids, sc.Elem(i).GetID())
	}

	expectedIDs := []string{
		"res5", "res7", "res2", "res9", "res4", "res6", "res10", "res8", "res1", "res3",
	}
	assert.Equal(t, expectedIDs, ids)

	// Sort with an empty list of sorting rules.
	sc.Sort([]string{})

	ids = []string{}
	for i := 0; i < sc.Len(); i++ {
		ids = append(ids, sc.Elem(i).GetID())
	}

	expectedIDs = []string{
		"res1", "res10", "res2", "res3", "res4", "res5", "res6", "res7", "res8", "res9",
	}
	assert.Equal(t, expectedIDs, ids)
}
