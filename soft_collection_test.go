package jsonapi_test

import (
	"fmt"
	"sort"
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
		Name:     "attr1",
		Type:     AttrTypeInt,
		Nullable: false,
	})
	typ.AddAttr(Attr{
		Name:     "attr2",
		Type:     AttrTypeString,
		Nullable: true,
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
		Name:     "attr3",
		Type:     AttrTypeBool,
		Nullable: false,
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
	sr := &SoftResource{Type: &Type{Name: "thirdtype"}}
	attr4 := Attr{
		Name:     "attr4",
		Type:     AttrTypeUint16,
		Nullable: true,
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
		Name:     "attr1",
		Type:     AttrTypeString,
		Nullable: false,
	})
	sc.GetType().AddAttr(Attr{
		Name:     "attr2",
		Type:     AttrTypeInt,
		Nullable: true,
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

func TestSoftCollectionRange(t *testing.T) {
	assert := assert.New(t)

	// Collection
	col := SoftCollection{}
	col.SetType(&Type{})
	col.AddAttr(Attr{
		Name:     "attr1",
		Type:     AttrTypeString,
		Nullable: false,
	})
	col.AddAttr(Attr{
		Name:     "attr2",
		Type:     AttrTypeInt,
		Nullable: false,
	})

	resources := []struct {
		id     string
		fields map[string]interface{}
	}{
		{
			id: "res1",
			fields: map[string]interface{}{
				"attr1": "string1",
				"attr2": 2,
			},
		}, {
			id: "res2",
			fields: map[string]interface{}{
				"attr1": "string2",
				"attr2": 2,
			},
		}, {
			id: "res3",
			fields: map[string]interface{}{
				"attr1": "string2",
				"attr2": 0,
			},
		}, {
			id: "res4",
			fields: map[string]interface{}{
				"attr1": "string2",
				"attr2": 2,
			},
		}, {
			id: "res5",
			fields: map[string]interface{}{
				"attr1": "string3",
				"attr2": 2,
			},
		}, {
			id: "res6",
			fields: map[string]interface{}{
				"attr1": "string3",
				"attr2": 4,
			},
		}, {
			id: "res7",
			fields: map[string]interface{}{
				"attr1": "string4",
				"attr2": 0,
			},
		}, {
			id: "res8",
			fields: map[string]interface{}{
				"attr1": "string5",
				"attr2": 2,
			},
		},
	}

	for _, res := range resources {
		sr := &SoftResource{}
		sr.SetType(col.GetType())
		sr.SetID(res.id)
		for field, val := range res.fields {
			sr.Set(field, val)
		}
		col.Add(sr)
	}

	// Range test 1
	rangd := col.Range(
		// IDs
		[]string{},
		// Filter
		nil,
		// Sort
		[]string{},
		// Fields
		[]string{},
		// PageSize
		10,
		// PageNumber
		0,
	)

	expectedIDs := []string{
		"res1", "res2", "res3", "res4", "res5", "res6", "res7", "res8",
	}
	ids := []string{}
	for i := 0; i < len(rangd); i++ {
		ids = append(ids, rangd[i].GetID())
	}
	assert.Equal(expectedIDs, ids, "range of IDs (1)")

	// Range test 2
	rangd = col.Range(
		// IDs
		[]string{"res1", "res2", "res3", "res4", "res5", "res6"},
		// Filter
		&Filter{Field: "attr2", Op: "=", Val: 2},
		// Sort
		[]string{"-attr1"},
		// Fields
		[]string{"attr1", "attr2"},
		// PageSize
		2,
		// PageNumber
		0,
	)

	expectedIDs = []string{"res5", "res2"}
	ids = []string{}
	for i := 0; i < len(rangd); i++ {
		ids = append(ids, rangd[i].GetID())
	}
	assert.Equal(expectedIDs, ids, "range of IDs (2)")

	// Range test 3
	assert.Equal(0, len(col.Range(nil, nil, nil, nil, 1, 100)), "range of IDs (3)")
}
func TestSoftCollectionSort(t *testing.T) {
	assert := assert.New(t)

	var (
		now = time.Now()
		sc  = &SoftCollection{}
	)

	// A collection of resources will be created and
	// one attribute will be added for each entry from
	// the following slice.
	// The point is to provoke all possible scenarios
	// for each attribute type.
	attrs := []struct {
		vals [2]interface{}
	}{
		// non-nullable
		{vals: [2]interface{}{"", "a"}},
		{vals: [2]interface{}{int(-1), int(0)}},
		{vals: [2]interface{}{int8(-1), int8(0)}},
		{vals: [2]interface{}{int16(-1), int16(0)}},
		{vals: [2]interface{}{int32(-1), int32(0)}},
		{vals: [2]interface{}{int64(-1), int64(0)}},
		{vals: [2]interface{}{uint(0), uint(1)}},
		{vals: [2]interface{}{uint8(0), uint8(1)}},
		{vals: [2]interface{}{uint16(0), uint16(1)}},
		{vals: [2]interface{}{uint32(0), uint32(1)}},
		{vals: [2]interface{}{uint64(0), uint64(1)}},
		{vals: [2]interface{}{false, true}},
		{vals: [2]interface{}{now, now.Add(time.Second)}},
		// nullable
		{vals: [2]interface{}{nilptr("string"), nilptr("string")}},
		{vals: [2]interface{}{nilptr("string"), ptr("a")}},
		{vals: [2]interface{}{ptr(""), nilptr("string")}},
		{vals: [2]interface{}{ptr(""), ptr("")}},
		{vals: [2]interface{}{ptr(""), ptr("a")}},
		{vals: [2]interface{}{nilptr("int"), nilptr("int")}},
		{vals: [2]interface{}{nilptr("int"), ptr(int(0))}},
		{vals: [2]interface{}{ptr(int(-1)), nilptr("int")}},
		{vals: [2]interface{}{ptr(int(-1)), ptr(int(-1))}},
		{vals: [2]interface{}{ptr(int(-1)), ptr(int(0))}},
		{vals: [2]interface{}{nilptr("int8"), nilptr("int8")}},
		{vals: [2]interface{}{nilptr("int8"), ptr(int8(0))}},
		{vals: [2]interface{}{ptr(int8(-1)), nilptr("int8")}},
		{vals: [2]interface{}{ptr(int8(-1)), ptr(int8(-1))}},
		{vals: [2]interface{}{ptr(int8(-1)), ptr(int8(0))}},
		{vals: [2]interface{}{nilptr("int16"), nilptr("int16")}},
		{vals: [2]interface{}{nilptr("int16"), ptr(int16(0))}},
		{vals: [2]interface{}{ptr(int16(-1)), nilptr("int16")}},
		{vals: [2]interface{}{ptr(int16(-1)), ptr(int16(-1))}},
		{vals: [2]interface{}{ptr(int16(-1)), ptr(int16(0))}},
		{vals: [2]interface{}{nilptr("int32"), nilptr("int32")}},
		{vals: [2]interface{}{nilptr("int32"), ptr(int32(0))}},
		{vals: [2]interface{}{ptr(int32(-1)), nilptr("int32")}},
		{vals: [2]interface{}{ptr(int32(-1)), ptr(int32(-1))}},
		{vals: [2]interface{}{ptr(int32(-1)), ptr(int32(0))}},
		{vals: [2]interface{}{nilptr("int64"), nilptr("int64")}},
		{vals: [2]interface{}{nilptr("int64"), ptr(int64(0))}},
		{vals: [2]interface{}{ptr(int64(-1)), nilptr("int64")}},
		{vals: [2]interface{}{ptr(int64(-1)), ptr(int64(-1))}},
		{vals: [2]interface{}{ptr(int64(-1)), ptr(int64(0))}},
		{vals: [2]interface{}{nilptr("uint"), nilptr("uint")}},
		{vals: [2]interface{}{nilptr("uint"), ptr(uint(0))}},
		{vals: [2]interface{}{ptr(uint(0)), nilptr("uint")}},
		{vals: [2]interface{}{ptr(uint(0)), ptr(uint(0))}},
		{vals: [2]interface{}{ptr(uint(0)), ptr(uint(1))}},
		{vals: [2]interface{}{nilptr("uint8"), nilptr("uint8")}},
		{vals: [2]interface{}{nilptr("uint8"), ptr(uint8(0))}},
		{vals: [2]interface{}{ptr(uint8(0)), nilptr("uint8")}},
		{vals: [2]interface{}{ptr(uint8(0)), ptr(uint8(0))}},
		{vals: [2]interface{}{ptr(uint8(0)), ptr(uint8(1))}},
		{vals: [2]interface{}{nilptr("uint16"), nilptr("uint16")}},
		{vals: [2]interface{}{nilptr("uint16"), ptr(uint16(0))}},
		{vals: [2]interface{}{ptr(uint16(0)), nilptr("uint16")}},
		{vals: [2]interface{}{ptr(uint16(0)), ptr(uint16(0))}},
		{vals: [2]interface{}{ptr(uint16(0)), ptr(uint16(1))}},
		{vals: [2]interface{}{nilptr("uint32"), nilptr("uint32")}},
		{vals: [2]interface{}{nilptr("uint32"), ptr(uint32(0))}},
		{vals: [2]interface{}{ptr(uint32(0)), nilptr("uint32")}},
		{vals: [2]interface{}{ptr(uint32(0)), ptr(uint32(0))}},
		{vals: [2]interface{}{ptr(uint32(0)), ptr(uint32(1))}},
		{vals: [2]interface{}{nilptr("uint64"), nilptr("uint64")}},
		{vals: [2]interface{}{nilptr("uint64"), ptr(uint64(0))}},
		{vals: [2]interface{}{ptr(uint64(0)), nilptr("uint64")}},
		{vals: [2]interface{}{ptr(uint64(0)), ptr(uint64(0))}},
		{vals: [2]interface{}{ptr(uint64(0)), ptr(uint64(1))}},
		{vals: [2]interface{}{nilptr("bool"), nilptr("bool")}},
		{vals: [2]interface{}{nilptr("bool"), ptr(false)}},
		{vals: [2]interface{}{ptr(false), nilptr("bool")}},
		{vals: [2]interface{}{ptr(false), ptr(false)}},
		{vals: [2]interface{}{ptr(false), ptr(true)}},
		{vals: [2]interface{}{nilptr("time.Time"), nilptr("time.Time")}},
		{vals: [2]interface{}{nilptr("time.Time"), ptr(now)}},
		{vals: [2]interface{}{ptr(now), ptr(now)}},
		{vals: [2]interface{}{ptr(now), ptr(now.Add(time.Second))}},
	}

	// Add attributes to type
	typ := &Type{Name: "type"}
	for i, t := range attrs {
		ti, null := GetAttrType(fmt.Sprintf("%T", t.vals[0]))
		typ.AddAttr(Attr{
			Name:     "attr" + strconv.Itoa(i),
			Type:     ti,
			Nullable: null,
		})
	}
	sc.SetType(typ)

	// Add resources
	for i := range attrs {
		sr := &SoftResource{}
		sr.SetType(typ)
		sr.SetID("id" + strconv.Itoa(i))
		for j := range attrs {
			if i != j {
				sr.Set("attr"+strconv.Itoa(j), attrs[j].vals[0])
			} else {
				sr.Set("attr"+strconv.Itoa(j), attrs[j].vals[1])
			}
		}
		sc.Add(sr)
	}

	// Sort collection
	rules := []string{}
	for i := 0; i < sc.Len(); i++ {
		reverse := ""
		if i%3 == 0 {
			reverse = "-"
		}
		rules = append(rules, reverse+"attr"+strconv.Itoa(i))
	}
	rules = append(rules, "id")
	sc.Sort(rules)

	// Sorted IDs from the collection
	ids := []string{}
	for i := 0; i < sc.Len(); i++ {
		ids = append(ids, sc.Elem(i).GetID())
	}

	expectedIDs := []string{
		"id0", "id3", "id6", "id9", "id12", "id20", "id24", "id25", "id27",
		"id35", "id39", "id40", "id42", "id50", "id54", "id55", "id57",
		"id69", "id70", "id72", "id10", "id13", "id16", "id18", "id21",
		"id23", "id26", "id28", "id31", "id33", "id36", "id38", "id41",
		"id43", "id46", "id48", "id51", "id53", "id56", "id58", "id61",
		"id63", "id64", "id65", "id66", "id67", "id68", "id71", "id73",
		"id75", "id76", "id74", "id62", "id60", "id59", "id52", "id49",
		"id47", "id45", "id44", "id37", "id34", "id32", "id30", "id29",
		"id22", "id19", "id17", "id15", "id14", "id11", "id8", "id7", "id5",
		"id4", "id2", "id1",
	}
	assert.Equal(expectedIDs, ids, fmt.Sprintf("sort with rules: %v", rules))

	// Sort with an empty list of sorting rules.
	sc.Sort([]string{})

	ids = []string{}
	for i := 0; i < sc.Len(); i++ {
		ids = append(ids, sc.Elem(i).GetID())
	}

	sort.Strings(expectedIDs)
	assert.Equal(expectedIDs, ids, "sort by ID")
}

func TestSoftCollectionMiscellaneous(t *testing.T) {
	assert := assert.New(t)

	sc := &SoftCollection{}
	assert.Nil(sc.Elem(99), "nonexistent element")
}
