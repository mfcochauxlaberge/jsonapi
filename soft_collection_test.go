package jsonapi_test

import (
	"testing"

	. "github.com/mfcochauxlaberge/jsonapi"

	"github.com/mitchellh/copystructure"
	"github.com/stretchr/testify/assert"
)

func TestSoftCollection(t *testing.T) {
	sc := &SoftCollection{}

	// Add type
	typ := Type{Name: "thistype"}
	_ = typ.AddAttr(Attr{
		Name:     "attr1",
		Type:     AttrTypeInt,
		Nullable: false,
	})
	_ = typ.AddAttr(Attr{
		Name:     "attr2",
		Type:     AttrTypeString,
		Nullable: true,
	})
	_ = typ.AddRel(Rel{
		FromName: "rel1",
		FromType: "thistype",
		ToOne:    true,
		ToName:   "rel2",
		ToType:   "othertype",
		FromOne:  true,
	})
	_ = typ.AddRel(Rel{
		FromName: "rel3",
		FromType: "thistype",
		ToOne:    false,
		ToName:   "rel4",
		ToType:   "othertype",
		FromOne:  true,
	})

	// Make a copy so that modifying the original typ
	// does not modify the SoftCollection's type.
	typcopy := copystructure.Must(copystructure.Copy(typ)).(Type)
	sc.SetType(&typcopy)

	assert.Equal(t, sc.Type, &typ)

	// Modify the SoftCollection's type and the local type
	// at the same time and check whether they still are
	// the same.
	attr3 := Attr{
		Name:     "attr3",
		Type:     AttrTypeBool,
		Nullable: false,
	}
	rel5 := Rel{
		FromName: "rel5",
		FromType: "thistype",
		ToOne:    true,
		ToName:   "rel6",
		ToType:   "othertype",
		FromOne:  false,
	}
	_ = typ.AddAttr(attr3)
	_ = sc.AddAttr(attr3)
	_ = typ.AddRel(rel5)
	_ = sc.AddRel(rel5)

	assert.Equal(t, sc.Type, &typ)

	// Add a SoftResource with more fields than those
	// specified in the SoftCollection.
	sr := &SoftResource{Type: &Type{Name: "thirdtype"}}
	attr4 := Attr{
		Name:     "attr4",
		Type:     AttrTypeUint16,
		Nullable: true,
	}
	sr.AddAttr(attr4)
	_ = typ.AddAttr(attr4)
	rel7 := Rel{
		FromName: "rel7",
		FromType: "thirdtype",
		ToOne:    true,
		ToName:   "rel8",
		ToType:   "othertype",
		FromOne:  true,
	}
	sr.AddRel(rel7)
	_ = typ.AddRel(rel7)
	rel8 := Rel{
		FromName: "rel8",
		FromType: "thirdtype",
		ToOne:    false,
		ToName:   "rel9",
		ToType:   "othertype",
		FromOne:  true,
	}
	sr.AddRel(rel8)
	_ = typ.AddRel(rel8)

	sc.Add(sr)

	assert.Equal(t, sc.Type, &typ)

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

	sc.Type.Name = "type1"
	_ = sc.Type.AddAttr(Attr{
		Name:     "attr1",
		Type:     AttrTypeString,
		Nullable: false,
	})
	_ = sc.Type.AddAttr(Attr{
		Name:     "attr2",
		Type:     AttrTypeInt,
		Nullable: true,
	})
	_ = sc.Type.AddRel(Rel{
		FromName: "rel1",
		ToOne:    true,
		ToType:   "type2",
	})

	sr := &SoftResource{}
	sr.SetType(sc.Type)
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

func TestSoftCollectionMiscellaneous(t *testing.T) {
	assert := assert.New(t)

	sc := &SoftCollection{}
	assert.Nil(sc.At(99), "nonexistent element")
}
