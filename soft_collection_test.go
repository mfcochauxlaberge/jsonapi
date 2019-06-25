package jsonapi_test

import (
	"fmt"
	"strconv"
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
	var (
		now = time.Now()
		sc  = &SoftCollection{}
	)

	attrs := []struct {
		typ  string
		vals [5]interface{}
	}{
		{
			typ:  "string",
			vals: [5]interface{}{"", "a", "b", "b", "c"},
		}, {
			typ:  "int",
			vals: [5]interface{}{-1, 0, 1, 1, 2},
		}, {
			typ:  "int8",
			vals: [5]interface{}{-1, 0, 1, 1, 2},
		}, {
			typ:  "int16",
			vals: [5]interface{}{-1, 0, 1, 1, 2},
		}, {
			typ:  "int32",
			vals: [5]interface{}{-1, 0, 1, 1, 2},
		}, {
			typ:  "int64",
			vals: [5]interface{}{-1, 0, 1, 1, 2},
		}, {
			typ:  "uint",
			vals: [5]interface{}{0, 1, 2, 2, 3},
		}, {
			typ:  "uint8",
			vals: [5]interface{}{0, 1, 2, 2, 3},
		}, {
			typ:  "uint16",
			vals: [5]interface{}{0, 1, 2, 2, 3},
		}, {
			typ:  "uint32",
			vals: [5]interface{}{0, 1, 2, 2, 3},
		}, {
			typ:  "uint64",
			vals: [5]interface{}{0, 1, 2, 2, 3},
		}, {
			typ:  "bool",
			vals: [5]interface{}{false, true, true, false, false},
		}, {
			typ: "time.Time",
			vals: [5]interface{}{
				now.Add(-time.Second),
				now,
				now.Add(time.Second),
				now.Add(time.Second),
				now.Add(2 * time.Second),
			},
		},
	}

	// Add attributes to type
	typ := Type{Name: "type"}
	for i, t := range attrs {
		ti, null := GetAttrType(t.typ)
		typ.AddAttr(Attr{
			Name: "attr" + strconv.Itoa(i),
			Type: ti,
			Null: null,
		})
	}
	sc.SetType(&typ)

	// Add resources
	for j := 0; j < 5; j++ {
		sr := NewSoftResource(typ, nil)
		sr.SetID("id" + strconv.Itoa(j))
		for i, t := range attrs {
			sr.Set("attr"+strconv.Itoa(i), t.vals[j])
		}
		sc.Add(sr)
	}

	for j := 0; j < 5; j++ {
		res := sc.Elem(j)
		fmt.Printf("Resource: %s (%s)\n", res.GetID(), res.GetType().Name)
		for _, field := range res.GetType().Fields() {
			fmt.Printf("  %s: %q (%T)\n", field, res.Get(field), res.Get(field))
		}
	}

	// Sort
	rules := []string{}
	for i := 0; i < sc.Len(); i++ {
		reverse := ""
		if i%3 == 0 {
			reverse = "-"
		}
		rules = append(rules, reverse+"attr"+strconv.Itoa(i))
	}
	sc.Sort(rules)

	// Sorted IDs from the collection
	ids := []string{}
	for i := 0; i < sc.Len(); i++ {
		ids = append(ids, sc.Elem(i).GetID())
	}

	expectedIDs := []string{
		"id4", "id2", "id3", "id1", "id0",
	}
	assert.Equal(t, expectedIDs, ids, fmt.Sprintf("rules: %v", rules))

	// // Sort with an empty list of sorting rules.
	// sc.Sort([]string{})

	// ids = []string{}
	// for i := 0; i < sc.Len(); i++ {
	// 	ids = append(ids, sc.Elem(i).GetID())
	// }

	// expectedIDs = []string{
	// 	"id1", "id10", "id2", "id3", "id4", "id5", "id6", "id7", "id8", "id9",
	// }
	// assert.Equal(t, expectedIDs, ids)
}
